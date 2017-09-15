FROM funcy/go:dev
RUN mkdir -p $GOPATH/src/github.com/denismakogon/voteapphot
ADD . $GOPATH/src/github.com/denismakogon/voteapphot
RUN cd $GOPATH/src/github.com/denismakogon/voteapphot && go build -o func
WORKDIR $GOPATH/src/github.com/denismakogon/voteapphot
ENTRYPOINT ["./func"]
