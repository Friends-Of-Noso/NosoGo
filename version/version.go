package version

import "fmt"

var (
	VersionMajor = 0
	VersionMinor = 0
	VersionPatch = 2
	Version      = "v0.0.2"
	Name         = "nosogo"
	// GitCommit is set with --ldflags "-X main.gitCommit=$(git rev-parse HEAD)"
	GitCommit string
	Title     = fmt.Sprintf("%s %s", Name, Version)
)

func init() {
	if GitCommit != "" {
		Version += "+" + GitCommit[:8]
	}
}
