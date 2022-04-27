package metadatabp

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type BaseMetadata string

const (
	BaseplateK8sNodeName  BaseMetadata = "baseplateK8sNodeName"
	BaseplateK8sNodeIP    BaseMetadata = "baseplateK8sNodeIP"
	BaseplateK8sPodName   BaseMetadata = "baseplateK8sPodName"
	BaseplateK8sPodIP     BaseMetadata = "baseplateK8sPodIP"
	BaseplateK8sNamespace BaseMetadata = "baseplateK8sNamespace"
)

var baseMetadataVariables = map[BaseMetadata]string{
	BaseplateK8sNodeName:  "BASEPLATE_K8S_METADATA_NODE_NAME",
	BaseplateK8sNodeIP:    "BASEPLATE_K8S_METADATA_NODE_IP",
	BaseplateK8sPodName:   "BASEPLATE_K8S_METADATA_POD_NAME",
	BaseplateK8sPodIP:     "BASEPLATE_K8S_METADATA_POD_IP",
	BaseplateK8sNamespace: "BASEPLATE_K8S_METADATA_NAMESPACE",
}

type Config struct {
	BaseK8sMetadata map[BaseMetadata]string
	k8sClient       K8SFetcher
}

// K8SFetcher is an interface for fetching data from the k8s api.
// The interface is solely here to make unit testing with the k8s client easier.
type K8SFetcher interface {
	GetPodStatus(ctx context.Context, namespace, podName string) (string, error)
}

// New returns a new instance of metadata.
func New(options ...func(*Config) error) (*Config, error) {
	config := &Config{}

	baseK8sMetadata, err := fetchBaseMetadata()
	if err != nil {
		return nil, err
	}

	config.BaseK8sMetadata = baseK8sMetadata

	for _, option := range options {
		err := option(config)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

// WithK8SClient adds the k8s client to config.
func WithK8SClient() func(*Config) error {
	return func(c *Config) error {
		client, err := NewK8SClient()
		c.k8sClient = client
		return err
	}
}

// GetBaseMetadata returns base k8s metadata based on the key provided.
func (c *Config) GetBaseMetadata(key BaseMetadata) string {
	return c.BaseK8sMetadata[key]
}

// fetchBaseMetadata fetches the base k8s metadata from the environment
// this base metadata is needed for fetching further metadata via the k8s api.
func fetchBaseMetadata() (map[BaseMetadata]string, error) {
	baseK8sMetadata := make(map[BaseMetadata]string)
	for k, v := range baseMetadataVariables {
		value, exists := os.LookupEnv(v)
		if !exists || strings.TrimSpace(value) == "" {
			return nil, fmt.Errorf("metadatapb:%s base k8s metadata value not present", v)
		}
		baseK8sMetadata[k] = value
	}
	return baseK8sMetadata, nil
}
