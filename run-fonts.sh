#!/bin/bash -ex
go run examples/test-fonts/main.go -msg "ABCDEFGHIJKLM
NOPQRSTUVWXYZ
abcdefghijklm
nopqrstuvwxyz
0123456789
~\`!@#$%^&*()-_=+,.
{}|[]\\;:'\"<>/?"
rm *.zip

for i in *.gto ; do echo $i ; gerbview $i ; done
