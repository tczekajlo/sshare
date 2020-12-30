package k8s

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sshare/logger"

	pb "sshare/protobuf"

	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) isDeploymentExists(name string) bool {
	log := logger.GetInstance()
	deploymentsClient := c.Client.AppsV1().Deployments(c.namespace)

	_, err := deploymentsClient.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {

		log.Debugw("Failed to get latest version of Deployment",
			"name", name,
			"error", err)
		return false
	}

	return true
}

func (c *Client) getDeployment(data *pb.BackendData) (*appsv1.Deployment, error) {
	client := c.Client.AppsV1().Deployments(c.namespace)

	resp, err := client.Get(context.TODO(), data.Name, metav1.GetOptions{})
	if err != nil {

		c.log.Debugw("Failed to get latest version of Deployment",
			"name", data.Name,
			"stream-id", data.StreamID,
			"error", err)
		return nil, err
	}

	return resp, nil
}

func (c *Client) getService(data *pb.BackendData) (*apiv1.Service, error) {
	client := c.Client.CoreV1().Services(c.namespace)
	serviceName := fmt.Sprintf("service-%s", data.Name)

	resp, err := client.Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		c.log.Debugw("Failed to get service",
			"name", serviceName,
			"stream-id", data.StreamID,
			"error", err)
		return nil, err
	}

	return resp, nil
}

func (c *Client) createDeployment(data *pb.BackendData) error {
	client := c.Client.AppsV1().Deployments(c.namespace)

	c.log.Infow("Creating deployment...",
		"stream-id", data.StreamID,
		"name", data.Name)

	result, err := client.Create(context.TODO(), getDeploymentSpec(data), metav1.CreateOptions{})
	if err != nil {
		c.log.Errorw("Cannot create deployment",
			"name", data.Name,
			"stream-id", data.StreamID,
			"error", err)
		return err
	}

	c.log.Infow("Created deployment",
		"stream-id", data.StreamID,
		"name", result.GetObjectMeta().GetName())

	return nil
}

func (c *Client) createService(data *pb.BackendData) error {
	client := c.Client.CoreV1().Services(c.namespace)

	c.log.Infow("Creating service...",
		"stream-id", data.StreamID,
		"name", data.Name)

	result, err := client.Create(context.TODO(), getServiceSpec(data), metav1.CreateOptions{})
	if err != nil {
		c.log.Errorw("Cannot create service",
			"name", data.Name,
			"stream-id", data.StreamID,
			"error", err)
		return err
	}

	c.log.Infow("Created service",
		"stream-id", data.StreamID,
		"name", result.GetObjectMeta().GetName())

	return nil
}

func (c *Client) createIngress(data *pb.BackendData) error {
	client := c.Client.NetworkingV1().Ingresses(c.namespace)

	c.log.Infow("Creating ingress...",
		"stream-id", data.StreamID,
		"name", data.Name)

	result, err := client.Create(context.TODO(), getIngressSpec(data), metav1.CreateOptions{})
	if err != nil {
		c.log.Errorw("Cannot create service",
			"name", data.Name,
			"stream-id", data.StreamID,
			"error", err)
		return err
	}

	c.log.Infow("Created ingress",
		"stream-id", data.StreamID,
		"name", result.GetObjectMeta().GetName())

	return nil
}

// Create creates a backend for a client request
func (c *Client) Create(data *pb.BackendData, opts ...interface{}) error {
	if err := c.createDeployment(data); err != nil {
		return fmt.Errorf("Cannot create deployment: %s, error: %s", data.Name, err)
	}

	if err := c.createService(data); err != nil {
		return fmt.Errorf("Cannot create service: %s, error: %s", data.Name, err)
	}

	if !data.OnlyTCP {
		if err := c.createIngress(data); err != nil {
			return fmt.Errorf("Cannot create ingress: %s, error: %s", data.Name, err)
		}
	}

	return nil
}

// IsReady determines if backend is ready
func (c *Client) IsReady(data *pb.BackendData, opts ...interface{}) (bool, error) {

	timeout := time.After(viper.GetDuration("server.backend-ready-timeout"))
	ticker := time.Tick(1000 * time.Millisecond)
	// Keep trying until we're time out or get a result or get an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			return false, errors.New("timed out")
		// Got a tick, we should check on checkSomething()
		case <-ticker:
			deploymentIsReady := false
			serviceIsReady := false

			deployment, err := c.getDeployment(data)
			if err != nil {
				return false, err
			}
			for _, dcondidtion := range deployment.Status.Conditions {
				if dcondidtion.Status == apiv1.ConditionFalse {
					continue
				}
				deploymentIsReady = true
			}

			service, err := c.getService(data)
			if err != nil {
				return false, err
			}
			for _, scondidtion := range service.Status.LoadBalancer.Ingress {
				if scondidtion.IP == "" {
					continue
				}
				serviceIsReady = true
			}

			if deploymentIsReady && serviceIsReady {
				return (deploymentIsReady && serviceIsReady), nil
			}
		}
	}
}

// Delete deletes the backend components created for a client request
func (c *Client) Delete(data *pb.BackendData, opts ...interface{}) error {
	serviceClient := c.Client.CoreV1().Services(c.namespace)
	deploymentClient := c.Client.AppsV1().Deployments(c.namespace)
	ingressClient := c.Client.NetworkingV1().Ingresses(c.namespace)

	serviceName := fmt.Sprintf("service-%s", data.Name)

	err := serviceClient.Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		c.log.Debugw("Failed to delete service",
			"name", serviceName,
			"stream-id", data.StreamID,
			"error", err)
		return err
	}

	err = deploymentClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
	if err != nil {
		c.log.Debugw("Failed to delete deployment",
			"name", data.Name,
			"stream-id", data.StreamID,
			"error", err)
		return err
	}

	if !data.OnlyTCP {
		err = ingressClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
		if err != nil {
			c.log.Debugw("Failed to delete ingress",
				"name", data.Name,
				"stream-id", data.StreamID,
				"error", err)
			return err
		}
	}

	return nil
}

// GetConnectionData returns data needed to establish a tunnel connection for a client
func (c *Client) GetConnectionData(data *pb.BackendData, opts ...interface{}) (*pb.Connection, error) {
	connData := &pb.Connection{}

	service, err := c.getService(data)
	if err != nil {
		return connData, err
	}

	ingress := getIngressSpec(data)

	connData.SSHHost = service.Status.LoadBalancer.Ingress[0].IP
	connData.SSHPort = service.Spec.Ports[0].Port
	connData.HTTPScheme = viper.GetBool("backend-https-enabled")
	connData.Domain = ingress.Spec.Rules[0].Host

	return connData, nil
}
