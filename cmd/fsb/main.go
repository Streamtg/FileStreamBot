package main

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/routes" // <-- Importa los handlers
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

const versionString = "3.1.0"

var rootCmd = &cobra.Command{
	Use:               "fsb [command]",
	Short:             "Telegram File Stream Bot",
	Long:              "Telegram Bot to generate direct streamable links for telegram media.",
	Example:           "fsb run --port 8080",
	Version:           versionString,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// Registrar handler HTTP para el reproductor multimedia
	http.HandleFunc("/watch", routes.WatchHandler)

	// Comandos de la CLI
	config.SetFlagsFromConfig(runCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(sessionCmd)
	rootCmd.SetVersionTemplate(fmt.Sprintf(`Telegram File Stream Bot version %s`, versionString))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
