/*
Copyright © 2020 tczekajlo

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
	"sshare/pkg/grpc"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client [PORT]",
	Short: "This command creates a secure tunnel that exposes your local port",
	Long: `The example below exposes local port 9090:

  $ sshare client 9090

Expose only TCP for port 9090:

  $ sshare client 9090 --tcp

  `,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		localPort, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			fmt.Printf("Cannot parse PORT: %v", err)
			os.Exit(1)
		}

		grpc.RunClient(int32(localPort))
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you can define flags and configuration settings.
	clientCmd.Flags().String("server-address", "localhost:50041", "server address")
	clientCmd.Flags().String("token", "", "authorization token")
	clientCmd.Flags().Bool("http-enable-cors", false, "enable CORS")
	clientCmd.Flags().Bool("https-redirect", false, "redirect HTTP to HTTPS")
	clientCmd.Flags().Bool("tls-disabled", false, "disable TLS for connection to the server")
	clientCmd.Flags().Bool("tcp", false, "expose TCP port (for a service that does not support HTTP protocol)")

	viper.BindPFlag("http-enable-cors", clientCmd.Flags().Lookup("http-enable-cors"))
	viper.BindPFlag("https-redirect", clientCmd.Flags().Lookup("https-redirect"))
	viper.BindPFlag("client.tls-disabled", clientCmd.Flags().Lookup("tls-disabled"))
	viper.BindPFlag("client.server-address", clientCmd.Flags().Lookup("server-address"))
	viper.BindPFlag("client.tcp", clientCmd.Flags().Lookup("tcp"))
	viper.BindPFlag("client.token", clientCmd.Flags().Lookup("token"))
}
