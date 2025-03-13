# Contentstack management SDK for Go
This is the Go SDK for the contentstack management API. It is primarily
developed to be used in our terraform provider for contenstack, see
https://github.com/labd/terraform-provider-contentstack

## Example

```go
cfg := management.ClientConfig{
    BaseURL:    "https://eu-api.contentstack.com/",
    HTTPClient: httpClient,
    AuthToken:  "foobar", // Optional
}

client, err := management.NewClient(cfg)
if err != nil {
    panic(err)
}

stackAuth := management.StackAuth{
    ApiKey:          "foobar", // Required
    ManagementToken: "secret", // Optional
    Branch:          "development", // Optional
}

instance, err := client.Stack(stackAuth)
if err != nil {
    panic(err)
}

webhooks, err := stack.WebHookFetchAll(context.TODO())
if err != nil {
    panic(err)
}

```
