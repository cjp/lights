#!/bin/sh

set -eux -o pipefail

GOOS=linux GOARCH=arm GOARM=5 go build -v github.com/cjp/lights
