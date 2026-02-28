package opensearch

import (
	"context"
	"strings"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

// ProductIndexMapping은 Product 인덱스의 매핑 정의를 반환한다.
func ProductIndexMapping() string {
	return `{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "id":          { "type": "keyword" },
      "name":        { "type": "text", "analyzer": "standard" },
      "description": { "type": "text", "analyzer": "standard" },
      "category":    { "type": "keyword" },
      "price":       { "type": "float" },
      "in_stock":    { "type": "boolean" },
      "tags":        { "type": "keyword" },
      "created_at":  { "type": "date" }
    }
  }
}`
}

// AccessLogIndexMapping은 AccessLog 인덱스의 매핑 정의를 반환한다.
func AccessLogIndexMapping() string {
	return `{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "timestamp":        { "type": "date" },
      "method":           { "type": "keyword" },
      "endpoint":         { "type": "keyword" },
      "status_code":      { "type": "integer" },
      "response_time_ms": { "type": "float" },
      "error_message":    { "type": "text", "fields": { "keyword": { "type": "keyword" } } },
      "client_ip":        { "type": "ip" },
      "user_agent":       { "type": "text" },
      "request_body":     { "type": "text" },
      "service_name":     { "type": "keyword" }
    }
  }
}`
}

// CreateIndex는 지정된 이름과 매핑으로 인덱스를 생성한다.
func CreateIndex(ctx context.Context, client *opensearchapi.Client, indexName, mapping string) error {
	_, err := client.Indices.Create(ctx, opensearchapi.IndicesCreateReq{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	})
	return err
}

// DeleteIndex는 지정된 인덱스를 삭제한다.
func DeleteIndex(ctx context.Context, client *opensearchapi.Client, indexName string) error {
	_, err := client.Indices.Delete(ctx, opensearchapi.IndicesDeleteReq{
		Indices: []string{indexName},
	})
	return err
}

// GetMapping은 지정된 인덱스의 매핑 정보를 반환한다.
func GetMapping(ctx context.Context, client *opensearchapi.Client, indexName string) (*opensearchapi.MappingGetResp, error) {
	return client.Indices.Mapping.Get(ctx, &opensearchapi.MappingGetReq{
		Indices: []string{indexName},
	})
}
