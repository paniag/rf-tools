#!/usr/bin/make
## makefile (for rf-tools)
## Copyright 2017 Mac Radigan
## All Rights Reserved
## Mac Radigan

.PHONY: build run bootstrap test
.DEFAULT_GOAL := build

target := rf

## 1024 bytes per frame x 3.5 frames = 3584 bytes
nbytes := 3584

build:
	go $@ $(target).go

run: build
	go $@ $(target).go

test:
	@yes `seq -s '' 9` | tr -d '\n' | head -c $(nbytes) | ./$(target)

test2:
	@cat ./data/test.dat | ./$(target)

bootstrap:
	go get -u github.com/go-mangos/mangos

## *EOF*
