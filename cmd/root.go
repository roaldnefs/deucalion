// Copyright © 2019 Roald Nefs <info@roaldnefs.com>

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// The default base URL for API requests to Prometheus.
	defaultBaseURL = "http://localhost:9090/"
)

var (
	// Used for command flags
	cfgFile string
	debug   bool

	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   "deucalion",
		Short: "A tool for querying a Prometheus instance for alerts and run commands based on the result.",
		Long:  `A tool for querying a Prometheus instance for alerts and run commands based on the result.`,
		Run: func(cmd *cobra.Command, args []string) {
			httpAPI := newAPI()
			handleAlerts(httpAPI)
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug logging")

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.deucalion.yaml)")

	// Base URL for API requests to Prometheus. Defaults to localhost, but can
	// be set to a external hosted Prometheus server. The base URL should
	// always be specified with a trailing slash.
	rootCmd.PersistentFlags().StringP("url", "u", defaultBaseURL, "Promtheus URL")
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))

	// Commands that are executed based on the Prometheus API resonse.
	rootCmd.PersistentFlags().StringP("silent", "s", "", "command to execute when alerts aren't firing")
	viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))
	rootCmd.PersistentFlags().StringP("firing", "f", "", "command to execute when alerts are firing")
	viper.BindPFlag("firing", rootCmd.PersistentFlags().Lookup("firing"))
	rootCmd.PersistentFlags().StringP("warning", "w", "", "command to execute when alerts are firing and below specified severity")
	viper.BindPFlag("warning", rootCmd.PersistentFlags().Lookup("warning"))

	// Only execute the firing command if the severity is matched.
	rootCmd.PersistentFlags().StringP("severity", "", "", "the severity label")
	viper.BindPFlag("severity", rootCmd.PersistentFlags().Lookup("severity"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".deucalion" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".deucalion")
	}

	// Read in environment variable that match.
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if debug {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}

// newAPI returns a Promtheus v1 API client.
func newAPI() v1.API {
	baseURL := viper.GetString("url")

	client, err := api.NewClient(api.Config{Address: baseURL})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	httpAPI := v1.NewAPI(client)

	return httpAPI
}

// handleAlerts
func handleAlerts(httpAPI v1.API) error {
	commandSilent := viper.GetString("silent")
	commandFiring := viper.GetString("firing")
	commandWarning := viper.GetString("warning")
	severity := viper.GetString("severity")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	alertsResult, err := httpAPI.Alerts(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	command := commandSilent

	// Change the command when alerts are firing
	for _, alert := range alertsResult.Alerts {
		// Get the severity (label) of the alert
		alertSeverity := string(alert.Labels["severity"])
		// Check if the alert is firing
		if alert.State == v1.AlertStateFiring {
			// Skip if a severity is given but does not match the label
			if severity != "" && alertSeverity != severity {
				command = commandWarning
				continue
			}

			command = commandFiring
			break
		}
	}

	err = execute(command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

// execute returns an error if it failes to execute the command
func execute(command string) error {
	commandSlice := strings.Fields(command)
	name := commandSlice[0]
	args := commandSlice[1:]

	cmd := exec.Command(name, args...)
	err := cmd.Run()

	return err
}
