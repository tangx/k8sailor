#!/bin/bash


set -eux

cd $(dirname $0)

find . -type f -name "*.md" |xargs sed -i 's@(./assets@(../assets@g'
# find . -type f -name "*.md" |xargs sed -i 's@(../assets@(../../assets@g'