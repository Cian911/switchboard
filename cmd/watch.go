package main

import (
	"fmt"
	"log"

	"github.com/cian911/switchboard/utils"
	"github.com/cian911/switchboard/watcher"
	"github.com/fsnotify/fsnotify"
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
			if viper.ConfigFileUsed() != "" && ws.Watchers != nil {
				registerMultiConsumers()
			} else {
				validateFlags()
				registerSingleConsumer()
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

	viper.BindPFlag("path", runCmd.PersistentFlags().Lookup("path"))
	viper.BindPFlag("destination", runCmd.PersistentFlags().Lookup("destination"))
	viper.BindPFlag("ext", runCmd.PersistentFlags().Lookup("ext"))

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

func validateFlags() {
	if !utils.ValidatePath(viper.GetString("path")) {
		log.Fatalf("Path cannot be found. Does the path exist?: %s", viper.GetString("path"))
	}

	if !utils.ValidatePath(viper.GetString("destination")) {
		log.Fatalf("Destination cannot be found. Does the path exist?: %s", viper.GetString("destination"))
	}

	if !utils.ValidateFileExt(viper.GetString("ext")) {
		log.Fatalf("Ext is not valid. A file extention should contain a '.': %s", viper.GetString("ext"))
	}
}

func registerMultiConsumers() {
	watch, _ := fsnotify.NewWatcher()
	var pw watcher.Producer = &watcher.PathWatcher{
		Watcher: *watch,
	}

	for i, v := range ws.Watchers {
		if i == 0 {
			// Register the path and create the watcher
			pw.(*watcher.PathWatcher).Path = v.Path
		} else {
			// Add paths to this watcher, so as we don't spawn multiple
			// watcher instances.
			pw.(*watcher.PathWatcher).AddPath(v.Path)
		}

		var pc watcher.Consumer = &watcher.PathConsumer{
			Path:        v.Path,
			Destination: v.Destination,
			Ext:         v.Ext,
		}

		pw.Register(&pc)
	}

	pw.Observe()
}

func registerSingleConsumer() {
	var pw watcher.Producer = &watcher.PathWatcher{
		Path: viper.GetString("path"),
	}

	var pc watcher.Consumer = &watcher.PathConsumer{
		Path:        viper.GetString("path"),
		Destination: viper.GetString("destination"),
		Ext:         viper.GetString("ext"),
	}

	pw.Register(&pc)
	pw.Observe()
}
