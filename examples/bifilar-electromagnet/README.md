# Bifilar Electromagnet

Nikola Tesla patented the bifilar coil in [U.S. Patent 512340](
https://teslauniverse.com/nikola-tesla/patents/us-patent-512340-coil-electro-magnets).

The point of the bifilar coil is to add internal capacitance to the coil
between the windings and maximize the voltage between every pair of adjacent
wires. I was thinking that if a coil could be wound in such a manner, why
couldn't an electromagnet?

## Electromagnet 1 - How it is wired

This electromagnet is wound with two wires. They remain separate until the
very end. First, `wire 1` is used to wind `coil 1`. Then `wire 2` is used
to wind `coil 2` directly on top of `coil 1`. Then `wire 1` is wrapped back
around (through the center air gap) to the front of the electromagnet spool
holder and is used to wind `coil 3`. Then `wire 2` is wrapped back around
and used to wind `coil 4`. They alternate back and forth until 14 total
coils are wrapped.

Each coil has 120 turns.

Finally, the end of `wire 1` (coils 1, 3, 5, 7, 9, 11, and 13) can be
connected to the start of `wire 2` (coils 2, 4, 6, 8, 10, 12, and 14)
and the end of `wire 2` can return to the circuit.

Optionally, instead of connecting the end of `wire 1` to the start of `wire 2`,
a capacitor (or variable capacitor) can be placed in between (in series)
in order to tune the capacitance to ideally make the resonant frequency
of the coil the same as its main oscillating frequency (in the case of
a brushless motor, for example).

Here is a diagram of the first winding:

![winding 1 diagram](coil1-winding-120turns-6920mm.png)

The second winding is kept separate from the first so that an external
capacitor could be added between the two windings. It is wound similarly
to the first coil, in the same direction:

![winding 2 diagram](coil2-winding-120turns-7664mm.png)

The third winding is a continuation of the wire from the first winding.
Here is a view from the front of the coil:

![winding 3 front](coil3-front.png)

Here is a view from the back of the coil:

![winding 3 rear](coil3-rear.png)

### Winding lengths (theoretical)

| Winding | Length (mm) |
|  :---:  |   :---:     |
|    1    |    6410     |
|    2    |    7164     |
|    3    |    7918     |
|    4    |    8672     |
|    5    |    9426     |
|    6    |   10179     |
|    7    |   10933     |
|    8    |   11687     |
|    9    |   12441     |
|   10    |   13195     |
|   11    |   13949     |
|   12    |   14703     |
|   13    |   15834     |
|   14    |   16211     |
|  Total  |  158722     |

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
