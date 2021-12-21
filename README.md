# alvarium-sdk-go

This is a re-implementation of the Alvarium SDK in Go. It borrows conceptually from the [original implementation](https://github.com/project-alvarium/go-sdk) while updating
functionality for [IOTA Streams](https://github.com/iotaledger/streams) support.

A direct link to the IOTA C bindings directory is [here](https://github.com/iotaledger/streams/tree/v1.1.0/bindings/c).

**Environment Variable Setup**

In order to execute the unit tests for the SDK, you will need to export an environment variable that points to the directory containing the C bindings once they have been compiled.
As an example from my local machine, I have put the following directive into my `~/.profile`

`export LD_LIBRARY_PATH=/home/developer/Documents/github/project-alvarium/alvarium-sdk-go/internal/iota/include`

This variable is necessary for `make test` to execute a build prior to running the suite of unit tests

# SDK Interface

The SDK provides a minimal API -- NewSdk(), Create(), Mutate(), Transit(), Publish() and BootstrapHandler().

### NewSdk()

```go
func NewSdk(annotators []annotator.Contract, cfg config.SdkInfo, logger logInterface.Logger) interfaces.Sdk
```

Used to instantiate a new SDK instance with the specified list of annotators.

Takes a list of annotators, a populated configuration and a logger instance. Returns an SDK instance.

### Create()

```go
func (s *sdk) Create(ctx context.Context, data []byte)
```

Used to register creation of new data with the SDK. Passes data through the SDK instance's list of annotators.

SDK instance method. Parameters include:

- ctx -- Provide a context that may be used by individual annotators
- data -- The data being created represented as a byte array

### Mutate()

```go
func (s *sdk) Mutate(ctx context.Context, old, new []byte)
```

Used to register mutation of existing data with the SDK. Passes data through the SDK instance's list of annotators.

SDK instance method. Parameters include:

- ctx -- Provide a context that may be used by individual annotators

- old -- The source data item that is being modified, represented as a byte array

- new -- The new data item resulting from the change, represented as a byte array

Calling this method will link the old data to the new in a lineage. Specific annotations will be applied to the `new` data element.

### Transit()

```go
func (s *sdk) Transit(ctx context.Context, data []byte)
```

Used to annotate data that is neither originated or modified but simply handed from one application to another.

SDK instance method. Parameters include:

- ctx -- Provide a context that may be used by individual annotators

- data -- The data being handled represented as a byte array

### Publish()

```go
func (s *sdk) Publish(ctx context.Context, data []byte)
```

Used to annotate data that is neither originated or modified but **before** being handed to another application.

SDK instance method. Parameters include:

- ctx -- Provide a context that may be used by individual annotators

- data -- The data being handled represented as a byte array

### BootstrapHandler()

```go
BootstrapHandler(ctx context.Context, wg *sync.WaitGroup) bool
```

SDK instance method. Ensures clean shutdown of the SDK and associated resources.
