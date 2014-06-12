#!/bin/bash

. ./functions

curl $debug -D /tmp/test.header -k -s -o /tmp/test.body \
	-H "Range: bytes=512-2047" \
	-H "User-Agent: VLC/2.1.0 LibVLC/2.1.0" \
	-H "Icy-Metadata: 1" \
	-H "Host: localhost:1443" \
	"$base_url/sample.avi"

valid_head 'c792f24861c2bfdeac072c8a77e63aa4e8003da3';
valid_body '2924ddb33820391a3f760323c0c8a5010260b12d';
