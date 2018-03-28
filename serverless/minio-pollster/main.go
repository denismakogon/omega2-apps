package main

import (
	"context"
	"fmt"
	"github.com/denismakogon/omega2-apps/serverless/minio-pollster/api"
	"os"
	"sync"
)

func start() error {

	ctx := context.Background()
	minio, err := api.New()
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
