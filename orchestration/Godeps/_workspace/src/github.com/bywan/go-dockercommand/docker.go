package dockercommand

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

type Docker struct {
	client *docker.Client
}

func NewDocker(endpoint string) (*Docker, error) {
	endpoint = resolveDockerEndpoint(endpoint)
	client, err := docker.NewClient(resolveDockerEndpoint(endpoint))
	if err != nil {
		return nil, err
	}

	if len(os.Getenv("DOCKER_CERT_PATH")) != 0 {
		cert, err := tls.LoadX509KeyPair(os.Getenv("DOCKER_CERT_PATH")+"/cert.pem", os.Getenv("DOCKER_CERT_PATH")+"/key.pem")
		if err != nil {
			log.Fatal(err)
		}

		caCert, err := ioutil.ReadFile(os.Getenv("DOCKER_CERT_PATH") + "/ca.pem")
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
		tr := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		client.HTTPClient.Transport = tr

	}

	return &Docker{client}, nil
}

func resolveDockerEndpoint(input string) string {
	if len(input) != 0 {
		return input
	}
	if len(os.Getenv("DOCKER_HOST")) != 0 {
		return os.Getenv("DOCKER_HOST")
	}
	return "unix:///var/run/docker.sock"
}
