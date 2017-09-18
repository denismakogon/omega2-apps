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
	cd serverless/emotion-recognition/emotion-results; \
	go build -race; \
	rm -fr emotion-recorder

dep-emotion-recorder:
	cd serverless/emotion-recognition/emotion-results; \
	dep ensure

dep-emotion-recorder-up:
	cd serverless/emotion-recognition/emotion-results; \
	dep ensure -update

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
	cd serverless/landmark-recognition/tweet-fail; \
	${GOPATH}/bin/dep ensure; \
	go build -race; \
	rm -fr tweet-fail; \
	docker build -t denismakogon/tweet-fail:latest .; \
	docker rmi denismakogon/tweet-fail:latest

ci-build-tweet-success:
	cd serverless/landmark-recognition/tweet-success; \
	${GOPATH}/bin/dep ensure; \
	go build -race; \
	rm -fr tweet-success; \
	docker build -t denismakogon/tweet-success:latest .; \
	docker rmi denismakogon/tweet-success:latest

ci-build-emotion-recorder:
	cd serverless/emotion-recognition/emotion-recorder; \
	${GOPATH}/bin/dep ensure; \
	go build; \
	rm -fr emotion-recorder; \
	docker build -t denismakogon/emotion-recorder:latest .; \
	docker rmi denismakogon/emotion-recorder:latest


all: dep-twitter-daemon dep-twitter-daemon-up build-twitter-daemon