#!/bin/bash

. ./functions

curl $debug -D /tmp/test.header -k -s -o /tmp/test.body \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1444" \
	"$base_url/sample.mkv"

valid_head '38c2d97c01e00d585a1a3426dbfa80ff39fb3d8e';
valid_body '48c9a4e3a24324e3d3ac6b284f7502878ace1909';
