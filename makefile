#!/usr/bin/make
## makefile (for go-wfg)
## Copyright 2017 Mac Radigan
## All Rights Reserved
## Mac Radigan

.PHONY: build run bootstrap
.DEFAULT_GOAL := build

target := wfg

build:
	go $@ $(target).go

run: build
	go $@ $(target).go

bootstrap:
	go get -u github.com/go-mangos/mangos

## *EOF*
