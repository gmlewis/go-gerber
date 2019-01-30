# aileron_heavyitalic

![aileron_heavyitalic](aileron_heavyitalic.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-gerber/gerber"
	_ "github.com/gmlewis/go-gerber/gerber/fonts/aileron_heavyitalic"
)

func main() {
	// ...
	Text(x, y, 1.0, message, "aileron_heavyitalic", pts),
	// ...
}
```