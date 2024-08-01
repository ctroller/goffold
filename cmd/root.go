package cmd

import (
	"embed"
	"os"

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
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := initTemplates()
		if err != nil {
			return err
		}

		tpl, err := template.LoadTemplate(args[0])
		if err != nil {
			return err
		}
		err = tpl.Execute()
		if err != nil {
			return err
		}

		return nil
	},
}

var DefaultTemplates embed.FS

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var templateDir string
var output string

func init() {
	rootCmd.Flags().StringVarP(&templateDir, "templates", "t", "", "path to templates")
	rootCmd.Flags().StringVarP(&output, "output", "o", ".", "output directory. Defaults to current directory")
	initResolvers()
}

func initResolvers() {
	dependencies.RegisterResolver(dependencies.DependencyResolver{
		Type:    "go",
		Handler: dependencies.GoDependencyHandler(inject.Defaults),
	})
}

func initTemplates() error {
	template.OsFs = afero.NewOsFs()
	if templateDir == "" {
		templateDir = "templates"
		template.TemplateFs = afero.FromIOFS{FS: DefaultTemplates}
	} else {
		template.TemplateFs = template.OsFs
	}

	template.TemplateDir = templateDir
	template.OutputDir = output

	return template.ValidateConfig()
}
