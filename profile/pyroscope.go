package profile

import (
	"github.com/pkg/errors"
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/skit-ai/vcore/env"
)

var (
	host, appName, releaseVersion string
)

func init() {
	host = env.String("PYROSCOPE_HOST", "")
	appName = env.String("APP_NAME", "")
	releaseVersion = env.String("RELEASE_VERSION", "")
}

var (
	errNoHost = errors.New("host not defined")
)

// InitPyroscope initialises profiling using Pyroscope.
func InitPyroscope() error {
	if host == "" {
		return errNoHost
	}

	// Initialize Pyroscope with default profileTypes
	err := initPyroscope(nil)
	if err != nil {
		return err
	}

	return nil
}

// InitPyroscopeWithProfiles initialises profiling wit specified ProfileTypes.
// List of profiles can be found here: https://github.com/grafana/pyroscope-golang/blob/main/pyroscope/types.go
func InitPyroscopeWithProfiles(profileTypes []pyroscope.ProfileType) error {
	if host == "" {
		return errNoHost
	}

	return initPyroscope(profileTypes)
}

func initPyroscope(profileTypes []pyroscope.ProfileType) error {
	_, err := pyroscope.Start(
		pyroscope.Config{
			ApplicationName: appName,
			// Pyroscope host to push metrics to
			ServerAddress: host,
			// Set application release version
			Tags: map[string]string{
				"app_version": releaseVersion,
			},
			Logger:       pyroscope.StandardLogger,
			ProfileTypes: profileTypes,
		})

	return err
}
