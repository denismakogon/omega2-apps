package main

import (
	"context"
	"fmt"
	"github.com/denismakogon/omega2-apps/serverless/minio-pollster/api"
	"net/url"
	"os"
	"sync"
)

func start() error {
	ctx := context.Background()
	minioURL := os.Getenv("MINIO_URL")
	u, err := url.Parse(minioURL)
	if err != nil {
		return err
	}

	minio, err := api.New(u)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	err = minio.DispatchObjects(ctx, wg, "emokognition-v2")
	return err
}

func main() {
	err := start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
