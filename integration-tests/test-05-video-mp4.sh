#!/bin/bash

. ./functions

curl $debug -D /tmp/test.header -k -s -o /tmp/test.body \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1444" \
	"$base_url/sample.mp4"

valid_head '8f89aa1928f5ea18968c8242516469e5f7ad7828';
valid_body '94554f1f8edc402582faffab00ce620cc34809f8';
