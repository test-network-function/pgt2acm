#!/bin/bash
rm -rf test/acmgen-output
mkdir test/acmgen-output
cp -r test/init-source-crs test/acmgen-output/source-crs
./pgt2acm -i test/pgt-input -o test/acmgen-output -s test/schema.json -k PtpConfig

if diff -r test/acmgen-output test/acmgen-expected-output; then
	echo "Test Passed"
else
	echo "Test failed"
fi
