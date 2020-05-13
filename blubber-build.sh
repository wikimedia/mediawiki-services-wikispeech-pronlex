#!/usr/bin/env bash

# clean up previous builds
docker rm wikispeech-pronlex-test
docker rmi --force wikispeech-pronlex-test

docker rm wikispeech-pronlex
docker rmi --force wikispeech-pronlex

# build docker
blubber .pipeline/blubber.yaml test | docker build --tag wikispeech-pronlex-test --file - .
blubber .pipeline/blubber.yaml production | docker build --tag wikispeech-pronlex --file - .

