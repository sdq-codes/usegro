package version

const (
	// AppSemVer is apps version.
	AppSemVer = "0.0.2"
)

var (
	// GitCommit is the current HEAD set using ldflags.
	GitCommit string

	// Version is the built version.
	Version string = AppSemVer

	// BuildTime is the build timestamp set using ldflags.
	BuildTime string

	// BuildUser is the user who built the binary set using ldflags.
	BuildUser string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}
