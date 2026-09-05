package version

// Build information populated at build time via linker flags.
var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

// Info is the build metadata exposed by operational endpoints and logs.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
}

// String returns the short version string.
func String() string {
	return Version
}

// Full returns the complete version with commit and build time.
func Full() string {
	return Version + " (" + Commit + " @ " + BuildTime + ")"
}

// Current returns the build metadata compiled into this binary.
func Current() Info {
	return Info{Version: Version, Commit: Commit, BuildTime: BuildTime}
}
