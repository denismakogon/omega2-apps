package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/denismakogon/omega2-apps/serverless/minio-pollster/api"
	"net/url"
	"os"
	"sync"
)

func withDefault(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func start() error {
	minioURL := withDefault("MINIO_URL", "s3://admin:password@localhost:9000/us-east-1/emotions")
	minioURL = *flag.String("s3-url", minioURL, "S3 API-compatible (minio, swift, etc.) store URL")

	flag.Parse()

	ctx := context.Background()
	u, err := url.Parse(minioURL)
	if err != nil {
		return err
	}

	minio, err := api.New(u)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	return minio.DispatchObjects(ctx, wg, "emokognition-v2")
}

func main() {
	err := start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
