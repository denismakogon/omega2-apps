FROM funcy/go:dev
RUN mkdir -p $GOPATH/src/github.com/denismakogon/tweet-dispatcher
ADD . $GOPATH/src/github.com/denismakogon/tweet-dispatcher
RUN cd $GOPATH/src/github.com/denismakogon/tweet-dispatcher && go build -o func
WORKDIR $GOPATH/src/github.com/denismakogon/tweet-dispatcher
ENTRYPOINT ["./func"]
