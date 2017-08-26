# Just builds
.PHONY: all

dep-twitter-daemon:
	cd serverless/twitter_daemon; \
	dep ensure

dep-twitter-daemon-up:
	cd serverless/twitter_daemon; \
	dep ensure -update

build-twitter-daemon:
	cd serverless/twitter_daemon; \
	go build -race; \
	rm -fr twitter_daemon

build-twitter-daemon-mipsle:
	cd serverless/twitter_daemon; \
	GOOS=linux GOARCH=mipsle go build -ldflags "-s -w" -compiler gc; \
	rm -fr twitter_daemon
