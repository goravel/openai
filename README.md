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

## Testing

Run command below to run test:

```bash
go test ./...
```
