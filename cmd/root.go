package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CMD struct {
	Configulation []Configuration `mapstructure:"cmd"`
}
type Configuration struct {
	Name       string   `mapstructure:"name"`
	Domain     string   `mapstructure:"domain"`
	Path       string   `mapstructure:"path"`
	AppFolder  string   `mapstructure:"app_folder"`
	Image      string   `mapstructure:"image"`
	CSS        string   `mapstructure:"css"`
	Script     string   `mapstructure:"script"`
	Font       string   `mapstructure:"font"`
	Constant   string   `mapstructure:"constant"`
	IgnoreDir  []string `mapstructure:"ignore_dir"`
	IgnoreFile []string `mapstructure:"ignore_file"`
}

var (
	cfgFile string
	cfgData CMD

	rootCmd = &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
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
		os.Exit(1)
	}

	err = viper.Unmarshal(&cfgData)
	if err != nil {
		fmt.Println("Error: unmarshal templatify.yaml")
		os.Exit(1)
	}
}
