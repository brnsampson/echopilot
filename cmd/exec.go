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
	"strings"

	//i NOTE: If you are forking this, then change this import to point to your repo.
	"github.com/brnsampson/echopilot/pkg/echoserver"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Executes the core functionality of your pkg code as a cli-based tool",
	Long: `Exec is the top-level command for any cli-based interactions with
the code in your pkg/ directory. This is opposed to the 'serve' function
which exposes the same functionality over REST or RPC interface.  `,
	Run: runExec,
}

func runExec(cmd *cobra.Command, args []string) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("Failed to initialize logger!")
		os.Exit(1)
	}
	sugar := logger.Sugar()
	defer sugar.Sync()

	// NOTE: To change the command being run, change the following:
	value := strings.Join(args, " ")
	result, err := echoserver.Echo(value)
	if err != nil {
		sugar.Errorf("Error in Echo: %v", err)
	}
	fmt.Println(result)
	os.Exit(0)
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// NOTE: add wny additional command line options here:
}
