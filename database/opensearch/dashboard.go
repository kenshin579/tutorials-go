package opensearch

import (
	"context"
	"fmt"
	"strings"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

// TopErrors는 에러 메시지 기준 Top N 에러를 조회한다.
func TopErrors(ctx context.Context, client *opensearchapi.Client, indexName string, topN int) (*opensearchapi.SearchResp, error) {
	query := fmt.Sprintf(`{
  "size": 0,
  "query": {
    "range": { "status_code": { "gte": 400 } }
  },
  "aggs": {
    "top_errors": {
      "terms": { "field": "error_message.keyword", "size": %d }
    }
  }
}`, topN)
	return client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{indexName},
		Body:    strings.NewReader(query),
	})
}

// StatusCodeDistribution은 HTTP 상태 코드 분포를 조회한다.
func StatusCodeDistribution(ctx context.Context, client *opensearchapi.Client, indexName string) (*opensearchapi.SearchResp, error) {
	query := `{
  "size": 0,
  "aggs": {
    "status_codes": {
      "terms": { "field": "status_code", "size": 20 }
    }
  }
}`
	return client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{indexName},
		Body:    strings.NewReader(query),
	})
}

// ErrorRateOverTime은 시간대별 에러(4xx, 5xx) 발생 추이를 조회한다.
func ErrorRateOverTime(ctx context.Context, client *opensearchapi.Client, indexName, interval string) (*opensearchapi.SearchResp, error) {
	query := `{
  "size": 0,
  "query": {
    "range": { "status_code": { "gte": 400 } }
  },
  "aggs": {
    "errors_over_time": {
      "date_histogram": { "field": "timestamp", "calendar_interval": "` + interval + `" }
    }
  }
}`
	return client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{indexName},
		Body:    strings.NewReader(query),
	})
}

// RequestCountByEndpoint는 엔드포인트별 호출 횟수를 조회한다.
func RequestCountByEndpoint(ctx context.Context, client *opensearchapi.Client, indexName string, topN int) (*opensearchapi.SearchResp, error) {
	query := fmt.Sprintf(`{
  "size": 0,
  "aggs": {
    "by_endpoint": {
      "terms": { "field": "endpoint", "size": %d }
    }
  }
}`, topN)
	return client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{indexName},
		Body:    strings.NewReader(query),
	})
}

// SlowestEndpoints는 평균 응답 시간이 가장 느린 엔드포인트 Top N을 조회한다.
func SlowestEndpoints(ctx context.Context, client *opensearchapi.Client, indexName string, topN int) (*opensearchapi.SearchResp, error) {
	query := fmt.Sprintf(`{
  "size": 0,
  "aggs": {
    "by_endpoint": {
      "terms": { "field": "endpoint", "size": %d, "order": { "avg_response_time": "desc" } },
      "aggs": {
        "avg_response_time": {
          "avg": { "field": "response_time_ms" }
        }
      }
    }
  }
}`, topN)
	return client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{indexName},
		Body:    strings.NewReader(query),
	})
}

// PercentileResponseTime은 응답 시간의 퍼센타일(P50, P95, P99)을 조회한다.
func PercentileResponseTime(ctx context.Context, client *opensearchapi.Client, indexName string) (*opensearchapi.SearchResp, error) {
	query := `{
  "size": 0,
  "aggs": {
    "response_time_percentiles": {
      "percentiles": {
        "field": "response_time_ms",
        "percents": [50, 95, 99]
      }
    }
  }
}`
	return client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{indexName},
		Body:    strings.NewReader(query),
	})
}
