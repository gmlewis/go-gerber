# fontawesome5freesolid

![fontawesome5freesolid](fontawesome5freesolid.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-gerber/gerber"
	_ "github.com/gmlewis/go-gerber/gerber/fonts/fontawesome5freesolid"
)

func main() {
	// ...
	Text(x, y, 1.0, message, "fontawesome5freesolid", pts),
	// ...
}
```