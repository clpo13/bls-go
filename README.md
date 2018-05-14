# bls-go

**bls-go** is a Go interface for the public API provided by the United States
[Bureau of Labor Statistics](https://www.bls.gov/).

## Requirements

- [Go](https://golang.org)
- BLS.gov API key (not strictly required, but highly recommended; you can
    request one [here](https://data.bls.gov/registrationEngine/))

## Installation

After cloning the repository, build and install the library with
`go install github.com/clpo13/bls-go`. See the [Usage](#Usage) section for how
to interact with the library.

Alternatively, you can call `go get -u github.com/clpo13/bls-go` to fetch and
install the latest version of the library directly to your GOPATH.

## Usage

In your `import` statement, add `"github.com/clpo13/bls-go"`. Now, you have
access to `blsgo.GetData`, `blsgo.Payload`, and `blsgo.ResultData`, as well as
a few other helper structs and functions. Run `go doc github.com/clpo13/bls-go`
to get some basic information on the available objects.

More detailed usage notes coming soon.

An example program using this library can be found at
<https://github.com/clpo13/bls-go-example.git>.

## Contributing

Issues and pull requests are always welcome. Please file any bug reports or
feature requests using the GitHub [issues page](https://github.com/clpo13/bls-go/issues).

### To do list

- [ ] Better error handling
- [ ] Online API docs

## License

This program is available under the terms of the Apache 2.0 license, the text
of which can be found in [LICENSE](LICENSE) or at
<https://www.apache.org/licenses/LICENSE-2.0>.
