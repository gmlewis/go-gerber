#!/bin/bash -ex
go run main.go -n $@
gerbview \
  bifilar-coil.gbl \
  bifilar-coil.gbs \
  bifilar-coil.gko \
  bifilar-coil.gtl \
  bifilar-coil.gto \
  bifilar-coil.gts \
  bifilar-coil.xln
