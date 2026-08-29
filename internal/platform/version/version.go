package version

// Build information populated at build time via linker flags.
var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

// String returns the short version string.
func String() string {
	return Version
}

// Full returns the complete version with commit and build time.
func Full() string {
	return Version + " (" + Commit + " @ " + BuildTime + ")"
}
