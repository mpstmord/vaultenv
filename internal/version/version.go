// Package version provides build-time version metadata for vaultenv.
package version

import "fmt"

// These variables are set at build time via -ldflags.
var (
	// Major is the major version component.
	Major = "0"
	// Minor is the minor version component.
	Minor = "1"
	// Patch is the patch version component.
	Patch = "0"
	// Commit is the short git commit hash.
	Commit = "unknown"
	// BuildDate is the UTC build timestamp.
	BuildDate = "unknown"
)

// Info holds structured version metadata.
type Info struct {
	Major     string
	Minor     string
	Patch     string
	Commit    string
	BuildDate string
}

// Get returns the current build's version Info.
func Get() Info {
	return Info{
		Major:     Major,
		Minor:     Minor,
		Patch:     Patch,
		Commit:    Commit,
		BuildDate: BuildDate,
	}
}

// String returns a human-readable version string, e.g. "v0.1.0 (abc1234, 2024-01-01)".
func (i Info) String() string {
	return fmt.Sprintf("v%s.%s.%s (commit=%s, built=%s)",
		i.Major, i.Minor, i.Patch, i.Commit, i.BuildDate)
}

// Short returns only the semver portion, e.g. "v0.1.0".
func (i Info) Short() string {
	return fmt.Sprintf("v%s.%s.%s", i.Major, i.Minor, i.Patch)
}
