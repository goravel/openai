package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/path"
)

func main() {
	setup := packages.Setup(os.Args)
	aiConfigPath := path.Config("ai.go")
	moduleImport := setup.Paths().Module().Import()
	serviceProvider := "&openai.ServiceProvider{}"
	aiProviderContract := "github.com/goravel/framework/contracts/ai"
	openAIFacadesImport := moduleImport + "/facades"
	env := `
OPENAI_API_KEY=
OPENAI_BASE_URL=
`
	provider := `map[string]any{
		"key": config.Env("OPENAI_API_KEY", ""),
		"models": map[string]any{
			"text": map[string]any{
				"default": "",
			},
			"audio": map[string]any{
				"default": "",
			},
			"transcription": map[string]any{
				"default": "",
			},
			"image": map[string]any{
				"default": "",
			},
		},
		"url": config.Env("OPENAI_BASE_URL", ""),
		"via": func() (ai.Provider, error) {
			return openaifacades.OpenAI("openai")
		},
	}`
	aiProvidersConfig := match.Config("ai.providers")

	setup.Install(
		modify.RegisterProvider(moduleImport, serviceProvider),

		modify.GoFile(aiConfigPath).Find(match.Imports()).Modify(
			modify.AddImport(aiProviderContract),
			modify.AddImport(openAIFacadesImport, "openaifacades"),
		).Find(aiProvidersConfig).Modify(modify.AddConfig("openai", provider)),

		modify.WhenFileExists(path.Base(".env"), modify.WhenFileNotContains(path.Base(".env"), "OPENAI_API_KEY", modify.File(path.Base(".env")).Append(env))),
		modify.WhenFileExists(path.Base(".env.example"), modify.WhenFileNotContains(path.Base(".env.example"), "OPENAI_API_KEY", modify.File(path.Base(".env.example")).Append(env))),
	).Uninstall(
		modify.WhenFileExists(aiConfigPath, modify.GoFile(aiConfigPath).
			Find(aiProvidersConfig).Modify(modify.RemoveConfig("openai")).
			Find(match.Imports()).Modify(
			modify.RemoveImport(aiProviderContract),
			modify.RemoveImport(openAIFacadesImport, "openaifacades"),
		)),

		modify.UnregisterProvider(moduleImport, serviceProvider),
	).Execute()
}
