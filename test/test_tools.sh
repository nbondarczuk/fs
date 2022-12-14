#!/bin/bash

CONTROL_FILE_NAME=
CONTROL_FILE_CONTENT=
CONTROL_FILE_FORMAT="{\"name\": \"%s\",\"check_sum\":\"%s\",\"size\": %d}"
TMP_FILE_NAME=

# make a control file for a given test file
# rv ->  CONTROL_FILE_NAME
make_control_file() {
	local filename filepathname cksum size
	filename="$1"
	filepathname="testdata/$filename"
	cksum=$(sha256sum "$filepathname" | awk '{print $1}')
	size=$(stat --format=%s "$filepathname")
	CONTROL_FILE_CONTENT=$(printf "$CONTROL_FILE_FORMAT" "$filename" "$cksum" "$size")
	CONTROL_FILE_NAME="${filepathname}.json"
	echo "$CONTROL_FILE_CONTENT" >"$CONTROL_FILE_NAME"
	return 0
}

# make a test file with size of N blocks
# rv -> TMP_FILE_NAME
make_test_file() {
	local blocks block_size filename 
	blocks=$(( "$1" * $RANDOM ))
	block_size="$2"
	filename=$(date +"%Y%m%d%H%M%SXXXXXX")
	TMP_FILE_NAME=$(mktemp "testdata/${filename}")
	dd bs="${block_size}" count="${blocks}" </dev/urandom >"${TMP_FILE_NAME}"
	return 0
}

# Extract key/value form parameters from the reply, format curl command and run it
upload_test_file() {
	local reply tmp_file_name
	tmp_file_name="$1"
	reply="$2"
	url=$(echo "${reply}" | jq '.data[0].url' | tr -d \")
	# Curl magic to put the file size in the request header constraints
	cmd="curl -X PUT -T ${tmp_file_name} -D - ${url}"
	echo "Running: $cmd"
	$cmd
	return $?
}	
