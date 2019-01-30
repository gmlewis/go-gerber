# italianno_regular

![italianno_regular](italianno_regular.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-gerber/gerber"
	_ "github.com/gmlewis/go-gerber/gerber/fonts/italianno_regular"
)

func main() {
	// ...
	Text(x, y, 1.0, message, "italianno_regular", pts),
	// ...
}
```