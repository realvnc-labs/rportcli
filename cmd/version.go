package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	Version   = "0.0.1"
	BuildTime = ""
	GitCommit = ""
	GitRef    = ""
)

func init() {
	rootCmd.AddCommand(versionCmd)
	BuildTime = time.Now().String()
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number of rportcli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version())
	},
}

func version() string {
	return fmt.Sprintf(`Version: %s
BuildDate  %s
GoVersion  %s
GitBranch  %s
GitCommit  %s
GoArch     %s
GoOS       %s
`, Version, BuildTime, runtime.Version(), GitRef, GitCommit, runtime.GOARCH, runtime.GOOS)
}
