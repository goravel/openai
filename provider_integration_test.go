package openai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	contractsai "github.com/goravel/framework/contracts/ai"
	contractsconfig "github.com/goravel/framework/contracts/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const integrationDefaultTextModel = "gpt-4.1-nano"

func TestProviderPromptIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping OpenAI integration test in short mode")
	}

	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY is not set")
	}

	model := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))
	if model == "" {
		model = integrationDefaultTextModel
	}

	providerConfig := contractsai.ProviderConfig{Key: apiKey}
	providerConfig.Models.Text.Default = model

	provider, err := NewOpenAI(integrationConfigStub{providerConfig: providerConfig}, "openai")
	require.NoError(t, err)
	require.NotNil(t, provider)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	response, err := provider.Prompt(ctx, contractsai.AgentPrompt{
		Agent: integrationAgentStub{},
		Input: "Reply with OK only.",
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Usage())

	assert.NotEmpty(t, strings.TrimSpace(response.Text()))
	assert.Positive(t, response.Usage().Total())
	assert.Empty(t, response.ToolCalls())
}

type integrationConfigStub struct {
	providerConfig contractsai.ProviderConfig
}

var _ contractsconfig.Config = integrationConfigStub{}

func (stub integrationConfigStub) Env(_ string, defaultValue ...any) any {
	return firstOrZero(defaultValue)
}

func (stub integrationConfigStub) EnvString(_ string, defaultValue ...string) string {
	return firstOrZero(defaultValue)
}

func (stub integrationConfigStub) EnvBool(_ string, defaultValue ...bool) bool {
	return firstOrZero(defaultValue)
}

func (stub integrationConfigStub) Add(string, any) {}

func (stub integrationConfigStub) Get(_ string, defaultValue ...any) any {
	return firstOrZero(defaultValue)
}

func (stub integrationConfigStub) GetString(_ string, defaultValue ...string) string {
	return firstOrZero(defaultValue)
}

func (stub integrationConfigStub) GetInt(_ string, defaultValue ...int) int {
	return firstOrZero(defaultValue)
}

func (stub integrationConfigStub) GetBool(_ string, defaultValue ...bool) bool {
	return firstOrZero(defaultValue)
}

func (stub integrationConfigStub) GetDuration(_ string, defaultValue ...time.Duration) time.Duration {
	return firstOrZero(defaultValue)
}

func (stub integrationConfigStub) UnmarshalKey(key string, rawVal any) error {
	if key != "ai.providers.openai" {
		return fmt.Errorf("unexpected config key: %s", key)
	}

	providerConfig, ok := rawVal.(*contractsai.ProviderConfig)
	if !ok {
		return fmt.Errorf("unexpected config type: %T", rawVal)
	}

	*providerConfig = stub.providerConfig

	return nil
}

type integrationAgentStub struct{}

func (integrationAgentStub) Instructions() string                 { return "" }
func (integrationAgentStub) Messages() []contractsai.Message      { return nil }
func (integrationAgentStub) Middleware() []contractsai.Middleware { return nil }
func (integrationAgentStub) Tools() []contractsai.Tool            { return nil }

func firstOrZero[T any](values []T) T {
	var zero T
	if len(values) == 0 {
		return zero
	}

	return values[0]
}
