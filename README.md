# Write Gerber RS274X files (for PCBs) in Go

This is an experimental package used to write Gerber (RS274X) files
and bundle them into a ZIP file to send to printed circuit board (PCB)
manufacturers, all from Go code.

## Examples

Please see the examples in the various directories:

* 2-layer single bifilar coil: [bifilar-coil](examples/bifilar-coil)
* 2-layer 2 bifilar coil: [dual-bifilar-coil](examples/dual-bifilar-coil)
* 3-layer 3 bifilar coil: [tri-bifilar-coil](examples/tri-bifilar-coil)
* 4-layer 4 bifilar coil: [quad-bifilar-coil](examples/quad-bifilar-coil)
* 5-layer 5 bifilar coil: [penta-bifilar-coil](examples/penta-bifilar-coil)
* 6-layer 6 bifilar coil: [hex-bifilar-coil](examples/hex-bifilar-coil)
* 20-layer 20 bifilar coil: [icosi-bifilar-coil](examples/icosi-bifilar-coil)

## Documentation
[![GoDoc](https://godoc.org/github.com/gmlewis/go-gerber/gerber?status.svg)](https://godoc.org/github.com/gmlewis/go-gerber/gerber)

----------------------------------------------------------------------

## Webfonts using `go-fonts`

Webfont support has been switched to using
[github.com/gmlewis/go-fonts](https://github.com/gmlewis/go-fonts).

Below are some example fonts but there are many more in the `go-fonts` repo
to choose from.

### AaarghNormal

![aaarghnormal](images/aaarghnormal.png)

### Fascinate_InlineRegular

![fascinate_inlineregular](images/fascinate_inlineregular.png)

### GoodDogRegular

![gooddogregular](images/gooddogregular.png)

### HelsinkiRegular

![helsinkiregular](images/helsinkiregular.png)

### LatoRegular

![latoregular](images/latoregular.png)

### OverlockRegular

![overlockregular](images/overlockregular.png)

### Pacifico

![pacifico](images/pacifico.png)

### Snickles

![snickles](images/snickles.png)

### UbuntuMonoRegular

![ubuntumonoregular](images/ubuntumonoregular.png)

----------------------------------------------------------------------

Enjoy!

----------------------------------------------------------------------

# License

Copyright 2019 Glenn M. Lewis. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
