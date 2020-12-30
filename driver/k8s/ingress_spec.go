package k8s

import (
	"fmt"
	pb "sshare/protobuf"

	"github.com/spf13/viper"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ingress = &v1.Ingress{
	ObjectMeta: metav1.ObjectMeta{},
	Spec: v1.IngressSpec{
		Rules: []v1.IngressRule{
			{
				IngressRuleValue: v1.IngressRuleValue{
					HTTP: &v1.HTTPIngressRuleValue{
						Paths: []v1.HTTPIngressPath{
							{
								Path: "/",
								Backend: v1.IngressBackend{
									Service: &v1.IngressServiceBackend{
										Port: v1.ServiceBackendPort{},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

func getIngressSpec(data *pb.BackendData) *v1.Ingress {
	var prefix v1.PathType = "Prefix"
	service := getServiceSpec(data)

	ingress.ObjectMeta.Name = data.Name
	ingress.Spec.Rules[0].Host = fmt.Sprintf("%s.%s", data.Name, viper.GetString("backend-domain"))
	ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0].PathType = &prefix
	ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0].Backend.Service.Name = service.ObjectMeta.Name
	ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0].Backend.Service.Port.Name = service.Spec.Ports[1].Name

	ingress.ObjectMeta.Annotations = setIngressAnnotations(data)

	return ingress
}

func setIngressAnnotations(data *pb.BackendData) map[string]string {
	annotations := map[string]string{}

	if data.HTTPOptions.CORSEnabled {
		// refs: https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#enable-cors
		annotations["nginx.ingress.kubernetes.io/enable-cors"] = "true"
	}

	if data.HTTPOptions.HTTPSRedirect {
		// ref: https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#server-side-https-enforcement-through-redirect
		annotations["nginx.ingress.kubernetes.io/ssl-redirect"] = "true"
	}
	return annotations
}
