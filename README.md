# plex

[![Go Reference](https://pkg.go.dev/badge/thde.io/plex.svg)](https://pkg.go.dev/thde.io/plex) [![test](https://github.com/thde/plex/actions/workflows/test.yml/badge.svg)](https://github.com/thde/plex/actions/workflows/test.yml)

Go client for the [Plex Media Server API](https://developer.plex.tv/pms/), generated from the official OpenAPI spec using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

## Installation

```sh
go get thde.io/plex
```

## Usage

```go
client, err := plex.New(
    "http://192.168.1.10:32400",
    plex.WithToken("your-plex-token"),
    plex.WithClientIdentifier("a-unique-stable-client-id"),
    plex.WithProduct("my-app", "1.0"),
)
if err != nil {
    log.Fatal(err)
}

resp, err := client.GetIdentityWithResponse(ctx)
```

Every Plex Media Server has its own address, so `New` takes the base URL of the server. Use the server's LAN address or its `https://<ip>.<id>.plex.direct:32400` URL.

Plex Media Server responds with XML by default. `New` sends `Accept: application/json` on every request, because the generated response types only decode JSON. If you build a client with `NewClientWithResponses` instead, set the header yourself.

This client covers the Media Server API only. Getting a token from plex.tv (PIN or JWT flow) is not part of the Media Server spec. See [Authenticating with Plex](https://developer.plex.tv/pms/#section/API-Info/Authenticating-with-Plex).

### Options

| Option                     | Header                                |
| -------------------------- | ------------------------------------- |
| `WithToken`                | `X-Plex-Token`                        |
| `WithClientIdentifier`     | `X-Plex-Client-Identifier`            |
| `WithProduct`              | `X-Plex-Product`, `X-Plex-Version`    |
| `WithUserAgent`            | `User-Agent`                          |

## OpenAPI spec

Plex does not publish the spec as a separate file. It is embedded in the Redoc documentation page. `cmd/plex-openapi` extracts it:

```sh
go run thde.io/plex/cmd/plex-openapi@latest > openapi.json
```

The spec is not checked in. `go generate` downloads it, and `overlay.yml` fixes names in the upstream spec that would collide or produce unreadable Go identifiers.

## Development

```sh
make generate  # download the spec and regenerate the client
make test      # run tests
make lint      # run linters
make lint-fix  # run linters and auto-fix issues
make update    # update dependencies and regenerate
```

## License

MIT
