# Write Gerber RS274X files (for PCBs) in Go

This is an experimental package used to write Gerber (RS274X) files
and bundle them into a ZIP file to send to printed circuit board (PCB)
manufacturers, all from Go code.

## Example

```bash
pushd examples/bifilar-coil
go run main.go
gerbview bifilar-coil.g* bifilar-coil.xln
```

generates:

![bifilar-coil](images/bifilar-coil.png)

and this shows the font and top/bottom silkscreen support:

![with silkscreen](images/go-gerber-with-silkscreen.png)

## Docs

[![GoDoc](https://godoc.org/github.com/gmlewis/go-gerber/gerber?status.svg)](https://godoc.org/github.com/gmlewis/go-gerber/gerber)

----------------------------------------------------------------------

## Webfonts

The SVG webfonts were copied from the excellent
[github.com/boldport/pcbmode](https://github.com/boldport/pcbmode/tree/master/pcbmode/fonts)
project and tweaked slightly.

Please note that these fonts are a work-in-progress and still have rendering
errors that need fixing.

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
