FROM funcy/go:dev
RUN go get -u github.com/golang/dep/cmd/dep
RUN mkdir -p $GOPATH/src/github.com/denismakogon/tweet-dispatcher
ADD . $GOPATH/src/github.com/denismakogon/tweet-dispatcher
WORKDIR $GOPATH/src/github.com/denismakogon/tweet-dispatcher
RUN $GOPATH/bin/dep ensure
RUN go build -o func
ENTRYPOINT ["./func"]
