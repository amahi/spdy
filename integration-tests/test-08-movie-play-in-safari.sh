#!/bin/bash

. ./functions

curl $debug -D /tmp/test.header -k -s -o /tmp/test.body \
	-H "User-Agent: AppleCoreMedia/1.0.0.11G63 (Macintosh; U; Intel Mac OS X 10_7_5; en_us)" \
	-H "Accept: */*" \
	-H "Range: bytes=4076148-4092531" \
	-H "Host: localhost:1444" \
	"$base_url/sample.avi"

valid_head '12beb434ca7492c296b0c203890e996c280bc712';
valid_body 'e7e3dc7b645f90ae6145ecf8e5bdd0e6c9478d92';

