#!/bin/bash

mkdir -v testdata testresults

./run-test-system.sh
./run-test-create.sh
./run-test-create-upload.sh
./run-test-create-upload-update-read.sh
./run-test-create-upload-update.sh

