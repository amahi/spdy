#!/bin/bash

. ./functions

curl $debug -D /tmp/test.header -k -s -o /tmp/test.body \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1443" \
	-H "Dnt: 1" \
	-H "Cache-Control: max-age=0" \
	-H "If-Modified-Since: Fri, 31 Oct 2013 23:08:54 GMT" \
	-H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; rv:24.0) Gecko/20100101 Firefox/24.0" \
	-H "Accept-Language: en-US,en;q=0.5" \
	-H "Cookie: locale=en" \
	-H "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8" \
	"$base_url/"

valid_head 'ec412aa8662157b16b6cf56247005c1d77205d3e';
