package cmd

import (
	"os"

	"github.com/ctroller/goffold/internal/dependencies"
	"github.com/ctroller/goffold/internal/inject"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goffold",
	Short: "Create golang projects with your preferred structure",
	Long: `Goffold is a CLI library to create golang projects with your preferred structure.
This application is a tool to generate the needed files to quickly create a Golang application, using predefined templates.`,
	Run: func(cmd *cobra.Command, args []string) { },
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
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goffold.yaml)")
	initResolvers()
}

func initResolvers() {
	dependencies.RegisterResolver(dependencies.DependencyResolver{ 
		Type: "go",
		Handler: dependencies.GoDependencyHandler(inject.Defaults),
	})
}
