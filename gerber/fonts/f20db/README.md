# f20db

![f20db](f20db.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-gerber/gerber"
	_ "github.com/gmlewis/go-gerber/gerber/fonts/f20db"
)

func main() {
	// ...
	Text(x, y, 1.0, message, "f20db", pts),
	// ...
}
```
