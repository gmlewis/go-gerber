# Quad (4-layer) bifilar coil

This is a design of a quad bifilar coil, meaning that when this design is
fabricated on a four-layer printed circuit board (PCB), there will be four
bifilar coils wired in series, one bifilar coil per layer.

The theory is that when you drive the coil with a sine wave at its
resonant frequency, it can transfer power at its greatest efficiency.

The beauty of this design is that there are multiple tap points along
the length of the coils that result in different resonant frequencies
due to the varying reactance. Thus, if your electronics drive the
coils properly from the different tap points, this single coil design
can potentially support the optimal energy transfer of multiple
frequencies.

Another experiment to try is to ground the second tap from the end,
drive the end tap with a sine wave at the resonant frequency of the
remainder of the coil, and then get the amplified benefits of the
remainder as in a secondary coil of a transformer.

## How it is wired

All coils are concentric and wind in the same direction.
Therefore the magnetic field from each coil section combines uniformly
with the other coils resulting in a stronger, cohesive field.

Here is a diagram showing how it is wired:

![quad-bifilar-coil-diagram](quad-bifilar-coil-diagram.png)

This shows the various layers on a small (n=20) coil to highlight
the wiring and various layers of the PCB:

![quad-bifilar-coil-layers](quad-bifilar-coil-layers.gif)

## Parametric design

In this design, coils can be created with varying trace widths, gaps
between traces, and number of spirals per coil. As a result, this
parametric design could theoretically be used for coils of any
manufacturable size (from microscopic on up).

## Example coils

In this section, we will document emperically-measured resistances
and resonant frequencies of fabricated coils as they become
available.

### Coil 1

![quad-bifilar-coil-1.jpg](/images/quad-bifilar-coil-1.jpg)

![IMG_20190608_193223.jpg](IMG_20190608_193223.jpg)

Coil 1 consists of:

```
PCB thickness = 1.6mm
trace width = 0.15mm (6 mils)
gap width = 0.15mm
number of coils per spiral = 100
```

| From point | To point | DC resistance (Ω) | Resonant Frequency |
|   :---:    |  :---:   |      :---:        |       :---:        |
|    TR      |  BR/TL   |      105.3 ?      |                    |
|   BR/TL    |  BL/3L   |      105.4 ?      |                    |
|   BL/3L    |  2R/3R   |      441          |                    |
|   2R/3R    |   2L     |      442          |                    |
|    TR      |   2L     |      989.7        |       ~25kHz       |


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
