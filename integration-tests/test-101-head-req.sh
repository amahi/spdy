#!/bin/bash

. ./functions

curl $debug -I -D /tmp/test-unsorted.header -k -s -o /tmp/test-unsorted.body \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1443" \
	"$base_url/"
check_same '/tmp/test-unsorted.header' '/tmp/test-unsorted.body';
