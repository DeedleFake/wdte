package debug

import (
	"errors"
	"runtime"
	"runtime/debug"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
)

var (
	ErrNoBuildInfo = errors.New("no build info read")
	ErrDepNotFound = errors.New("WDTE dependency not found")
)

// Version is a WDTE function with the following signature:
//
//    version
//
// It returns the current version of WDTE, as determined by Go's
// module system. If reading build info fails, ErrNoBuildInfo is
// returned. If the build info is read successfully but the version
// couldn't be determined, ErrDepNotFound is returned.
func Version(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("version")

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return wdte.Error{
			Err:   ErrNoBuildInfo,
			Frame: frame,
		}
	}

	for _, dep := range info.Deps {
		if dep.Path != "github.com/DeedleFake/wdte" {
			continue
		}

		return wdte.String(dep.Version)
	}

	return wdte.Error{
		Err:   ErrDepNotFound,
		Frame: frame,
	}
}

var Scope = wdte.S().Map(map[wdte.ID]wdte.Func{
	"version":   wdte.GoFunc(Version),
	"goVersion": wdte.String(runtime.Version()),
	"race":      wdte.Bool(RaceEnabled),
})

func init() {
	std.Register("debug", Scope)
}
