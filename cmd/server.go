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
	"net"
	"path/filepath"
	"sshare/grpc"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/homedir"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		grpc.RunServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "path to the kubeconfig file")
	serverCmd.Flags().String("backend-domain", "sshare.io", "domain name that is used for public access")
	serverCmd.Flags().StringP("namespace", "n", apiv1.NamespaceDefault, "namespace scope where SSHD instances are created")
	serverCmd.Flags().Bool("backend-https-enabled", false, "set true if backend supports HTTPs connection")
	serverCmd.Flags().Duration("backend-ready-timeout", (120 * time.Second), "time after which the backend is reported as not ready")
	serverCmd.Flags().Bool("tls-enabled", false, "enable TLS")
	serverCmd.Flags().String("tls-cert", filepath.Join(homedir.HomeDir(), ".sshare", "cert.pem"), "The TLS cert file")
	serverCmd.Flags().String("tls-key", filepath.Join(homedir.HomeDir(), ".sshare", "key.pem"), "The TLS key file")
	serverCmd.Flags().String("tls-ca", filepath.Join(homedir.HomeDir(), ".sshare", "ca.pem"), "The TLS CA file")
	serverCmd.Flags().String("driver", "kubernetes", "driver that is used to create backend")
	serverCmd.Flags().IP("address", net.ParseIP("0.0.0.0"), "Address to listen on")
	serverCmd.Flags().Int32("port", 50041, "Port to listen on")
	serverCmd.Flags().Bool("in-cluster", false, "run server in Kubernetes cluster")

	viper.BindPFlag("kubeconfig", serverCmd.Flags().Lookup("kubeconfig"))
	viper.BindPFlag("backend-domain", serverCmd.Flags().Lookup("backend-domain"))
	viper.BindPFlag("namespace", serverCmd.Flags().Lookup("namespace"))
	viper.BindPFlag("driver", serverCmd.Flags().Lookup("driver"))
	viper.BindPFlag("backend-https-enabled", serverCmd.Flags().Lookup("backend-https-enabled"))
	viper.BindPFlag("server.backend-ready-timeout", serverCmd.Flags().Lookup("backend-ready-timeout"))
	viper.BindPFlag("server.tls-enabled", serverCmd.Flags().Lookup("tls-enabled"))
	viper.BindPFlag("server.tls-cert", serverCmd.Flags().Lookup("tls-cert"))
	viper.BindPFlag("server.tls-key", serverCmd.Flags().Lookup("tls-key"))
	viper.BindPFlag("server.tls-ca", serverCmd.Flags().Lookup("tls-ca"))
	viper.BindPFlag("server.address", serverCmd.Flags().Lookup("address"))
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.in-cluster", serverCmd.Flags().Lookup("in-cluster"))
}
