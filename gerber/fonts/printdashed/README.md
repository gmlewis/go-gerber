# printdashed

![printdashed](printdashed.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-gerber/gerber"
	_ "github.com/gmlewis/go-gerber/gerber/fonts/printdashed"
)

func main() {
	// ...
	Text(x, y, 1.0, message, "printdashed", pts),
	// ...
}
```
