package cmd

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/spf13/cobra"
)

// Version is the Version of the cli to be overwritten by goreleaser in the CI run with the Version of the release in github
var Version string

func getVersion() string {
	noVersionAvailable := "No version info available for this build, run 'innoctl help version' for additional info"

	if len(Version) != 0 {
		return Version
	}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return noVersionAvailable
	}

	// If no main version is available, Go defaults it to (devel)
	if bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}

	var vcsRevision string
	var vcsTime time.Time
	for _, setting := range bi.Settings {
		switch setting.Key {
		case "vcs.revision":
			vcsRevision = setting.Value
		case "vcs.time":
			vcsTime, _ = time.Parse(time.RFC3339, setting.Value)
		}
	}

	if vcsRevision != "" {
		return fmt.Sprintf("%s, (%s)", vcsRevision, vcsTime)
	}

	return noVersionAvailable
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display application version information.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		version := getVersion()
		fmt.Printf("Inno CLI version: %v\n", version)
	},
}
