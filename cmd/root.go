package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nathaponb-epic/templatify/cmd/version"
	"github.com/nathaponb-epic/templatify/pkg/utils"
)

var (
	cfgFile string
	cfgData utils.Commands

	rootCmd = &cobra.Command{
		Use:   "templatify",
		Short: "ADMD templates manupulating CLI",
		RunE: func(cmd *cobra.Command, args []string) error {

			// prompt to user choose a command
			cmds := []string{"cdnify", "localify"}
			prompt := promptui.Select{
				Label: "Select command",
				Items: cmds,
			}
			_, result, err := prompt.Run()
			if err != nil {
				return err
			}

			var cdnifyObj utils.Configuration
			var localifyObj utils.Configuration

			yamlConfigCmds := make(map[string]bool)
			for _, v := range cfgData.Configulation {

				yamlConfigCmds[v.Name] = true

				if v.Name == "cdnify" {
					cdnifyObj = v
				} else if v.Name == "localify" {
					localifyObj = v
				}
			}

			// run command based on user select option prompt
			switch result {
			case "cdnify":

				// check for cmd config
				if !yamlConfigCmds["cdnify"] {
					return errors.New("cdnify attribute not found on config file")
				}

				return utils.Walker(cdnifyObj)

			case "localify":

				// check for cmd config
				if !yamlConfigCmds["localify"] {
					return errors.New("localify attribute not found on config file")
				}

				return utils.Walker(localifyObj)
			default:
				return nil
			}
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(version.VersionCommand)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

		home, err := os.Getwd()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("templatify")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error: templatify.yaml not found")
		// os.Exit(1)
	}

	err = viper.Unmarshal(&cfgData)
	if err != nil {
		fmt.Println("Error: unmarshal templatify.yaml")
		// os.Exit(1)
	}
}
