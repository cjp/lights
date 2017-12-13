#!/bin/sh

set -eux -o pipefail

tempdir=$(mktemp -d)
echo "${tempdir}" | grep -q '/tmp.'
cd "${tempdir}"

GOOS=linux GOARCH=arm GOARM=5 go build -v github.com/cjp/lights
gzip lights
scp lights.gz raspi:
rm -f lights
gunzip -f lights.gz
ssh raspi 'killall -q lights || true'
ssh raspi "nohup ./lights > /dev/null 2>&1 &"

cd
rm -rf "${tempdir}"
