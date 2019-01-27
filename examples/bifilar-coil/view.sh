#!/bin/bash -ex

for i in $(cat <<EOF
aaarghnormal
fascinate_inlineregular
gooddogregular
helsinkiregular
latoregular
overlockregular
pacifico
snickles
ubuntumonoregular
EOF
) ; do echo $i
go run main.go -font $i
mv bifilar-coil.zip bifilar-coil.${i}.zip
# gerbview bifilar-coil.gbl bifilar-coil.gbs bifilar-coil.gtl bifilar-coil.gts bifilar-coil.gbo bifilar-coil.gko bifilar-coil.gto bifilar-coil.xln
done
