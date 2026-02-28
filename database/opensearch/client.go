package opensearch

import (
	"context"
	"fmt"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NewOpenSearchClient는 주어진 endpoint로 OpenSearch 클라이언트를 생성한다.
func NewOpenSearchClient(endpoint string) (*opensearchapi.Client, error) {
	client, err := opensearchapi.NewClient(opensearchapi.Config{
		Client: opensearch.Config{
			Addresses: []string{endpoint},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("opensearch client 생성 실패: %w", err)
	}
	return client, nil
}

// StartOpenSearchContainer는 testcontainers로 OpenSearch 컨테이너를 시작하고 endpoint를 반환한다.
func StartOpenSearchContainer(ctx context.Context) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        "opensearchproject/opensearch:2.11.1",
		ExposedPorts: []string{"9200/tcp"},
		Env: map[string]string{
			"discovery.type":          "single-node",
			"DISABLE_SECURITY_PLUGIN": "true",
			"OPENSEARCH_JAVA_OPTS":    "-Xms512m -Xmx512m",
		},
		WaitingFor: wait.ForHTTP("/").WithPort("9200/tcp"),
	}

	opensearchC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("opensearch 컨테이너 시작 실패: %w", err)
	}

	endpoint, err := opensearchC.Endpoint(ctx, "http")
	if err != nil {
		return nil, "", fmt.Errorf("endpoint 가져오기 실패: %w", err)
	}

	return opensearchC, endpoint, nil
}
