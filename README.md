# OpenAI

The OpenAI provider for `facades.AI()` of Goravel.

## Version

| goravel/openai | goravel/framework |
|----------------|-------------------|
| v1.18.x        | v1.18.x           |

## Install

Run the command below in your project to install the package automatically:

```bash
./artisan package:install github.com/goravel/openai
```

This registers the service provider and updates `config/ai.go` so `ai.providers.openai.via` resolves through `openaifacades.OpenAI("openai")`.

Or check [the setup file](./setup/setup.go) to install the package manually.

## Custom Failover

The provider marks these OpenAI HTTP errors as failoverable by default:

| Status | Reason |
|--------|--------|
| `429 Too Many Requests` | `rate_limited` |
| `402 Payment Required` | `insufficient_credits` |
| `503 Service Unavailable` | `provider_overloaded` |

Configure `failover` rules to add OpenAI-specific error message mappings. Plain strings use substring matching, and slash-delimited strings use Go regular expressions.

```go
"openai": map[string]any{
	"key": config.Env("OPENAI_API_KEY", ""),
	"failover": map[string][]string{
		"context_length_exceeded": {
			"maximum context length",
			"/(?i)context.*length/",
		},
	},
	"via": func() (ai.Provider, error) {
		return openaifacades.OpenAI("openai")
	},
}
```

## Testing

Run command below to run all tests:

```bash
go test ./...
```

Run the live OpenAI smoke test with a real API key:

```bash
OPENAI_API_KEY=your-key go test -run '^TestProviderPromptIntegration$' -v ./...
```

The smoke test skips automatically when `OPENAI_API_KEY` is not set.
