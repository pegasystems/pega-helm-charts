package addons

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentMetadata struct {
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	metav1.TypeMeta
}
