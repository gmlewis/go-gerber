#!/bin/bash -ex
go run main.go -step 0.1 -n 20
gerbview \
  icosi-bifilar-coil.g2l \
  icosi-bifilar-coil.g3l \
  icosi-bifilar-coil.g4l \
  icosi-bifilar-coil.g5l \
  icosi-bifilar-coil.gbl \
  icosi-bifilar-coil.gbs \
  icosi-bifilar-coil.gko \
  icosi-bifilar-coil.gtl \
  icosi-bifilar-coil.gto \
  icosi-bifilar-coil.gts \
  icosi-bifilar-coil.xln
