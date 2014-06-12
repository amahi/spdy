#!/bin/bash

. ./functions

curl $debug -D /tmp/test.header -k -s -o /tmp/test-unsorted.body \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1443" \
	"$base_url/"

# the order of the files returned is based on the inode order and may change in each system
# this works around that
sort --dictionary /tmp/test-unsorted.body > /tmp/test.body

valid_head '57c7644f36d06a104d40efcd2b925fa408fd2082';
valid_body 'edd16e957becc92b2beed5a42b77cfedeaebf0c1';
