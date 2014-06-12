#!/bin/bash

. ./functions

curl $debug -D /tmp/test-unsorted.header -l -k -s -o /tmp/test-unsorted.body \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1443" \
	--data "hello=world" \
	"$base_url/"

# the order of the files returned is based on the inode order and may change in each system
# this works around that
sort --dictionary /tmp/test-unsorted.body > /tmp/test.body
sort --dictionary /tmp/test-unsorted.header > /tmp/test.header

#valid_head '62bde673bbe165824839e10364e759b36af6be0b';
valid_body '3f71d544d4b153fc80474a9724b0cddd04f1b971';

echo -n "hello=world" > "/tmp/postdat2"
echo -n "Data Receive: "
check_same "/tmp/postdat" "/tmp/postdat2"
