# abeezee_italic

![abeezee_italic](abeezee_italic.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-gerber/gerber"
	_ "github.com/gmlewis/go-gerber/gerber/fonts/abeezee_italic"
)

func main() {
	// ...
	Text(x, y, 1.0, message, "abeezee_italic", pts),
	// ...
}
```