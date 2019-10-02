/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/coreos/go-systemd/daemon"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var state string
var unset bool

// notifyCmd represents the notify command
var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "send an sd-notify message.",
	Long: `This command will send updates if the service is configured to be
run as a systemd unit with Type=notify.`,
	Run: runCommand,
}

func runCommand(cmd *cobra.Command, args []string) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("Failed to initialize logger!")
		os.Exit(1)
	}
	sugar := logger.Sugar()
	defer sugar.Sync()

	// If NOTIFY_SOCKET is unset then skip.
	if os.Getenv("NOTIFY_SOCKET") != "" {
		switch state {
		case "ready":
			_, err := daemon.SdNotify(unset, daemon.SdNotifyReady)
			if err != nil {
				sugar.Infof("Error when attempting to set systemd ready state %v", err)
			}
		case "stopping":
			_, err := daemon.SdNotify(unset, daemon.SdNotifyStopping)
			if err != nil {
				sugar.Infof("Error when attempting to set systemd stopping state %v", err)
			}
		case "reloading":
			_, err := daemon.SdNotify(unset, daemon.SdNotifyReloading)
			if err != nil {
				sugar.Infof("Error when attempting to set systemd reloading state %v", err)
			}
		case "watchdog":
			_, err := daemon.SdNotify(unset, daemon.SdNotifyWatchdog)
			if err != nil {
				sugar.Infof("Error when attempting to update systemd watchdog timestamp %v", err)
			}
		}
	} else {
		sugar.Infof("NOTIFY_SOCKET not defined. Skipping systemd notify.")
	}
}

func init() {
	systemdCmd.AddCommand(notifyCmd)
	notifyCmd.Flags().StringVarP(&state, "state", "s", "ready", "The state to notify systemd with")
	notifyCmd.Flags().BoolVarP(&unset, "unset", "u", false, "whether to unset the environment i.e. unset NOTIFY_SOCKET")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// notifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// notifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
