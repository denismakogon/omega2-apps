# Just builds
.PHONY: all test dep build test-log-datastore

dep:
	dep ensure

dep-up:
	dep ensure -update

build:
	go build -race
	rm -fr twitter_daemon

build-mipsle:
	GOOS=linux GOARCH=mipsle go build -ldflags "-s -w" -compiler gc
	rm -fr twitter_daemon
