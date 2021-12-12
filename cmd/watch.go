package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	shortDesc = "Run switchboard application."
	longDesc  = "Run the switchboard application passing in the path, destination, and file type you'd like to watch for."
)

var (
	configFile string
	ws         Watchers
)

type Watchers struct {
	Watchers []Watcher `yaml:"watchers,mapstructure"`
}

type Watcher struct {
	Path        string `yaml:"path"`
	Destination string `yaml:"destination"`
	Ext         string `yaml:"ext"`
	Operation   string `yaml:"operation"`
}

func Watch() {
	var runCmd = &cobra.Command{
		Use:   "watch",
		Short: shortDesc,
		Long:  longDesc,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(ws)

			if viper.ConfigFileUsed() != "" && ws.Watchers != nil {
				/* var pw watcher.Producer = &watcher.PathWatcher{ */
				/*   Path: path, */
				/* } */
				/*  */
				/* var pc watcher.Consumer = &watcher.PathConsumer{ */
				/*   Path:        path, */
				/*   Destination: destination, */
				/*   Ext: "", */
				/* } */
				/*  */
				/* pw.Register(&pc) */
				/* pw.Observe() */
			}

		},
	}

	initCmd(*runCmd)
}

func initCmd(runCmd cobra.Command) {
	cobra.OnInitialize(initConfig)

	runCmd.PersistentFlags().StringP("path", "p", "", "Path you want to watch.")
	runCmd.PersistentFlags().StringP("destination", "d", "", "Path you want files to be relocated.")
	runCmd.PersistentFlags().StringP("ext", "e", "", "File type you want to watch for.")
	runCmd.PersistentFlags().StringVar(&configFile, "config", "", "Pass an optional config file containing multipe paths to watch.")

	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(&runCmd)
	rootCmd.Execute()
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)

		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())

			err := viper.Unmarshal(&ws)

			if err != nil {
				log.Fatalf("Unable to decode config file. Please check that it is in correct format: %v", err)
			}

			if ws.Watchers == nil || len(ws.Watchers) == 0 {
				log.Fatalf("Unable to decode config file. Please check that it is in the correct format.")
			}
		}
	}

}
