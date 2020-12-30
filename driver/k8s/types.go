package k8s

import (
	"sshare/logger"
	"sshare/types"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client stores data for the client deriver
type Client struct {
	Client    *kubernetes.Clientset
	namespace string
	log       *zap.SugaredLogger
}

// New initializes a new client for the driver
func New() types.DriverAdapter {
	log := logger.GetInstance()
	client := &Client{
		log: log,
	}
	client.init()

	return client
}

func (c *Client) init() {
	// create the clientset
	clientset, err := kubernetes.NewForConfig(c.config())
	if err != nil {
		panic(err.Error())
	}

	c.namespace = viper.GetString("namespace")

	c.Client = clientset
}

func (c *Client) config() *rest.Config {
	var err error
	var config *rest.Config
	if viper.GetBool("server.in-cluster") {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	} else {
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", viper.GetString("kubeconfig"))
		if err != nil {
			panic(err.Error())
		}
	}

	return config
}
