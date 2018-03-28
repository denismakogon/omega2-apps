package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/denismakogon/omega2-apps/serverless/minio-pollster/common"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type store struct {
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	bucket     string
	config     *MinioConfig
}

func (m *MinioConfig) createStore() *store {
	client := s3.New(session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(m.AccessKeyID, m.SecretAccessKey, ""),
		Endpoint:         aws.String(m.Endpoint),
		Region:           aws.String(m.Region),
		DisableSSL:       aws.Bool(!m.UseSSL),
		S3ForcePathStyle: aws.Bool(true),
	})))
	return &store{
		client:     client,
		config:     m,
		uploader:   s3manager.NewUploaderWithClient(client),
		downloader: s3manager.NewDownloaderWithClient(client),
	}
}

type MinioConfig struct {
	Bucket          string `json:"bucket"`
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	UseSSL          bool   `json:"use_ssl"`
}

func (m *MinioConfig) FromURL(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	endpoint := u.Host

	var accessKeyID, secretAccessKey string
	if u.User != nil {
		accessKeyID = u.User.Username()
		secretAccessKey, _ = u.User.Password()
	}
	useSSL := u.Query().Get("ssl") == "true"

	strs := strings.SplitN(u.Path, "/", 3)
	if len(strs) < 3 {
		return errors.New("must provide bucket name and region in path of s3 api url. e.g. s3://s3.com/us-east-1/my_bucket")
	}
	region := strs[1]
	bucketName := strs[2]
	if region == "" {
		return errors.New("must provide non-empty region in path of s3 api url. e.g. s3://s3.com/us-east-1/my_bucket")
	} else if bucketName == "" {
		return errors.New("must provide non-empty bucket name in path of s3 api url. e.g. s3://s3.com/us-east-1/my_bucket")
	}

	m.Bucket = bucketName
	m.Endpoint = endpoint
	m.Region = region
	m.AccessKeyID = accessKeyID
	m.SecretAccessKey = secretAccessKey
	m.UseSSL = useSSL

	return nil
}

func (m *MinioConfig) ToMap() (map[string]interface{}, error) {
	return common.ToMap(m)
}

func withDefault(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func New() (*store, error) {
	m := &MinioConfig{}

	err := m.FromURL(withDefault("MINIO_URL",
		"s3://admin:password@localhost:9000/us-east-1/emotions"))
	if err != nil {
		return nil, err
	}
	logFields, err := m.ToMap()
	if err != nil {
		return nil, err
	}

	logrus.WithFields(logFields).Info("checking / creating s3 bucket")

	store := m.createStore()

	_, err = store.client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(m.Bucket)})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyOwnedByYou, s3.ErrCodeBucketAlreadyExists:
				// bucket already exists, NO-OP
			default:
				return nil, fmt.Errorf("failed to create bucket %s: %s", m.Bucket, aerr.Message())
			}
		} else {
			return nil, fmt.Errorf("unexpected error creating bucket %s: %s", m.Bucket, err.Error())
		}
	}

	return store, nil
}

func (s *store) asyncDispatcher(ctx context.Context, wg sync.WaitGroup, log *logrus.Entry, input *s3.ListObjectsInput,
	req *http.Request, httpClient *http.Client, fnToken string) error {

	result, err := s.client.ListObjectsWithContext(ctx, input)
	if err != nil {
		return err
	}
	fields := logrus.Fields{}
	fields["current_key"] = *result.Marker
	fields["objects_found"] = len(result.Contents)
	if result.NextMarker != nil {
		fields["next_query_key"] = *result.NextMarker
	}
	log = log.WithFields(fields)
	if len(result.Contents) > 0 {
		wg.Add(len(result.Contents))
		for _, object := range result.Contents {

			go func(wg sync.WaitGroup, object *s3.Object) {
				defer wg.Done()

				getObjectReq, _ := s.client.GetObjectRequest(
					&s3.GetObjectInput{
						Bucket: &s.config.Bucket,
						Key:    object.Key,
					},
				)

				urlStr, err := getObjectReq.Presign(15 * time.Minute)
				if err != nil {
					log.Fatal(err.Error())
					os.Exit(1)
				}

				log.Info("Sending the object: ", s.config.Bucket+"/"+*object.Key)
				log.Info("Presigned object URL: ", urlStr)

				payload := &common.RequestPayload{MediaURL: urlStr}

				err = common.DoRequest(log, payload, req, httpClient, fnToken)
				if err != nil {
					log.Fatal(err.Error())
					os.Exit(1)
				}

			}(wg, object)
		}
		input.SetMarker(*result.NextMarker)
	}

	return nil
}

func (s *store) DispatchObjects(ctx context.Context, wg sync.WaitGroup, appName string) error {
	log := logrus.WithFields(logrus.Fields{"bucketName": s.config.Bucket})
	config := map[string]string{}

	pgConf := new(common.PostgresConfig)
	err := pgConf.FromEnv()
	if err != nil {
		return err
	}
	log.Info("Postgres config provisioned")

	config, err = common.Append(pgConf, config)
	if err != nil {
		return err
	}

	m := &MinioConfig{}
	err = m.FromURL(os.Getenv("INTERNAL_MINIO_URL"))
	if err != nil {
		return err
	}
	log.Info("Internal Minio config provisioned")

	config, err = common.Append(m, config)
	if err != nil {
		return err
	}

	fnAPIURL, fnToken, err := setupEmokognitionV2(ctx, appName, config)
	if err != nil {
		return err
	}

	input := &s3.ListObjectsInput{
		Bucket:  aws.String(s.config.Bucket),
		MaxKeys: aws.Int64(10),
		Marker:  aws.String(""),
	}

	detect, err := http.NewRequest(
		http.MethodPost, fmt.Sprintf("%s/r/%s/detect", fnAPIURL, appName),
		nil)
	if err != nil {
		return err
	}
	httpClient := common.SetupHTTPClient()
	for {

		err = s.asyncDispatcher(ctx, wg, log, input, detect, httpClient, fnToken)
		if err != nil {
			return err
		}

		time.Sleep(5 * time.Second)
	}

	wg.Wait()

	return nil
}
