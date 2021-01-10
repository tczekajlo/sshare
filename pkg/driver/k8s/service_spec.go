package k8s

import (
	"fmt"
	pb "sshare/protobuf"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var service = &apiv1.Service{
	ObjectMeta: metav1.ObjectMeta{},
	Spec: apiv1.ServiceSpec{
		Type: apiv1.ServiceTypeLoadBalancer,
		Ports: []apiv1.ServicePort{
			{
				Name:     "sshd",
				Protocol: apiv1.ProtocolTCP,
				Port:     2222,
				TargetPort: intstr.IntOrString{
					StrVal: "sshd",
				},
			},
			{
				Name:     "remote",
				Protocol: apiv1.ProtocolTCP,
				TargetPort: intstr.IntOrString{
					StrVal: "remote",
				},
			},
		},
	},
}

func getServiceSpec(data *pb.BackendData) *apiv1.Service {
	deployment := getDeploymentSpec(data)

	service.ObjectMeta.Name = fmt.Sprintf("service-%s", data.Name)
	service.Spec.Selector = map[string]string{
		"app": data.Name,
	}
	service.Spec.Ports[1].Port = deployment.Spec.Template.Spec.Containers[0].Ports[1].ContainerPort

	return service
}
