# vcore

Golang Common utility functions

## vcore/errors

The errors package is a wrapper around the brilliant [pkg/errors](https://github.com/pkg/errors) library.
The major difference between [pkg/errors](https://github.com/pkg/errors) and this library is support for the following:

* Wrapping of an error into a custom error so as to add a stacktrace.
* Allows creation of errors with a stack which have no cause(nil cause).
* Every error created supports the `fatality` interface which is meant to inform us if the error is fatal or not.
* Every error created supports the `tagged` interface which returns any tags(`map[string]string`) associated with an error.
* The stacktrace of the error will return the stactrace starting from the deepest cause.

### Basic Usage

#### Create Error without cause:

This error is non-fatal. As marked by the last bool flag.

```go
errors.NewError("Error without a cause", nil, false)
```

#### Create Error with a cause:

```go
cause := errors.NewError("Error without a cause", nil, false)
errWithCause := errors.NewError("Error with a cause", cause, false)
```

#### Check if an error is fatal

```go
cause := errors.NewError("Error without a cause", nil, false)
if cause.Fatal{
    fmt.Println("This error is fatal.")
}
```

#### Get the stacktrace of the error and print it:

```go
cause := errors.NewError("Error without a cause", nil, false)
errWithCause := errors.NewError("Error with a cause", cause, false)
fmt.Println(errWithCause.Stacktrace())
```

Alternatively,

```go
cause := errors.NewError("Error without a cause", nil, false)
errWithCause := errors.NewError("Error with a cause", cause, false)
errWithCause.PrintStackTrace()
```

## vcore/crypto

The crypto module is meant to help services implement various cryptographic functions with ease.
Current features include -

1. Encryption of []byte and string. Supported techniques -
    - AES-256-GCM
2. Decryption of []byte. Supported techniques -
    - AES-256-GCM

AES-256 is PCI DSS compliant, as it is a recognised industry standard encryption.

This module exports the following functions -
1. EncryptBytes: Encrypt a bytearray. Example usage -
``` go
x := []byte("hello world")
enc := EncryptBytes(x)
fmt.Println(enc)
```
2. EncryptString: Encrypt a string. Example usage -
``` go
x := "hello world"
enc := EncryptString(x)
fmt.Println(enc)
```
3. DecryptBytes: Decrypt a bytearray. Example usage -
``` go
y := []byte{180, 27, 0, 28, 249, 65, 157, 217, 78, 134, 227, 25, 135, 180, 197, 2, 170, 235, 128, 7, 99, 202, 202, 210, 149, 75, 209, 157, 114, 129, 236, 206, 62, 132, 175, 42, 26, 224, 26}
p := DecryptBytes(y)
fmt.Println(string(p))
```

### Key management

Vault is used to generate the encrypted data key when an environment/client is set up. The encrypted data key is passed to vcore as an environment variable.
Vcore then calls Vault APIs to decrypt the data key and proceed with the encryption/decryption.

### Environment Variables needed

The following environment variables are needed to utilize the crypto module -
``` sh
export VAULT_URI="http://localhost:8200"
export VAULT_ROLE_ID="****"
export VAULT_SECRET_ID="****"
export VAULT_APPROLE_MOUNTPATH="approle"
export ENCRYPTED_DATA_KEY="****"
export VAULT_DATA_KEY_NAME="datakey-name"
```

Note: the above environment variables are just examples, set up vault and replace the actual values above.

## vcore/log

The log package is a basic wrapper on the standard log package  in Go's stdlib.
The log package logs to STDOUT and supports log levels(`int`).

The log levels currently supported are:

* 0 - ERROR
* 1 - WARN
* 2 - INFO
* 3 - DEBUG
* 4 - TRACE

To log using this package, one needs to use an instance of the `log.Logger` struct.

The `log.Logger` struct supports the following methods which can be used to log messages:

* `log.Trace(args ...interface{})`
* `log.Tracef(format string, args ...interface{})`
* `log.Debug(args ...interface{})`
* `log.Debugf(format string, args ...interface{})`
* `log.Info(args ...interface{})`
* `log.Infof(format string, args ...interface{})`
* `log.Warn(args ...interface{})`
* `log.Warnf(format string, args ...interface{})`
* `log.Error(err error, args ...interface{})`
* `log.Errorf(err error, format string, args ...interface{})`

Each of these methods are wrappers that correspond to a log level. This enforces the user to take cognizance of the log 
level of whatever they are attempting to log.

To set the log level on a `log.Logger` struct, make use of the `log.SetLevel(level int)` function.

### Default Logger

To quickly start logging messages, make use of the default logger(default level `WARN`). This can be done by simply 
calling the functions stated above.

Eg. To add a trace log

```go
headers := make(map[string]string)
for k := range req.Header {
    headers[strings.ToLower(k)] = req.Header.Get(k)
}
log.SetLevel(log.DEBUG)
log.Tracef("Headers: %s", headers)
```

Here, we directly make use of the default logger. Please note, since the log level is set to DEBUG here, this trace message will not be logged.

### Custom Logger

```go
customLogger := log.Logger{log.DEBUG}
customLogger.Debug("This is a debug message")
```

## vcore/events

### Sending Cost Tracker Event

First export `AWS_CREDENTIALS`

```sh
export AWS_ACCESS_KEY_ID=AKID
export AWS_SECRET_ACCESS_KEY=SECRET
export AWS_REGION=us-east-1
```

Then use events package from the lib
```go
import (
    "github.com/Vernacular-ai/vcore/events"
)

// If you don't want to export AWS_CREDENTIALS or already have them under different name, call
// `SetAWSCredentials` while initializing the module
err = events.SetAWSCredentials(awsAccessKey, awsSecretKey, awsRegion)

// send the actual cost event
events.SendCostEvent(
    events.NewCostEvent(events.ASR, events.GOOGLE, "client-uuid", "flow-uuid", "call-uuid", "conv-uuid")
)

// if you want to count a single event with multiple hits
events.SendCostEvent(
    events.NewCostEventWithNumHits(events.ASR, events.GOOGLE, "client-uuid", "flow-uuid", "call-uuid", "conv-uuid", 2)
)
```


## vcore/transport

### vcore/transport/amqp
### vcore/transport/redis

## vcore/utils

The vcore/utils package contains basic utility functions and file utilities for downloading, reading and writing to files.

## vcore/vorm