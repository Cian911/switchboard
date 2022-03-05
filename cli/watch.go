package cli

import (
	"fmt"
	"log"
	"regexp"

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
	configFile   string
	ws           Watchers
	regexPattern *regexp.Regexp
	regexErr     error
)

// Watchers is a struct that contains a list of watchers.
// in yaml format
type Watchers struct {
	// Watchers is a list of watchers
	Watchers        []Watcher `yaml:"watchers,mapstructure"`
	PollingInterval int       `yaml:pollingInterval`
}

// Watcher is a struct that contains a path, destination, and file extention and event operation.
// in yaml format
type Watcher struct {
	// Path is the path you want to watch
	Path string `yaml:"path"`
	// Destination is the path you want files to be relocated
	Destination string `yaml:"destination"`
	// Ext is the file extention you want to watch for
	Ext string `yaml:"ext"`
	// Operation is the event operation you want to watch for
	// CREATE, MODIFY, REMOVE, CHMOD, WRITE itc.
	Operation string `yaml:"operation"`
}

// Watch is the main function that runs the watcher.
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
	runCmd.PersistentFlags().IntP("poll", "", 60, "Specify a polling time in seconds.")
	runCmd.PersistentFlags().StringVar(&configFile, "config", "", "Pass an optional config file containing multiple paths to watch.")
	runCmd.PersistentFlags().StringP("regex-pattern", "r", "", "Pass a regex pattern to watch for any files mathcing this pattern.")

	viper.BindPFlag("path", runCmd.PersistentFlags().Lookup("path"))
	viper.BindPFlag("destination", runCmd.PersistentFlags().Lookup("destination"))
	viper.BindPFlag("ext", runCmd.PersistentFlags().Lookup("ext"))
	viper.BindPFlag("poll", runCmd.PersistentFlags().Lookup("poll"))
	viper.BindPFlag("regex-pattern", runCmd.PersistentFlags().Lookup("regex-pattern"))

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

	if len(viper.GetString("regex-pattern")) > 0 {
		// Validate regex pattern
		regexPattern, regexErr = utils.ValidateRegexPattern(viper.GetString("regex-pattern"))

		if regexErr != nil {
			log.Fatalf("Regex pattern is not valid. Please check it again: %v", regexErr)
		}
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

	pi := viper.GetInt("poll")

	if &ws.PollingInterval != nil {
		pi = ws.PollingInterval
	}

	log.Println("Observing")
	pw.Observe(pi)
}

func registerSingleConsumer() {
	var pw watcher.Producer = &watcher.PathWatcher{
		Path: viper.GetString("path"),
	}

	var pc watcher.Consumer = &watcher.PathConsumer{
		Path:        viper.GetString("path"),
		Destination: viper.GetString("destination"),
		Ext:         viper.GetString("ext"),
		Pattern:     *regexPattern,
	}

	pw.Register(&pc)

	log.Println("Observing")
	pw.Observe(viper.GetInt("poll"))
}
