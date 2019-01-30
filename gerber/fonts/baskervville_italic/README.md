# baskervville_italic

![baskervville_italic](baskervville_italic.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-gerber/gerber"
	_ "github.com/gmlewis/go-gerber/gerber/fonts/baskervville_italic"
)

func main() {
	// ...
	Text(x, y, 1.0, message, "baskervville_italic", pts),
	// ...
}
```
