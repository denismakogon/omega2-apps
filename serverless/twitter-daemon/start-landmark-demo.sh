#!/usr/bin/env bash

source ${GOPATH}/../omage2_twitter.new.rc

export TwitterBotType="landmark-recognition"
export InitialTweetID=978976524434067456

./twitter-daemon
