# lobstertwo_bold

![lobstertwo_bold](lobstertwo_bold.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-gerber/gerber"
	_ "github.com/gmlewis/go-gerber/gerber/fonts/lobstertwo_bold"
)

func main() {
	// ...
	Text(x, y, 1.0, message, "lobstertwo_bold", pts),
	// ...
}
```