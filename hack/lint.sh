#!/bin/bash
# Copyright 2019 Iguazio
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -e

cd $GOPATH/src/github.com/v3io/v3ctl

echo Installing impi
go get -u github.com/pavius/impi/cmd/impi

echo Linting imports with impi
$GOPATH/bin/impi \
    --local github.com/v3io/v3ctl \
    --scheme stdLocalThirdParty \
    ./pkg/...

echo Getting all packages
go get ./...

echo Installing gometalinter
go get -u gopkg.in/alecthomas/gometalinter.v2
$GOPATH/bin/gometalinter.v2 --install

echo Linting with gometalinter
$GOPATH/bin/gometalinter.v2 \
    --deadline=300s \
    --disable-all \
    --enable-gc \
    --enable=deadcode \
    --enable=goconst \
    --enable=gofmt \
    --enable=golint \
    --enable=gosimple \
    --enable=ineffassign \
    --enable=interfacer \
    --enable=misspell \
    --enable=staticcheck \
    --enable=unconvert \
    --enable=varcheck \
    --enable=vet \
    --enable=vetshadow \
    --enable=errcheck \
    --exclude="_test.go" \
    --exclude="comment on" \
    --exclude="error should be the last" \
    --exclude="should have comment" \
    --skip=pkg/platform/kube/apis \
    --skip=pkg/platform/kube/client \
    ./pkg/...
