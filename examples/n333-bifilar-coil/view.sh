#!/bin/bash -ex
go run main.go -step 0.1 -n 20
gerbview \
  n333-bifilar-coil.g2l \
  n333-bifilar-coil.g3l \
  n333-bifilar-coil.g4l \
  n333-bifilar-coil.g5l \
  n333-bifilar-coil.gbl \
  n333-bifilar-coil.gbs \
  n333-bifilar-coil.gko \
  n333-bifilar-coil.gtl \
  n333-bifilar-coil.gto \
  n333-bifilar-coil.gts \
  n333-bifilar-coil.xln
