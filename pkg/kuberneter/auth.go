package kuberneter

import (
	"os"

	b64 "encoding/base64"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// K8sConfig saves the config for k8s auth
type K8sConfig struct {
	Config *rest.Config
}

func getK8sClient() (client.Client, error) {
	// authentication through k8s config file
	if len(os.Getenv("K8S_CONFIG_FILE")) > 0 {
		config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("K8S_CONFIG_FILE"))

		if err != nil {
			log.Error().Msgf("failed to get k8s client through k8s config file: %s", err)
			return nil, err
		}

		log.Info().Msg("auth to k8s API through k8s config file")

		return client.New(config, client.Options{})
	}

	c := &K8sConfig{
		Config: &rest.Config{},
	}

	// authentication through k8s service account token or k8s client certificate
	if len(os.Getenv("K8S_HOST")) > 0 && c.hasCertificateAuthority() {
		c.Config.Host = os.Getenv("K8S_HOST")

		// authentication through k8s service account token
		if c.hasServiceAccountToken() {
			log.Info().Msg("auth to k8s API through k8s service account token")
			return client.New(c.Config, client.Options{})
		}

		// authentication through k8s client certificate
		if c.hasClientCertificate() {
			log.Info().Msg("auth to k8s API through k8s client certificate")
			return client.New(c.Config, client.Options{})
		}
	}

	log.Error().Msg("failed to get k8s client. check the k8s cluster auth information")
	return nil, errors.New("failed to get k8s client")
}

func (c *K8sConfig) hasCertificateAuthority() bool {
	if len(os.Getenv("K8S_CA_FILE")) > 0 {
		c.Config.TLSClientConfig.CAFile = os.Getenv("K8S_CA_FILE")
		return true
	}

	if len(os.Getenv("K8S_CA_DATA")) > 0 {
		caDataDecoded, err := b64.StdEncoding.DecodeString(os.Getenv("K8S_CA_DATA"))
		if err != nil {
			log.Error().Msgf("failed to decode K8S_CA_DATA: %s", err)
			return false
		}
		c.Config.TLSClientConfig.CAData = caDataDecoded
		return true
	}

	return false
}

func (c *K8sConfig) hasServiceAccountToken() bool {
	if len(os.Getenv("K8S_SA_TOKEN_FILE")) > 0 {
		c.Config.BearerTokenFile = os.Getenv("K8S_SA_TOKEN_FILE")
		return true
	}

	if len(os.Getenv("K8S_SA_TOKEN_DATA")) > 0 {
		c.Config.BearerToken = os.Getenv("K8S_SA_TOKEN_DATA")
		return true
	}

	return false
}

func (c *K8sConfig) hasClientCertificate() bool {
	hasCert := false

	if len(os.Getenv("K8S_CERT_FILE")) > 0 {
		c.Config.TLSClientConfig.CertFile = os.Getenv("K8S_CERT_FILE")
		hasCert = true
	}

	if len(os.Getenv("K8S_CERT_DATA")) > 0 {
		certDataDecoded, err := b64.StdEncoding.DecodeString(os.Getenv("K8S_CERT_DATA"))
		if err != nil {
			log.Error().Msgf("failed to decode K8S_CERT_DATA: %s", err)
			return false
		}
		c.Config.TLSClientConfig.CertData = certDataDecoded
		hasCert = true
	}

	if hasCert {
		if len(os.Getenv("K8S_KEY_FILE")) > 0 {
			c.Config.TLSClientConfig.KeyFile = os.Getenv("K8S_KEY_FILE")
			return true
		}

		if len(os.Getenv("K8S_KEY_DATA")) > 0 {
			keyDataDecoded, err := b64.StdEncoding.DecodeString(os.Getenv("K8S_KEY_DATA"))
			if err != nil {
				log.Error().Msgf("failed to decode K8S_KEY_DATA: %s", err)
				return false
			}
			c.Config.TLSClientConfig.KeyData = keyDataDecoded
			return true
		}
	}

	return false
}
