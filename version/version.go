package version

import "fmt"

var (
	VersionMajor = 0
	VersionMinor = 1
	VersionPatch = 0
	Version      = "v0.1.0"
	Name         = "NosoGo"
	// GitCommit is set with --ldflags "-X main.gitCommit=$(git rev-parse HEAD)"
	GitCommit string
	Title     = fmt.Sprintf("%s %s", Name, Version)
)

func init() {
	if GitCommit != "" {
		Version += "+" + GitCommit[:8]
	}
}
