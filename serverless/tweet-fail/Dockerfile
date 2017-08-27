FROM funcy/go:dev
RUN mkdir -p $GOPATH/src/github.com/denismakogon/tweet-task
ADD . $GOPATH/src/github.com/denismakogon/tweet-task
RUN cd $GOPATH/src/github.com/denismakogon/tweet-task && go build -o func
WORKDIR $GOPATH/src/github.com/denismakogon/tweet-task
ENTRYPOINT ["./func"]
