package version

// Version information injected by GoReleaser
var (
	Version = "dev"
	Build   = "dev"
	Commit  = "none"
	Branch  = "main"
)

// VersionInfo holds the version information
type VersionInfo struct {
	Version string
	Build   string
	Commit  string
	Branch  string
}

// GetVersionInfo returns the current version information
func GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version: Version,
		Build:   Build,
		Commit:  Commit,
		Branch:  Branch,
	}
}