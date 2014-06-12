#!/bin/bash

. ./functions

curl $debug -D /tmp/test.header -k -s -o /tmp/test.body \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1444" \
	"$base_url/sample.avi"

valid_head 'ff3b9222a0d2fb43ce1d80352d2bcb5b7e719b81';
valid_body '56c8ca5e53defe7bac756bffb9893099a95a67fd';
