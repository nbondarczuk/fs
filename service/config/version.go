package config

import (
	"fmt"
	"os"
)

var (
	GitCommitHash, BuiltAt, BuiltBy, BuiltOn string
)

func setVersionInfo(gitCommitHash, builtAt, builtBy, builtOn string) {
	GitCommitHash, BuiltAt, BuiltBy, BuiltOn = gitCommitHash, builtAt, builtBy, builtOn
}

// Log all version relevant info collected by the loader
func logVersionInfo() {
	// version info
	Log.Info("Version info")
	Log.Info("GitCommitHash: " + GitCommitHash)
	Log.Info("      BuiltAt: " + BuiltAt)
	Log.Info("      BuiltBy: " + BuiltBy)
	Log.Info("      BuiltOn: " + BuiltOn)
}

// OneLineInfo composes version infor string using Version/Branch/Revision
// variables compiled into the binary executable
func OneLineVersionInfo() string {
	return fmt.Sprintf("%s_%s_%s_%s", GitCommitHash, BuiltAt, BuiltBy, BuiltOn)
}

func printVersionInfoAndExit() {
	fmt.Printf("%s\n", OneLineVersionInfo())
	os.Exit(0)
}
