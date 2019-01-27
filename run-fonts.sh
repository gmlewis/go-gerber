#!/bin/bash -ex
go run examples/test-fonts/main.go -msg "ABCDEF
GHIJKL
MNOPQR
STUVWX
YZ
abcdef
ghijkl
mnopqr
stuvwx
yz
01234
56789"
rm *.zip
gerbview aaarghnormal.gto
gerbview helsinkiregular.gto
gerbview ubuntumonoregular.gto
gerbview fascinate_inlineregular.gto
gerbview latoregular.gto
gerbview webfontb8enx5qc.gto
gerbview gooddogregular.gto
gerbview overlockregular.gto
gerbview webfontrtffmgjf.gto
