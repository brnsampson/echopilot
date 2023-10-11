/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"strings"
    "time"
    "os"

	"github.com/brnsampson/echopilot/rpc/echo"
	"github.com/brnsampson/echopilot/pkg/option"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runClient,
}

func runClient(cmd *cobra.Command, args []string) {
	fmt.Println("client called")
    flags := cmd.Flags()

    addr, err := flags.GetString("addr")
	if err != nil {
		fmt.Printf("Error reading addr flag: %v", err)
        os.Exit(1)
	}

    timeout, err := flags.GetInt("timeout")
	if err != nil {
		fmt.Printf("Error reading timeout flag: %v", err)
        os.Exit(1)
	}

    t := option.NewOption(time.Duration(timeout) * time.Second)
    tlsSkipVerify, err := flags.GetBool("tlsSkipVerify")
	if err != nil {
        fmt.Printf("Error reading tlsSkipVerify flag: %v", err)
        os.Exit(1)
	}
    sv := option.NewOption(tlsSkipVerify)

	client, err := echo.NewRemoteEchoClient(addr, t, sv)
	if err != nil {
		fmt.Printf("Error while creating client: %v", err)
        os.Exit(1)
	}

    request := echo.NewStringRequest(strings.Join(args, " "))

	result, err := client.EchoString(request)
	if err != nil {
		fmt.Printf("Error during client request: %v", err)
        os.Exit(1)
	}

	fmt.Println(result)
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	clientCmd.Flags().String("addr", "127.0.0.1:8080", "Address of the echo server")
	clientCmd.Flags().Int("timeout", 10, "Request timeout (in seconds)")
	clientCmd.Flags().Bool("tlsSkipVerify", false, "Skip TLS verification when connecting to GRPC server. Useful when running server with self signed certs.")
}
