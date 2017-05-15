package srcpool

import (
	"time"
)

// nowFunc returns the current time; it's overridden in tests.
var nowFunc = time.Now
