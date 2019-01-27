#!/bin/bash -ex
go run examples/test-fonts/main.go -all
rm *.zip

for i in *.gto ; do echo $i ; gerbview $i ; done
