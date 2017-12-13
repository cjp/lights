#!/bin/sh

set -eux -o pipefail

gzip lights
scp lights.gz raspi:
rm -f lights
gunzip -f lights.gz
ssh raspi 'killall -q lights || true'
ssh raspi "nohup ./lights > /dev/null 2>&1 &"
