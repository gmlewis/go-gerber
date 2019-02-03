#!/bin/bash -ex
go run main.go
gerbview \
  quad-bifilar-coil.g2l \
  quad-bifilar-coil.g3l	\
  quad-bifilar-coil.gbl	\
  quad-bifilar-coil.gbs	\
  quad-bifilar-coil.gko	\
  quad-bifilar-coil.gtl	\
  quad-bifilar-coil.gto	\
  quad-bifilar-coil.gts \
  quad-bifilar-coil.xln
