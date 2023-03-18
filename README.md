# PqCom

This repository includes the implementations for dilithium and kyber post quantum algorithms. It also includes an implementation of post-quantum safe comunication protocol and an application for exchaning files/chatting utilizing this protocol.

## How to compile

Install go [here](https://go.dev/doc/install).

While in the root git directory run with:

```
go build -o pqcom
```

## How to run

Either run the compiled binary from [compilation](#how-to-build) or run with:

```sh
go run .
```

See the flag `-h` for additional information.

## How to test

Install go [here](https://go.dev/doc/install).

While the root git directory run with:

```
go test -v ./...
```

## Wireshark integration

The folder `wireshark` contains a `README` about wireshark dissector integration for the implemented protocol.

<!-- TODO: update reame -->
<!-- TODO: Add pre-commit -->
<!-- TODO: check all todos -->
