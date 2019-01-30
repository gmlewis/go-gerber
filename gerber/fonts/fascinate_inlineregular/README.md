# fascinate_inlineregular

![fascinate_inlineregular](fascinate_inlineregular.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-gerber/gerber"
	_ "github.com/gmlewis/go-gerber/gerber/fonts/fascinate_inlineregular"
)

func main() {
	// ...
	Text(x, y, 1.0, message, "fascinate_inlineregular", pts),
	// ...
}
```