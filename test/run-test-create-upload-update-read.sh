#!/bin/bash

TS=$(date +"%Y%m%d%H%M%S")
HOST=localhost
PORT=1234
HEADER="Content-Type: application/json"
TENANT=xxx
DEVICE=yyy
BLOCK_SIZE=1
DATA=
TMP_FILE_NAME=

# shellcheck source=test_tools.sh
. test_tools.sh

CONTROL_FILE_FORMAT="{\"name\": \"%s\",\"check_sum\":\"%s\",\"size\": %d,\"status\":\"U\"}"

# Perform one test, sending the request to an API, checki the reply and try upload a file
function run_test() {
	local testno fn rfn status rc reply msg rfn url id

	testno="$1"

	#
	# Make a test file with random content of size N blocks
	#
	make_test_file "${testno}" "${BLOCK_SIZE}"
	fn="${TMP_FILE_NAME}"
	bfn="$(basename "${fn}")"
	make_control_file "${bfn}"
	DATA="${CONTROL_FILE_NAME}"

	#
	# Send request to the fs API
	#
	url="http://${HOST}:${PORT}/api/v1/files?tenant=${TENANT}&device=${DEVICE}"
	echo "Running: curl -X POST -H ${HEADER} -d@${DATA} ${url}"
	reply=$(curl -X POST -H "${HEADER}" -d@"${DATA}" "${url}" 2>/dev/null)
	rc=$?
	if [ "${rc}" -ne 0 ]; then
		echo "### Error ###"
		return 1
	fi	
	msg=$(echo "${reply}" | jq)
	rfn="testresults/results-create-${TS}-step-1-create.lst"
	echo "Result: ${rc} -> ${msg}" | tee "${rfn}"
	# Check the results
	status=$(echo "${msg}" | jq '.status')
	if [ "${status}" != "true" ]; then
		echo "### Error ###"
		return 1
	fi
	id=$(echo "${reply}" | jq '.data[0].ID' | tr -d \")
	
	#
	# Upload test file testing the presigned URL
	#
	upload_test_file "${TMP_FILE_NAME}" "${reply}"
	rc=$?
	if [ "${rc}" -ne 0 ]; then
		echo "### Error ###"
		return 1
	fi
	
	#
	# Do the update of the file after upload using the same file as for create
	#
	url="http://${HOST}:${PORT}/api/v1/files/${id}"
	echo "Running: curl -X PATCH -H ${HEADER} -d@${DATA} ${url}"
	reply=$(curl -X PATCH -H "${HEADER}" -d@"${DATA}" "${url}" 2>/dev/null)
	rc=$?
	if [ "${rc}" -ne 0 ]; then
		echo "### Error ###"
		return 1
	fi
	msg=$(echo "${reply}" | jq)
	rfn="testresults/results-create-upload-update-read-${TS}-step-2-update.lst"
	echo "Result: ${rc} -> ${msg}" | tee "${rfn}"
	# Check the results of the update
	status=$(echo "${msg}" | jq '.status')
	if [ "${status}" != "true" ]; then
		echo "### Error ###"
		return 1
	fi

	#
	# Read by id
	#
	url="http://${HOST}:${PORT}/api/v1/files/${id}"
	echo "Running: curl -X GET -H ${HEADER} ${url}"
	reply=$(curl -X GET -H "${HEADER}" "${url}" 2>/dev/null)
	rc=$?
	if [ "${rc}" -ne 0 ]; then
		echo "### Error ###"
		return 1
	fi
	msg=$(echo "${reply}" | jq)
	rfn="testresults/results-create-upload-update-read-${TS}-step-3-read-id.lst"
	echo "Result: $rc -> $msg"  | tee "${rfn}"
	# Check the results of the update
	status=$(echo "${msg}" | jq '.status')
	if [ "${status}" != "true" ]; then
		echo "### Error ###"
		return 1
	fi

	#
	# Read by filter on attributes like TENANT and DEVICE
	#
	url="http://${HOST}:${PORT}/api/v1/files?tenant=${TENANT}&device=${DEVICE}"
	echo "Running: curl -X GET -H ${HEADER} ${url}"
	reply=$(curl -X GET -H "${HEADER}" "${url}" 2>/dev/null)
	rc=$?
	if [ "${rc}" -ne 0 ]; then
		echo "### Error ###"
		return 1
	fi
	msg=$(echo "${reply}" | jq)
	rfn="testresults/results-create-upload-update-read-${TS}-step-4-read-filter.lst"
	echo "Result: $rc -> $msg"  | tee "${rfn}"
	# Check the results of the update
	status=$(echo "${msg}" | jq '.status')
	if [ "${status}" != "true" ]; then
		echo "### Error ###"
		return 1
	fi

	echo "Success"
	
	return 0
}

n=${1:-1}

while true
do
	if test "$n" -gt 0
	then
		run_test "$n"
		rc=$?
		if [ "${rc}" -ne 0 ]; then
			echo "### Stop ###"
		fi
	else
		break
	fi
	n=$(( n - 1 ))
done	

exit 0
