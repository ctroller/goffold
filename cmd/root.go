package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ctroller/goffold/internal/dependencies"
	"github.com/ctroller/goffold/internal/inject"
	"github.com/ctroller/goffold/internal/template"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goffold <template>",
	Short: "Create golang projects with your preferred structure",
	Long: `Goffold is a CLI library to create golang projects with your preferred structure.
This application is a tool to generate the needed files to quickly create a Golang application, using predefined templates.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}

		if !template.IsExisting(args[0]) {
			return fmt.Errorf("unknown template %s in template path %s", args[0], template.TemplatesPath)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) { 
		fmt.Println("Hello, World! args are: " + strings.Join(args, " "))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&template.TemplatesPath, "templates", "t", "./templates", "path to templates")
	template.Fs = afero.NewOsFs()

	initResolvers()
	loadTemplates()
}

func initResolvers() {
	dependencies.RegisterResolver(dependencies.DependencyResolver{ 
		Type: "go",
		Handler: dependencies.GoDependencyHandler(inject.Defaults),
	})
}

func loadTemplates() {
	// Load templates from the default location
	template.InitTemplates()
}