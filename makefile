#!/usr/bin/make
## makefile (for rf-tools)
## Copyright 2017 Mac Radigan
## All Rights Reserved
## Mac Radigan

.PHONY: build run bootstrap test
.DEFAULT_GOAL := build

target := rf

build:
	go $@ $(target).go

run: build
	go $@ $(target).go

test:
	@echo "TEST" | ./$(target)

bootstrap:
	go get -u github.com/go-mangos/mangos

## *EOF*
