# sshare

sshare is an easy way to expose your local service to the world by using your own infrastructure.

```
$ sshare client 8080 --https-redirect
sshare 1.0.0 🚀

Address: https://ca3150dd-6fca-41e6-8896-8b20a87bd925.my.domain -> 0.0.0.0:8080
HTTPs redirect: enabled
Address: http://ca3150dd-6fca-41e6-8896-8b20a87bd925.my.domain -> https://ca3150dd-6fca-41e6-8896-8b20a87bd925.my.domain
```

## Requirements

- Kubernetes cluster >= v1.8 with nginx-ingress

## How does is work?

Sshare is using Kubernetes as a backend to create resources that allow exposing your local service.

1. Sshare client sends a request to a `sshare server`
2. Sshare server creates needed resources (OpenSSH server, service, ingress) in the Kubernetes cluster
3. When the backend is ready in the Kubernetes cluster, `sshare client` creates an SSH tunnel between the remote (sshd server) host and the client
4. Sshare forwards connections from the remote server to the local service that was exposed

## Usage

```
Sshare is an easy way to share your local server with the world

Usage:
  sshare [command]

Available Commands:
  client      This command creates a secure tunnel that exposes your local port
  completion  Generate completion script
  help        Help about any command
  server      Runs server that creates a backend for client request

Flags:
      --config string      config file (default is $HOME/.sshare.yaml)
  -h, --help               help for sshare
      --log-level string   logging level (debug,info,warn,error) (default "info")

Use "sshare [command] --help" for more information about a command.
```

### Server

```
Runs server that creates a backend for client request.

Usage:
  sshare server [flags]

Flags:
      --address ip                       address to listen on (default 0.0.0.0)
      --auth-token string                define authorization token that is required from a client
      --backend-domain string            domain name that is used for public access (default "sshare.io")
      --backend-https-enabled            set true if backend supports HTTPs connection
      --backend-ready-timeout duration   time after which the backend is reported as not ready (default 2m0s)
      --client-session-timeout int32     time in seconds after which a session for client is closed (0 means no limit)
      --driver string                    driver that is used to create backend (default "kubernetes")
  -h, --help                             help for server
      --in-cluster                       run server in Kubernetes cluster
      --kubeconfig string                path to the kubeconfig file (default "~/.kube/config")
  -n, --namespace string                 namespace scope where SSHD instances are created (default "default")
      --port int32                       port to listen on (default 50041)
      --metrics-port int32               port that metrics are exposed on (default 2112)
      --tls-ca string                    The TLS CA file (default "~/.sshare/ca.pem")
      --tls-cert string                  The TLS cert file (default "~/.sshare/cert.pem")
      --tls-enabled                      enable TLS for connection between client and server
      --tls-key string                   The TLS key file (default "~/.sshare/key.pem")
      --tls-port int32                   port to listen on for TLS connection (default 50040)

Global Flags:
      --config string      config file (default is $HOME/.sshare.yaml)
      --log-level string   logging level (debug,info,warn,error) (default "info")
```

### Client

```
The example below exposes local port 9090:

  $ sshare client 9090

Expose only TCP for port 9090:

  $ sshare client 9090 --tcp

Usage:
  sshare client [PORT] [flags]

Flags:
  -h, --help                    help for client
      --http-enable-cors        enable CORS
      --https-redirect          redirect HTTP to HTTPS
      --server-address string   server address (default "localhost:50041")
      --tcp                     expose TCP port (for a service that does not support HTTP protocol)
      --tls-disabled            disable TLS for connection to the server
      --token string            authorization token

Global Flags:
      --config string      config file (default is $HOME/.sshare.yaml)
      --log-level string   logging level (debug,info,warn,error) (default "info")
```

## Example of usage

### Run sshare server on a local machine

In the following example we run `sshare server` and expose a nginx container by using `sshare client`.

The example assumes that you have access to a Kubernetes cluster with configured nginx-ingress. The `sshare server` uses a default Kubernetes configuration available on your local machine.

1. Run `sshare server`

```
$ sshare server
{"level":"info","ts":1609336286.353205,"caller":"grpc/server.go:243","msg":"sshare gRPC server","version":"1.0.0","address":"0.0.0.0:50041"}
{"level":"info","ts":1609336286.353504,"caller":"grpc/server.go:212","msg":"TLS is disabled. Skipping TLS Responder"}
{"level":"info","ts":1609336286.353702,"caller":"grpc/server.go:295","msg":"Running Prometheus metrics","address":":2112","endpoint":"/metrics"}
```

2. Run a nginx container locally

```
$ docker run --rm -p 8080:80 nginx
```

3. Run `sshare client`

```
$ sshare client 8080 --tls-disabled
sshare 1.0.0 🚀

Address: http://74ba1e74-01d4-4b60-be9a-46992e392ee9.sshare.io -> 0.0.0.0:8080

```

### Connect to a sshare server runs on a Kubernetes cluster

1. [Here](examples/kubernetes) you can find an example of how to setup `sshare server` on a Kubernetes cluster
2. If `sshare server` is ready you can connect to the server from a local machine

```
sshare client 8080 --server-address server.sshare.mydomain.com:50041
sshare 1.0.0 🚀

Address: https://d22336c3-740a-478d-bdf9-d438d197c957.sshare.mydomain.com -> 0.0.0.0:8080
Address: http://d22336c3-740a-478d-bdf9-d438d197c957.sshare.mydomain.com -> 0.0.0.0:8080
```

That's all :)