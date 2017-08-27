# Just builds
.PHONY: all

dep-twitter-daemon:
	cd serverless/twitter-daemon; \
	dep ensure

dep-twitter-daemon-up:
	cd serverless/twitter-daemon; \
	dep ensure -update

build-twitter-daemon:
	cd serverless/twitter-daemon; \
	go build -race; \
	rm -fr twitter_daemon

build-twitter-daemon-mipsle:
	cd serverless/twitter-daemon; \
	GOOS=linux GOARCH=mipsle go build -ldflags "-s -w" -compiler gc; \
	rm -fr twitter_daemon

ci-build-twitter-daemon:
	cd serverless/twitter-daemon; \
	dep ensure; \
	go build -race; \
	rm -fr twitter_daemon

ci-build-twitter-daemon-mipsle:
	cd serverless/twitter-daemon; \
	dep ensure; \
	GOOS=linux GOARCH=mipsle go build -ldflags "-s -w" -compiler gc; \
	rm -fr twitter_daemon
