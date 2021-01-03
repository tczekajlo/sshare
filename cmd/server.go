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
	Short: "Runs server that creates a backend for client request",
	Long:  `Runs server that creates a backend for client request.`,
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
	serverCmd.Flags().Int32("client-session-timeout", 0, "time in seconds after which a session for client is closed (0 means no limit)")
	serverCmd.Flags().Bool("tls-enabled", false, "enable TLS for connection between client and server")
	serverCmd.Flags().String("tls-cert", filepath.Join(homedir.HomeDir(), ".sshare", "cert.pem"), "The TLS cert file")
	serverCmd.Flags().String("tls-key", filepath.Join(homedir.HomeDir(), ".sshare", "key.pem"), "The TLS key file")
	serverCmd.Flags().String("tls-ca", filepath.Join(homedir.HomeDir(), ".sshare", "ca.pem"), "The TLS CA file")
	serverCmd.Flags().Int32("tls-port", 50040, "port to listen on for TLS connection")
	serverCmd.Flags().String("driver", "kubernetes", "driver that is used to create backend")
	serverCmd.Flags().IP("address", net.ParseIP("0.0.0.0"), "address to listen on")
	serverCmd.Flags().Int32("port", 50041, "port to listen on")
	serverCmd.Flags().Int32("metrics-port", 2112, "port that metrics are exposed on")
	serverCmd.Flags().Bool("in-cluster", false, "run server in Kubernetes cluster")
	serverCmd.Flags().String("auth-token", "", "define authorization token that is required from a client")

	viper.BindPFlag("kubeconfig", serverCmd.Flags().Lookup("kubeconfig"))
	viper.BindPFlag("backend-domain", serverCmd.Flags().Lookup("backend-domain"))
	viper.BindPFlag("namespace", serverCmd.Flags().Lookup("namespace"))
	viper.BindPFlag("driver", serverCmd.Flags().Lookup("driver"))
	viper.BindPFlag("backend-https-enabled", serverCmd.Flags().Lookup("backend-https-enabled"))
	viper.BindPFlag("server.backend-ready-timeout", serverCmd.Flags().Lookup("backend-ready-timeout"))
	viper.BindPFlag("server.client-session-timeout", serverCmd.Flags().Lookup("client-session-timeout"))
	viper.BindPFlag("server.tls-enabled", serverCmd.Flags().Lookup("tls-enabled"))
	viper.BindPFlag("server.tls-cert", serverCmd.Flags().Lookup("tls-cert"))
	viper.BindPFlag("server.tls-key", serverCmd.Flags().Lookup("tls-key"))
	viper.BindPFlag("server.tls-ca", serverCmd.Flags().Lookup("tls-ca"))
	viper.BindPFlag("server.address", serverCmd.Flags().Lookup("address"))
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.tls-port", serverCmd.Flags().Lookup("tls-port"))
	viper.BindPFlag("server.in-cluster", serverCmd.Flags().Lookup("in-cluster"))
	viper.BindPFlag("server.metrics-port", serverCmd.Flags().Lookup("metrics-port"))
	viper.BindPFlag("server.auth-token", serverCmd.Flags().Lookup("auth-token"))
}
