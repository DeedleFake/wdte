// Package all is a convience package that imports the entire standard
// library, thus registering it with std.Import.
package all

import (
	_ "github.com/DeedleFake/wdte/std/arrays"
	_ "github.com/DeedleFake/wdte/std/io"
	_ "github.com/DeedleFake/wdte/std/io/file"
	_ "github.com/DeedleFake/wdte/std/math"
	_ "github.com/DeedleFake/wdte/std/stream"
	_ "github.com/DeedleFake/wdte/std/strings"
)
