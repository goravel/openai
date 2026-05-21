package facades

import (
	"fmt"

	contractsai "github.com/goravel/framework/contracts/ai"

	"github.com/goravel/openai"
)

func OpenAI(provider string) (contractsai.Provider, error) {
	if openai.App == nil {
		return nil, fmt.Errorf("please register openai service provider")
	}

	instance, err := openai.App.MakeWith(openai.Binding, map[string]any{
		"provider": provider,
	})
	if err != nil {
		return nil, err
	}

	return instance.(contractsai.Provider), nil
}
