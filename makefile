#!/usr/bin/make
## makefile (for rf-tools)
## Copyright 2017 Mac Radigan
## All Rights Reserved

.PHONY: build run bootstrap \
   alt-1 alt-2              \
   seq-1 seq-3   seq-inf    \
  test-1 test-3 test-inf
.DEFAULT_GOAL := build

target := rf

## 1024 bytes per frame x 3.5 frames = 3584 bytes
nbytes := 3584

build:
	go $@ $(target).go

run: build
	go $@ $(target).go

alt-1:
	@yes `seq -s '' 9` | tr -d '\n' | head -c $(nbytes) | ./$(target)

alt-2:
	@./tests/fiducial.m | ./$(target)

## 1 sequence, N=1024, 3 patterns, 256 resdiual bytes
seq-1:
	@./tests/fiducial.m 1024 256 3 1
test-1:
	@./tests/fiducial.m 1024 256 3 1 | ./$(target)

## 3 sequences, N=1024, 2 patterns, 256 resdiual bytes
seq-3:
	@./tests/fiducial.m 1024 256 2 3
test-3:
	@./tests/fiducial.m 1024 256 2 3 | ./$(target)

## infinite sequence, N=1024, 15 patterns
seq-inf:
	@./tests/fiducial.m 1024 0 15 0
test-inf:
	@./tests/fiducial.m 1024 0 15 0 | ./$(target)

bootstrap:
	go get -u github.com/golang/glog

## *EOF*
