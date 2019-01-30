# printclearly

![printclearly](printclearly.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-gerber/gerber"
	_ "github.com/gmlewis/go-gerber/gerber/fonts/printclearly"
)

func main() {
	// ...
	Text(x, y, 1.0, message, "printclearly", pts),
	// ...
}
```
