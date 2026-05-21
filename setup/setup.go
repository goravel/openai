package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/path"
)

func main() {
	setup := packages.Setup(os.Args)
	aiConfigPath := path.Config("ai.go")
	appConfigPath := path.Config("app.go")
	moduleImport := setup.Paths().Module().Import()
	serviceProvider := "&openai.ServiceProvider{}"
	aiProviderContract := "github.com/goravel/framework/contracts/ai"
	openAIFacadesImport := moduleImport + "/facades"
	via := `func() (ai.Provider, error) {
			return openaifacades.OpenAI("openai")
		}`
	openAIProviderConfig := match.Config("ai.providers.openai")

	setup.Install(
		modify.When(func(_ map[string]any) bool {
			return !env.IsBootstrapSetup()
		}, modify.GoFile(appConfigPath).
			Find(match.Imports()).Modify(modify.AddImport(moduleImport)).
			Find(match.Providers()).Modify(modify.Register(serviceProvider))),

		modify.When(func(_ map[string]any) bool {
			return env.IsBootstrapSetup()
		}, modify.RegisterProvider(moduleImport, serviceProvider)),

		modify.GoFile(aiConfigPath).Find(match.Imports()).Modify(
			modify.AddImport(aiProviderContract),
			modify.AddImport(openAIFacadesImport, "openaifacades"),
		).Find(openAIProviderConfig).Modify(modify.AddConfig("via", via)),
	).Uninstall(
		modify.WhenFileExists(aiConfigPath, modify.GoFile(aiConfigPath).
			Find(openAIProviderConfig).Modify(modify.ReplaceConfig("via", `""`)).
			Find(match.Imports()).Modify(
			modify.RemoveImport(aiProviderContract),
			modify.RemoveImport(openAIFacadesImport, "openaifacades"),
		)),

		modify.When(func(_ map[string]any) bool {
			return !env.IsBootstrapSetup()
		}, modify.GoFile(appConfigPath).
			Find(match.Providers()).Modify(modify.Unregister(serviceProvider)).
			Find(match.Imports()).Modify(modify.RemoveImport(moduleImport))),

		modify.When(func(_ map[string]any) bool {
			return env.IsBootstrapSetup()
		}, modify.UnregisterProvider(moduleImport, serviceProvider)),
	).Execute()
}
