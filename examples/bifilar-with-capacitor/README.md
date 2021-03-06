# Single (2-layer) bifilar coil

This is a design of a bifilar coil, meaning that when this design is
fabricated on a two-layer printed circuit board (PCB), there will be a
single bifilar coil.

This designs differs from the others in that a single coil is devoted
to one layer and the board itself is the dielectric between the top and
bottom coils.

This design requires that an external wire connect the inner terminal
of one coil to the outer terminal of the other coil (with the added
benefit of being able to insert a tuning capacitor in between the
two coils.)

The theory is that when you drive the coil with a sine wave at its
resonant frequency, it can transfer power at its greatest efficiency.

The beauty of this design is that it is parametric.

## How it is wired

Both coils are concentric and wind in the same direction.
Therefore the magnetic field from each coil section combines uniformly
with the other coil resulting in a stronger, cohesive field.

Additionally, this design differs in that the two coils are directly
sandwiched together via the board's dielectric such that the capacitance
is stronger between the coils.

This shows the various layers on a small (100mm x 100mm) coil to highlight
the wiring and various layers of the PCB:

![bifilar-with-capacitor-layers](bifilar-with-capacitor-layers.gif)

## Parametric design

In this design, coils can be created with varying board size,
trace widths, and gaps between traces. As a result, this
parametric design could theoretically be used for coils of any
manufacturable size (from microscopic on up).

## Example coils

In this section, we will document emperically-measured resistances
and resonant frequencies of fabricated coils as they become
available.

### Coil 1

Here is a large (400mm x 400mm) PCB with 92 windings per coil:

![big-bifilar-cap-layers](big-bifilar-cap-layers.gif)

and measurements on the coil:

![20190610_154844.jpg](20190610_154844.jpg)
![20190610_153519.jpg](20190610_153519.jpg)
![20190610_151523.jpg](20190610_151523.jpg)
![20190610_155238.jpg](20190610_155238.jpg)

It has 2.0mm trace width, 2.0mm pad diameter, 1.0mm drill hole diameter, and 0.15mm trace gap.
Its board thickness is 0.6mm.

Please visit [Master Ivo's YouTube channel](https://www.youtube.com/channel/UCVZQj8Cl06UDheK3B4N7d9w)
to see some novel uses of bifilar coils.

### Coil 2

Here is a small (100mm x 100mm) PCB with 21 windings per coil:

It has 2.0mm trace width, 2.0mm pad diameter, 1.0mm drill hole diameter, and 0.15mm trace gap.

It measures 2.8Ω, 0.07mH, and 222pF. Its board thickness is 1.6mm.

![IMG_20190609_190535.jpg](IMG_20190609_190535.jpg)

![IMG_20190609_192140.jpg](IMG_20190609_192140.jpg)

### Coil 3

Coil 3 is a large (400mm x 400mm) PCB with 45 windings per coil.

It has 4.0mm trace width, 4.0mm pad diameter, 1.5mm drill hole diameter, and 0.15mm trace gap.

It measures 9.1Ω, 1.23mH, and 7460pF. Its board thickness is 0.8mm.

![IMG_20190624_183659.jpg](IMG_20190624_183659.jpg)

![IMG_20190624_183713.jpg](IMG_20190624_183713.jpg)

![IMG_20190624_183949.jpg](IMG_20190624_183949.jpg)

![IMG_20190624_185214.jpg](IMG_20190624_185214.jpg)

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
