#!/bin/bash -ex
go run main.go -n $@
gerbview \
  hex-bifilar-coil.g2l \
  hex-bifilar-coil.g3l \
  hex-bifilar-coil.g4l \
  hex-bifilar-coil.g5l \
  hex-bifilar-coil.gbl \
  hex-bifilar-coil.gbs \
  hex-bifilar-coil.gko \
  hex-bifilar-coil.gtl \
  hex-bifilar-coil.gto \
  hex-bifilar-coil.gts \
  hex-bifilar-coil.xln
