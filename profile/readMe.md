# Profile

Profile package is used to perform memory and CPU profiling using [Pyroscope](https://github.com/grafana/pyroscope). Once enabled, the application pushes periodic metric to a pyroscope server.

## Usage
Set the following mandatory env variables to enable profiling:
```
# Mandatory
PYROSCOPE_HOST="host:port"
APP_NAME="app_name"

# Optional
RELEASE_VERSION="v1.2.3"
```

Add the above mandatory env variables and call InitPyroscope() in main.go to enable profiling with default profileTypes.

Use InitPyroscopeWithProfiles() to enable profiling with specific profileTypes.

## Features

    // these profile types are enabled by default:
    pyroscope.ProfileCPU,
    pyroscope.ProfileAllocObjects,
    pyroscope.ProfileAllocSpace,
    pyroscope.ProfileInuseObjects,
    pyroscope.ProfileInuseSpace,

    // these profile types are optional and can used with InitPyroscopeWithProfiles():
    pyroscope.ProfileGoroutines,
    pyroscope.ProfileMutexCount,
    pyroscope.ProfileMutexDuration,
    pyroscope.ProfileBlockCount,
    pyroscope.ProfileBlockDuration,
