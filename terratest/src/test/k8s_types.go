package test

import "k8s.io/apimachinery/pkg/apis/meta/v1"

type deploymentMetadata struct {
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	v1.TypeMeta
}
