package opensearch

import (
	"context"
	"fmt"
	"strings"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

// Aggregate는 Aggregation 쿼리를 수행한다.
func Aggregate(ctx context.Context, client *opensearchapi.Client, indexName, query string) (*opensearchapi.SearchResp, error) {
	return client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{indexName},
		Body:    strings.NewReader(query),
	})
}

// AvgAggregation은 특정 필드의 평균값을 계산한다.
func AvgAggregation(ctx context.Context, client *opensearchapi.Client, indexName, field string) (*opensearchapi.SearchResp, error) {
	query := `{
  "size": 0,
  "aggs": {
    "avg_value": {
      "avg": { "field": "` + field + `" }
    }
  }
}`
	return Aggregate(ctx, client, indexName, query)
}

// TermsAggregation은 특정 필드의 그룹별 집계를 수행한다.
func TermsAggregation(ctx context.Context, client *opensearchapi.Client, indexName, field string, size int) (*opensearchapi.SearchResp, error) {
	query := `{
  "size": 0,
  "aggs": {
    "group_by": {
      "terms": { "field": "` + field + `", "size": ` + fmt.Sprintf("%d", size) + ` }
    }
  }
}`
	return Aggregate(ctx, client, indexName, query)
}

// DateHistogramAggregation은 날짜 필드를 기준으로 시계열 집계를 수행한다.
func DateHistogramAggregation(ctx context.Context, client *opensearchapi.Client, indexName, field, interval string) (*opensearchapi.SearchResp, error) {
	query := `{
  "size": 0,
  "aggs": {
    "over_time": {
      "date_histogram": { "field": "` + field + `", "calendar_interval": "` + interval + `" }
    }
  }
}`
	return Aggregate(ctx, client, indexName, query)
}

// NestedAggregation은 Bucket 안에 Metric을 중첩하여 집계한다.
// 예: 카테고리별 평균 가격
func NestedAggregation(ctx context.Context, client *opensearchapi.Client, indexName, bucketField, metricField string) (*opensearchapi.SearchResp, error) {
	query := `{
  "size": 0,
  "aggs": {
    "group_by": {
      "terms": { "field": "` + bucketField + `", "size": 10 },
      "aggs": {
        "avg_metric": {
          "avg": { "field": "` + metricField + `" }
        }
      }
    }
  }
}`
	return Aggregate(ctx, client, indexName, query)
}

