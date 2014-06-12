#!/bin/bash

. ./functions

curl $debug -D /tmp/test.header -k -s -o /tmp/test.body \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1444" \
	"$base_url/image.jpg"

valid_head '0909c9cefd359d1cbda8f71a477abfcd6e75ef52';
valid_body 'a37afec0825c483a906f32adbb70528b6d5867b4';
