FROM funcy/go:dev
RUN mkdir -p $GOPATH/src/github.com/denismakogon/tweet-success
ADD . $GOPATH/src/github.com/denismakogon/tweet-success
RUN cd $GOPATH/src/github.com/denismakogon/tweet-success && go build -o func
WORKDIR $GOPATH/src/github.com/denismakogon/tweet-success
ENTRYPOINT ["./func"]
