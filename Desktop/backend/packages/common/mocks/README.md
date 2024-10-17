Creating a mock?

These mocks are for 3rd party interfaces and implementations.

Only place mocks in this folder for which this code base DOES NOT own the actual interface and implementation.

All mocks in this folder should have a corresponding command under the `auto-gen-third-party-mocks` task in the [Taskfile.yml](../../../Taskfile.yml)

Use [mockery](https://github.com/vektra/mockery) to auto generate mocks into a mocks folder within the repective package.
