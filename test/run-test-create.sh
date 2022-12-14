#!/bin/bash

TS=$(date +"%Y%m%d%H%M%S")
HOST=localhost
PORT=1234
HEADER="Content-Type: application/json"
TENANT=xxx
DEVICE=yyy
DATA=

# shellcheck source=test_tools.sh
. test_tools.sh

# Perform one test, sending the request to an API, checki the reply and try upload a file
function run_test() {
	local testno fn bfn status rc reply msg rfn url

	testno="$1"
	
	# Make a test file with random content of size N blocks
	make_test_file "${testno}"
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
	rfn="testresults/results-create-${TS}-step-1-create.ls"
	echo "Result: ${rc} -> ${msg}" | tee "${rfn}"
	# Check the results
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
		if [ ${rc} -ne 0 ]; then
			echo "### Stop ###"
		fi
	else
		break
	fi
	n=$(( n - 1 ))
done	

exit 0
