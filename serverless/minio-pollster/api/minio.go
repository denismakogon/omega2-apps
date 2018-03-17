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
	"github.com/denismakogon/omega2-apps/serverless/common"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type store struct {
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	bucket     string
}

func createStore(bucketName, endpoint, region, accessKeyID, secretAccessKey string, useSSL bool) *store {
	client := s3.New(session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(!useSSL),
		S3ForcePathStyle: aws.Bool(true),
	})))
	return &store{
		client:     client,
		uploader:   s3manager.NewUploaderWithClient(client),
		downloader: s3manager.NewDownloaderWithClient(client),
		bucket:     bucketName,
	}
}

func New(u *url.URL) (*store, error) {
	endpoint := u.Host

	var accessKeyID, secretAccessKey string
	if u.User != nil {
		accessKeyID = u.User.Username()
		secretAccessKey, _ = u.User.Password()
	}
	useSSL := u.Query().Get("ssl") == "true"

	strs := strings.SplitN(u.Path, "/", 3)
	if len(strs) < 3 {
		return nil, errors.New("must provide bucket name and region in path of s3 api url. e.g. s3://s3.com/us-east-1/my_bucket")
	}
	region := strs[1]
	bucketName := strs[2]
	if region == "" {
		return nil, errors.New("must provide non-empty region in path of s3 api url. e.g. s3://s3.com/us-east-1/my_bucket")
	} else if bucketName == "" {
		return nil, errors.New("must provide non-empty bucket name in path of s3 api url. e.g. s3://s3.com/us-east-1/my_bucket")
	}

	logrus.WithFields(logrus.Fields{"bucketName": bucketName, "region": region, "endpoint": endpoint, "access_key_id": accessKeyID, "useSSL": useSSL}).Info("checking / creating s3 bucket")
	store := createStore(bucketName, endpoint, region, accessKeyID, secretAccessKey, useSSL)

	_, err := store.client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucketName)})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyOwnedByYou, s3.ErrCodeBucketAlreadyExists:
				// bucket already exists, NO-OP
			default:
				return nil, fmt.Errorf("failed to create bucket %s: %s", bucketName, aerr.Message())
			}
		} else {
			return nil, fmt.Errorf("unexpected error creating bucket %s: %s", bucketName, err.Error())
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
	log.Info("Current query key: ", *result.Marker)
	log.Info("Next query key: ", *result.NextMarker)
	log.Info("Found object: ", len(result.Contents))
	if len(result.Contents) > 0 {
		wg.Add(len(result.Contents))
		for _, object := range result.Contents {

			go func(wg sync.WaitGroup, object *s3.Object) {
				defer wg.Done()

				target := &aws.WriteAtBuffer{}
				log.Info("Pulling the content of the object: ", s.bucket+"/"+*object.Key)
				size, err := s.downloader.DownloadWithContext(ctx, target, &s3.GetObjectInput{
					Bucket: aws.String(s.bucket),
					Key:    object.Key,
				})
				if err != nil {
					log.Fatal(err.Error())
					os.Exit(1)
				}
				req.Header.Set("Content-Length", strconv.FormatInt(size, 10))
				//payload := &api.RequestPayload{MediaContent: string(target.Bytes())}
				//err = api.DoRequest(payload, req, httpClient, fnToken)
				//if err != nil {
				//	log.Fatal(err.Error())
				//	os.Exit(1)
				//}

			}(wg, object)
		}
		input.SetMarker(*result.NextMarker)
	}

	return nil
}

func (s *store) DispatchObjects(ctx context.Context, wg sync.WaitGroup, appName string) error {
	fnAPIURL, fnToken, err := setupEmokognitionV2(ctx, appName)
	if err != nil {
		return err
	}

	input := &s3.ListObjectsInput{
		Bucket:  aws.String(s.bucket),
		MaxKeys: aws.Int64(10),
		Marker:  aws.String(""),
	}

	detect, err := http.NewRequest(
		http.MethodPost, fmt.Sprintf("%s/r/%s/detect-v2", fnAPIURL, appName),
		nil)
	if err != nil {
		return err
	}
	httpClient := common.SetupHTTPClient()
	log := logrus.WithFields(logrus.Fields{"bucketName": s.bucket})
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
