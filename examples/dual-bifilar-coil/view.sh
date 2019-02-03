#!/bin/bash -ex
go run main.go
gerbview \
  dual-bifilar-coil.gbl \
  dual-bifilar-coil.gbs \
  dual-bifilar-coil.gko \
  dual-bifilar-coil.gtl \
  dual-bifilar-coil.gto \
  dual-bifilar-coil.gts \
  dual-bifilar-coil.xln
