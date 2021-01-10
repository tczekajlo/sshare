package k8s

import (
	pb "sshare/protobuf"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var deployment = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{},
	Spec: appsv1.DeploymentSpec{
		Replicas: int32Ptr(1),
		Selector: &metav1.LabelSelector{},
		Template: apiv1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{},
			Spec: apiv1.PodSpec{
				Containers: []apiv1.Container{
					{
						Name:  "openssh-server",
						Image: "tczekajlo/docker-openssh-server:latest",
						Ports: []apiv1.ContainerPort{
							{
								Name:          "sshd",
								Protocol:      apiv1.ProtocolTCP,
								ContainerPort: 2222,
							},
							{
								Name:     "remote",
								Protocol: apiv1.ProtocolTCP,
							},
						},
					},
				},
			},
		},
	},
}

func getDeploymentSpec(data *pb.BackendData) *appsv1.Deployment {
	deployment.ObjectMeta.Name = data.Name
	deployment.Spec.Selector.MatchLabels = map[string]string{
		"app": data.Name,
	}

	deployment.Spec.Template.ObjectMeta.Labels = map[string]string{
		"app": data.Name,
	}

	deployment.Spec.Template.Spec.Containers[0].Env = []apiv1.EnvVar{
		{
			Name:  "PUBLIC_KEY",
			Value: data.SshPublicKey,
		},
		{
			Name:  "USER_NAME",
			Value: "sshare",
		},
	}

	deployment.Spec.Template.Spec.Containers[0].Ports[1].ContainerPort = data.Connection.LocalPort

	return deployment
}

func int32Ptr(i int32) *int32 { return &i }
