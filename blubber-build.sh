#!/usr/bin/env bash

# clean up previous builds
docker rm wikispeech-pronlex-test
docker rmi --force wikispeech-pronlex-test

docker rm wikispeech-pronlex
docker rmi --force wikispeech-pronlex

# build docker
docker build --tag wikispeech-pronlex-test --file .pipeline/blubber.yaml --target test .
docker build --tag wikispeech-pronlex --file .pipeline/blubber.yaml --target production .
