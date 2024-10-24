package buildinfo

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/piplabs/story/lib/log"
)

const (
	VersionMajor = 0          // Major version component of the current release
	VersionMinor = 12         // Minor version component of the current release
	VersionPatch = 2          // Patch version component of the current release
	VersionMeta  = "unstable" // Version metadata to append to the version string
)

// Version returns the version of the whole story-monorepo and all binaries built from this git commit.
func Version() string {
	return fmt.Sprintf("v%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}

// VersionWithMeta holds the textual version string including the metadata.
func VersionWithMeta() string {
	v := Version()
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}

	return v
}

// Instrument logs the version, git commit hash, and timestamp from the runtime build info.
// It also sets metrics.
func Instrument(ctx context.Context) {
	commit, timestamp := get()

	version := Version()
	versionWithMeta := VersionWithMeta()

	log.Info(ctx, "Version info",
		"version", versionWithMeta,
		"git_commit", commit,
		"git_timestamp", timestamp,
	)

	versionGauge.WithLabelValues(version).Set(1)
	commitGauge.WithLabelValues(commit).Set(1)

	ts, _ := time.Parse(time.RFC3339, timestamp)
	timestampGauge.Set(float64(ts.Unix()))
}

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information of this binary",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			commit, timestamp := get()

			var sb strings.Builder
			_, _ = sb.WriteString("Version       " + VersionWithMeta())
			_, _ = sb.WriteString("\n")
			_, _ = sb.WriteString("Git Commit    " + commit)
			_, _ = sb.WriteString("\n")
			_, _ = sb.WriteString("Git Timestamp " + timestamp)
			_, _ = sb.WriteString("\n")

			cmd.Print(sb.String())
		},
	}
}

// get returns the git commit hash and timestamp from the runtime build info.
func get() (hash string, timestamp string) { //nolint:nonamedreturns // Disambiguate identical return types.
	hash, timestamp = "unknown", "unknown"
	hashLen := 7

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return hash, timestamp
	}

	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			if len(s.Value) < hashLen {
				hashLen = len(s.Value)
			}
			hash = s.Value[:hashLen]
		} else if s.Key == "vcs.time" {
			timestamp = s.Value
		}
	}

	return hash, timestamp
}
