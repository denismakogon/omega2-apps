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
	rm -fr twitter-daemon

build-twitter-daemon-mipsle:
	cd serverless/twitter-daemon; \
	GOOS=linux GOARCH=mipsle go build -ldflags "-s -w" -compiler gc; \
	rm -fr twitter-daemon

build-emotion-recorder:
	cd serverless/emotion-results; \
	go build -race; \
	rm -fr emotion-recorder

ci-build-twitter-daemon:
	cd serverless/twitter-daemon; \
	${GOPATH}/bin/dep ensure; \
	go build -race; \
	rm -fr twitter-daemon

ci-build-twitter-daemon-mipsle:
	cd serverless/twitter-daemon; \
	${GOPATH}/bin/dep ensure; \
	GOOS=linux GOARCH=mipsle go build -ldflags "-s -w" -compiler gc; \
	rm -fr twitter-daemon

ci-build-tweet-fail:
	cd serverless/tweet-fail; \
	${GOPATH}/bin/dep ensure; \
	go build -race; \
	rm -fr tweet-fail

ci-build-tweet-success:
	cd serverless/tweet-success; \
	${GOPATH}/bin/dep ensure; \
	go build -race; \
	rm -fr tweet-success

ci-build-tweet-dispatcher:
	cd serverless/tweet-dispatcher; \
	${GOPATH}/bin/dep ensure; \
	go build -race; \
	rm -fr tweet-dispatcher

all: dep-twitter-daemon dep-twitter-daemon-up build-twitter-daemon