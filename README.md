# Contentstack management SDK for Go

A small Go client for the [Contentstack Content Management API](https://www.contentstack.com/docs/developers/apis/content-management-api).
It covers the parts of the API needed to manage a stack as code: content
types, global fields, entries, environments, locales and webhooks. It is
mainly developed for, and used by, the
[Terraform provider for Contentstack](https://github.com/labd/terraform-provider-contentstack),
so it does not cover the full API.

## Installation

```sh
go get github.com/labd/contentstack-go-sdk
```

All functionality lives in the `management` package:

```go
import "github.com/labd/contentstack-go-sdk/management"
```

## Requirements

- Go 1.27 or newer (see `go.mod`)

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/labd/contentstack-go-sdk/management"
)

func main() {
	ctx := context.Background()

	client, err := management.NewClient(management.ClientConfig{
		BaseURL: "https://eu-api.contentstack.com/",
	})
	if err != nil {
		log.Fatal(err)
	}

	stack, err := client.Stack(&management.StackAuth{
		ApiKey:          os.Getenv("CONTENTSTACK_API_KEY"),
		ManagementToken: os.Getenv("CONTENTSTACK_MANAGEMENT_TOKEN"),
		Branch:          "main", // optional
	})
	if err != nil {
		log.Fatal(err)
	}

	webhooks, err := stack.WebHookFetchAll(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, wh := range webhooks {
		fmt.Println(wh.UID, wh.Name)
	}

	env, err := stack.EnvironmentCreate(ctx, management.EnvironmentInput{
		Name: "staging",
		URLs: []management.EnvironmentUrl{
			{Locale: "en-us", URL: "https://staging.example.com"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("created environment", env.UID)
}
```

### Supported resources

These methods are available on a `*management.StackInstance`:

| Resource      | Methods                                                                                           |
| ------------- | ------------------------------------------------------------------------------------------------- |
| Content types | `ContentTypeCreate`, `ContentTypeUpdate`, `ContentTypeDelete`, `ContentTypeFetch`, `ContentTypeFetchAll` |
| Global fields | `GlobalFieldCreate`, `GlobalFieldUpdate`, `GlobalFieldDelete`, `GlobalFieldFetch`, `GlobalFieldFetchAll` |
| Entries       | `EntryCreate`, `EntryUpdate`, `EntryDelete`, `EntryFetch`, `EntryFetchAll`                        |
| Environments  | `EnvironmentCreate`, `EnvironmentUpdate`, `EnvironmentDelete`, `EnvironmentFetch`, `EnvironmentFetchAll` |
| Locales       | `LocaleCreate`, `LocaleUpdate`, `LocaleDelete`, `LocaleFetch`, `LocaleFetchAll`                   |
| Webhooks      | `WebHookCreate`, `WebHookUpdate`, `WebHookDelete`, `WebHookFetch`, `WebHookFetchAll`              |

On the `*management.Client` itself there are `Login` and `Stacks` (lists the
stacks the user has access to; this needs an auth token).

### Errors

API errors (HTTP 404 and 422) are returned as `*management.ErrorMessage`,
which carries the Contentstack `ErrorCode`, `ErrorMessage` and per-field
`Errors`:

```go
_, err := stack.ContentTypeFetch(ctx, "blog_post")
var apiErr *management.ErrorMessage
if errors.As(err, &apiErr) && apiErr.ErrorCode == 404 {
	// content type does not exist
}
```

## Configuration and authentication

The SDK reads no environment variables; everything is passed explicitly.

`management.ClientConfig`:

- `BaseURL` (required): the management API host for your region, for example
  `https://api.contentstack.io/` (North America) or
  `https://eu-api.contentstack.com/` (Europe).
- `HTTPClient` (optional): a custom `*http.Client`. Defaults to `&http.Client{}`.
- `AuthToken` (optional): a user auth token.

Stack requests need the stack's API key, plus either a **management token**
(`StackAuth.ManagementToken`, sent as the `authorization` header) or a user
**auth token** (`ClientConfig.AuthToken`, sent as the `authtoken` header).
When both are set, the management token is used. `client.Stack` returns an
error when neither is set.

You can also get an auth token by logging in with user credentials:

```go
err := client.Login(ctx, management.UserCredentials{
	Email:    "user@example.com",
	Password: "secret",
})
```

To log all HTTP requests and responses, use the bundled debug transport:

```go
client, err := management.NewClient(management.ClientConfig{
	BaseURL:    "https://eu-api.contentstack.com/",
	HTTPClient: &http.Client{Transport: management.DebugTransport},
})
```

## Development

The repository uses [Task](https://taskfile.dev/):

```sh
task build      # go build ./...
task test       # go test ./...
task fmt        # go fmt ./...
task coverage   # race-enabled coverage report, opened as HTML
```

CI runs `go test -race -v ./...`, `golangci-lint` and `govulncheck` on every
pull request.

## Releasing

Changes are tracked with [changie](https://changie.dev/). Add a change file
for every user-facing change:

```sh
changie new
```

On every push to `main`, the
[changie-release-action](https://github.com/labd/changie-release-action)
opens or updates a release pull request that batches the unreleased changes
and updates `CHANGELOG.md`. Merging that pull request tags the new version.

## License

[MIT](LICENSE)
