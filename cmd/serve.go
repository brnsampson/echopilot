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
	"os"

	// NOTE: if you fork this repo you will need to change this path.
	"github.com/brnsampson/echopilot/internal/appserver"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server in the foreground.",
	Long:  `Start a reloadable GRPC server with REST gateway.`,
	Run:   runServe,
}

func runServe(cmd *cobra.Command, args []string) {
	// NOTE: to chenge the behavior of this function change the remainder of this function.
	var srv interface {
		BlockingRun() int
	}

    srv, err := appserver.NewAppServer(cmd.Flags())
    if err != nil {
        os.Exit(1)
    }
	exitCode := srv.BlockingRun()
	os.Exit(exitCode)
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// NOTE: add any additional flags here.
	serveCmd.Flags().String("config", "", "Location of a environment config file. All flags can be set via file.")
	serveCmd.Flags().String("host", "localhost", "Address to bind GRPC server")
	serveCmd.Flags().String("ip", "127.0.0.1", "Address to bind REST gateway for grpc server")
	serveCmd.Flags().Int("port", 3000, "Address to bind REST gateway for grpc server")
	serveCmd.Flags().String("tlsCert", "", "Location of server certificate for TLS")
	serveCmd.Flags().String("tlsKey", "", "Location of server key for TLS")
	serveCmd.Flags().Bool("tlsEnabled", true, "Enable tls")
	serveCmd.Flags().Bool("tlsSkipVerify", false, "Skip TLS verification between REST proxy and GRPC server. Almost never needed.")
}
