#LDFLAGS := $(shell go run buildscripts/gen-ldflags.go)
PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
#BUILD_LDFLAGS := '$(LDFLAGS)'
TAG := latest

HOST ?= $(shell uname)
CPU ?= $(shell uname -m)
# list of pkgs for the project without vendor
PKGS=$(shell go list ./... | grep -v /vendor/)

# if no host is identifed (no uname tool)
# we assume a Linux-64bit build
ifeq ($(HOST),)
  HOST = Linux
endif

# identify CPU
ifeq ($(CPU), x86_64)
  HOST := $(HOST)64
else
ifeq ($(CPU), amd64)
  HOST := $(HOST)64
else
ifeq ($(CPU), i686)
  HOST := $(HOST)32
endif
endif
endif


#############################################
# now we find out the target OS for
# which we are going to compile in case
# the caller didn't yet define OS himself
ifndef (OS)
  ifeq ($(HOST), Linux64)
    arch = gcc
  else
  ifeq ($(HOST), Linux32)
    arch = 32
  else
  ifeq ($(HOST), Darwin64)
    arch = clang
  else
  ifeq ($(HOST), Darwin32)
    arch = clang
  else
  ifeq ($(HOST), FreeBSD64)
    arch = gcc
  endif
  endif
  endif
  endif
  endif
endif

# This how we want to name the binary output
BINARY=$(CURDIR)/bin/server

# These are the values we want to pass for VERSION and BUILD
VERSION=0.1.0-${HOST}-${CPU}
BUILD=`git rev-parse HEAD`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

all: checks getdeps vendor build install

checks:
	@(env bash $(PWD)/buildscripts/check_tools.sh)

getdeps:
	@echo "Installing glide:"
	@curl https://glide.sh/get | sh

# Builds the project
build:
	go build ${LDFLAGS} -o ${BINARY} server.go

# Builds & run the project
run:
	go run ${LDFLAGS} server.go

format:
	@go fmt $(PKGS)

# Installs our project: copies binaries
install:
	go install ${LDFLAGS}

vendor:
	glide install

jwt-cert-dev:
	openssl genrsa -out ./settings/keys/jwt-private-dev.pem 1024
	openssl rsa -in ./settings/keys/jwt-private-dev.pem -outform PEM -pubout -out ./settings/keys/jwt-public-dev.pem

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
