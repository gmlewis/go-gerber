# veggieburger_bold

![veggieburger_bold](veggieburger_bold.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-gerber/gerber"
	_ "github.com/gmlewis/go-gerber/gerber/fonts/veggieburger_bold"
)

func main() {
	// ...
	Text(x, y, 1.0, message, "veggieburger_bold", pts),
	// ...
}
```