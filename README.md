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

See the flag `-h` for additional information and the flag `-h` in other commands, for example `./pqcom app chat -h`.

## Benchmarking

This application is capable of benchmarking all the available algorithms by using the the tests defined in `crypto/kem_test.go` and `crypto/sign_test.go`. Addionally the application [`benchstat`](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) can be used on the results. Few examples of benchmarks:

1. Benchmark only the sum of all KEM and Digital signature operations together

   ```sh
   go test -count=8 -timeout=24h -bench="Kem$|Signature$" -run=ˆ# ./... | tee out.txt && benchstat out.txt
   ```

2. Benchmark all operations and separate operations for all KEMs and Digital signatures

   ```sh
   go test -count=8 -timeout=24h -bench="Kem|Signature" -run=ˆ# ./... | tee out.txt && benchstat out.txt
   ```

## Examples for communicating

The listening port and target/listening address can be changed using flags. For more specific information check `./pqcom app -h`

### Configuration

Before chatting or file exchange can be done a confoguration file has to be specified. Path to a configuration file can be specified in 3 ways:

1. `ENV` variable called `PQCOM_CONFIG`
2. Using the `--config` flag
3. Default config location at `$HOME/.config/pqcom/pqcom.json`

To generate example config files for server and client run

```sh
./pqcom config gen
```

To specify different algorithms for communicating run

```sh
./pqcom config gen -h
```

and see available flags. To check the list available algorithsm run

```sh
./pqcom config list
```

### Chatting

To run the app in chat mode, the clint runs

```sh
./pqcom app chat -c --config pqcom_client.json
```

and the server runs:

```sh
./pqcom app chat -l --config pqcom_server.json
```

### File exchange

To exchange files, the receiver has to run

```sh
./pqcom app receive --config pqcom_server.json > output.txt
```

and the client that sends the file has to run

```sh
cat test.txt | ./pqcom app send --config pqcom_client.json
```

## How to test

Install go [here](https://go.dev/doc/install).

While the root git directory run with:

```
go test -v ./...
```

## Wireshark integration

The folder `wireshark` contains a `README` about wireshark dissector integration for the implemented protocol.
