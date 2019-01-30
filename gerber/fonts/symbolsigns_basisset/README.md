# symbolsigns_basisset

![symbolsigns_basisset](symbolsigns_basisset.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-gerber/gerber"
	_ "github.com/gmlewis/go-gerber/gerber/fonts/symbolsigns_basisset"
)

func main() {
	// ...
	Text(x, y, 1.0, message, "symbolsigns_basisset", pts),
	// ...
}
```
