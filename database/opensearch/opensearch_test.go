package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

var (
	client      *opensearchapi.Client
	opensearchC testcontainers.Container
)

const (
	productIndex   = "test-products"
	accessLogIndex = "test-access-logs"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var endpoint string
	var err error
	opensearchC, endpoint, err = StartOpenSearchContainer(ctx)
	if err != nil {
		fmt.Printf("OpenSearch 컨테이너 시작 실패: %v\n", err)
		os.Exit(1)
	}

	client, err = NewOpenSearchClient(endpoint)
	if err != nil {
		fmt.Printf("OpenSearch 클라이언트 생성 실패: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	if opensearchC != nil {
		_ = opensearchC.Terminate(ctx)
	}
	os.Exit(code)
}

// --- Helper ---

func setupProductIndex(t *testing.T) {
	t.Helper()
	err := CreateIndex(context.Background(), client, productIndex, ProductIndexMapping())
	require.NoError(t, err)
}

func teardownProductIndex(t *testing.T) {
	t.Helper()
	_ = DeleteIndex(context.Background(), client, productIndex)
}

func setupAccessLogIndex(t *testing.T) {
	t.Helper()
	err := CreateIndex(context.Background(), client, accessLogIndex, AccessLogIndexMapping())
	require.NoError(t, err)
}

func teardownAccessLogIndex(t *testing.T) {
	t.Helper()
	_ = DeleteIndex(context.Background(), client, accessLogIndex)
}

func loadProducts(t *testing.T) []Product {
	t.Helper()
	data, err := os.ReadFile("testdata/products.json")
	require.NoError(t, err)
	var products []Product
	require.NoError(t, json.Unmarshal(data, &products))
	return products
}

func loadAccessLogs(t *testing.T) []AccessLog {
	t.Helper()
	data, err := os.ReadFile("testdata/access_logs.json")
	require.NoError(t, err)
	var logs []AccessLog
	require.NoError(t, json.Unmarshal(data, &logs))
	return logs
}

func indexProducts(t *testing.T, products []Product) {
	t.Helper()
	ctx := context.Background()
	docs := make(map[string]any)
	for _, p := range products {
		docs[p.ID] = p
	}
	require.NoError(t, BulkIndex(ctx, client, productIndex, docs))
}

func indexAccessLogs(t *testing.T, logs []AccessLog) {
	t.Helper()
	ctx := context.Background()
	docs := make(map[string]any)
	for i, l := range logs {
		docs[fmt.Sprintf("log-%d", i+1)] = l
	}
	require.NoError(t, BulkIndex(ctx, client, accessLogIndex, docs))
}

// --- 인덱스 관리 테스트 ---

func TestCreateIndex(t *testing.T) {
	ctx := context.Background()
	err := CreateIndex(ctx, client, productIndex, ProductIndexMapping())
	assert.NoError(t, err)
	defer teardownProductIndex(t)
}

func TestDeleteIndex(t *testing.T) {
	setupProductIndex(t)

	ctx := context.Background()
	err := DeleteIndex(ctx, client, productIndex)
	assert.NoError(t, err)
}

func TestGetMapping(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	ctx := context.Background()
	resp, err := GetMapping(ctx, client, productIndex)
	require.NoError(t, err)

	indexMapping, exists := resp.Indices[productIndex]
	assert.True(t, exists)
	assert.NotEmpty(t, indexMapping.Mappings)
}

// --- 문서 CRUD 테스트 ---

func TestIndexDocument(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	ctx := context.Background()
	product := Product{
		ID:          "test-1",
		Name:        "Test Product",
		Description: "A test product",
		Category:    "test",
		Price:       9.99,
		InStock:     true,
		Tags:        []string{"test"},
		CreatedAt:   time.Now(),
	}

	err := IndexDocument(ctx, client, productIndex, product.ID, product)
	assert.NoError(t, err)

	resp, err := GetDocument(ctx, client, productIndex, product.ID)
	require.NoError(t, err)
	assert.True(t, resp.Found)
}

func TestBulkIndex(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	searchResp, err := Search(context.Background(), client, productIndex, `{"query":{"match_all":{}}}`)
	require.NoError(t, err)
	assert.Equal(t, len(products), searchResp.Hits.Total.Value)
}

func TestGetDocument(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	resp, err := GetDocument(ctx, client, productIndex, "1")
	require.NoError(t, err)
	assert.True(t, resp.Found)

	var p Product
	require.NoError(t, json.Unmarshal(resp.Source, &p))
	assert.Equal(t, "Go Programming Language", p.Name)
}

func TestUpdateDocument(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	err := UpdateDocument(ctx, client, productIndex, "1", map[string]any{
		"price": 29.99,
	})
	assert.NoError(t, err)

	resp, err := GetDocument(ctx, client, productIndex, "1")
	require.NoError(t, err)

	var p Product
	require.NoError(t, json.Unmarshal(resp.Source, &p))
	assert.Equal(t, 29.99, p.Price)
}

func TestDeleteDocument(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	err := DeleteDocument(ctx, client, productIndex, "1")
	assert.NoError(t, err)

	// 삭제된 문서를 조회하면 404 에러가 반환된다
	_, err = GetDocument(ctx, client, productIndex, "1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestDeleteByQuery(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	query := `{"query":{"term":{"category":"books"}}}`
	resp, err := DeleteByQuery(ctx, client, productIndex, query)
	require.NoError(t, err)
	assert.Greater(t, resp.Deleted, 0)
}

// --- 검색 쿼리 테스트 ---

func TestMatchQuery(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	resp, err := MatchQuery(ctx, client, productIndex, "name", "Go Programming")
	require.NoError(t, err)
	assert.Greater(t, resp.Hits.Total.Value, 0)
}

func TestTermQuery(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	resp, err := TermQuery(ctx, client, productIndex, "category", "books")
	require.NoError(t, err)
	assert.Greater(t, resp.Hits.Total.Value, 0)
}

func TestMatchPhraseQuery(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	resp, err := MatchPhraseQuery(ctx, client, productIndex, "description", "Go programming language")
	require.NoError(t, err)
	assert.Greater(t, resp.Hits.Total.Value, 0)
}

func TestMultiMatchQuery(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	resp, err := MultiMatchQuery(ctx, client, productIndex, "programming", []string{"name", "description"})
	require.NoError(t, err)
	assert.Greater(t, resp.Hits.Total.Value, 0)
}

func TestBoolQuery(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	query := `{
  "query": {
    "bool": {
      "must": [
        { "match": { "category": "books" } }
      ],
      "filter": [
        { "range": { "price": { "lte": 40 } } }
      ]
    }
  }
}`
	resp, err := Search(ctx, client, productIndex, query)
	require.NoError(t, err)
	assert.Greater(t, resp.Hits.Total.Value, 0)

	for _, hit := range resp.Hits.Hits {
		var p Product
		require.NoError(t, json.Unmarshal(hit.Source, &p))
		assert.Equal(t, "books", p.Category)
		assert.LessOrEqual(t, p.Price, 40.0)
	}
}

func TestRangeQuery(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	query := `{
  "query": {
    "range": {
      "price": {
        "gte": 40,
        "lte": 100
      }
    }
  }
}`
	resp, err := Search(ctx, client, productIndex, query)
	require.NoError(t, err)
	assert.Greater(t, resp.Hits.Total.Value, 0)

	for _, hit := range resp.Hits.Hits {
		var p Product
		require.NoError(t, json.Unmarshal(hit.Source, &p))
		assert.GreaterOrEqual(t, p.Price, 40.0)
		assert.LessOrEqual(t, p.Price, 100.0)
	}
}

func TestSearchWithSort(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	query := `{
  "query": { "match_all": {} },
  "sort": [
    { "price": { "order": "asc" } }
  ]
}`
	resp, err := Search(ctx, client, productIndex, query)
	require.NoError(t, err)
	assert.Greater(t, len(resp.Hits.Hits), 1)

	var prevPrice float64
	for _, hit := range resp.Hits.Hits {
		var p Product
		require.NoError(t, json.Unmarshal(hit.Source, &p))
		assert.GreaterOrEqual(t, p.Price, prevPrice)
		prevPrice = p.Price
	}
}

func TestSearchWithPagination(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	// 첫 번째 페이지: 3건
	query := `{
  "query": { "match_all": {} },
  "from": 0,
  "size": 3,
  "sort": [{ "price": { "order": "asc" } }]
}`
	page1, err := Search(ctx, client, productIndex, query)
	require.NoError(t, err)
	assert.Equal(t, 3, len(page1.Hits.Hits))

	// 두 번째 페이지: 3건
	query = `{
  "query": { "match_all": {} },
  "from": 3,
  "size": 3,
  "sort": [{ "price": { "order": "asc" } }]
}`
	page2, err := Search(ctx, client, productIndex, query)
	require.NoError(t, err)
	assert.Equal(t, 3, len(page2.Hits.Hits))
	assert.NotEqual(t, page1.Hits.Hits[0].ID, page2.Hits.Hits[0].ID)
}

func TestSearchWithHighlight(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	query := `{
  "query": {
    "match": { "description": "programming" }
  },
  "highlight": {
    "fields": {
      "description": {}
    }
  }
}`
	resp, err := Search(ctx, client, productIndex, query)
	require.NoError(t, err)
	assert.Greater(t, resp.Hits.Total.Value, 0)

	for _, hit := range resp.Hits.Hits {
		highlights, exists := hit.Highlight["description"]
		assert.True(t, exists)
		assert.Greater(t, len(highlights), 0)
	}
}

// --- Aggregation 테스트 ---

func TestAvgAggregation(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	resp, err := AvgAggregation(ctx, client, productIndex, "price")
	require.NoError(t, err)

	var aggs map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(resp.Aggregations, &aggs))

	var avgResult struct {
		Value float64 `json:"value"`
	}
	require.NoError(t, json.Unmarshal(aggs["avg_value"], &avgResult))
	assert.Greater(t, avgResult.Value, 0.0)
}

func TestTermsAggregation(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	resp, err := TermsAggregation(ctx, client, productIndex, "category", 10)
	require.NoError(t, err)

	var aggs map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(resp.Aggregations, &aggs))

	var termsResult struct {
		Buckets []struct {
			Key      string `json:"key"`
			DocCount int    `json:"doc_count"`
		} `json:"buckets"`
	}
	require.NoError(t, json.Unmarshal(aggs["group_by"], &termsResult))
	assert.Greater(t, len(termsResult.Buckets), 0)
}

func TestDateHistogramAggregation(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	resp, err := DateHistogramAggregation(ctx, client, productIndex, "created_at", "month")
	require.NoError(t, err)

	var aggs map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(resp.Aggregations, &aggs))

	var histResult struct {
		Buckets []struct {
			KeyAsString string `json:"key_as_string"`
			DocCount    int    `json:"doc_count"`
		} `json:"buckets"`
	}
	require.NoError(t, json.Unmarshal(aggs["over_time"], &histResult))
	assert.Greater(t, len(histResult.Buckets), 0)
}

func TestNestedAggregation(t *testing.T) {
	setupProductIndex(t)
	defer teardownProductIndex(t)

	products := loadProducts(t)
	indexProducts(t, products)

	ctx := context.Background()
	resp, err := NestedAggregation(ctx, client, productIndex, "category", "price")
	require.NoError(t, err)

	var aggs map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(resp.Aggregations, &aggs))

	var nestedResult struct {
		Buckets []struct {
			Key       string `json:"key"`
			DocCount  int    `json:"doc_count"`
			AvgMetric struct {
				Value float64 `json:"value"`
			} `json:"avg_metric"`
		} `json:"buckets"`
	}
	require.NoError(t, json.Unmarshal(aggs["group_by"], &nestedResult))
	assert.Greater(t, len(nestedResult.Buckets), 0)

	for _, bucket := range nestedResult.Buckets {
		assert.Greater(t, bucket.AvgMetric.Value, 0.0)
	}
}

// --- Dashboard 로그 분석 테스트 ---

func TestTopErrors(t *testing.T) {
	setupAccessLogIndex(t)
	defer teardownAccessLogIndex(t)

	logs := loadAccessLogs(t)
	indexAccessLogs(t, logs)

	ctx := context.Background()
	resp, err := TopErrors(ctx, client, accessLogIndex, 5)
	require.NoError(t, err)

	var aggs map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(resp.Aggregations, &aggs))

	var termsResult struct {
		Buckets []struct {
			Key      string `json:"key"`
			DocCount int    `json:"doc_count"`
		} `json:"buckets"`
	}
	require.NoError(t, json.Unmarshal(aggs["top_errors"], &termsResult))
	assert.Greater(t, len(termsResult.Buckets), 0)
	t.Logf("Top errors: %+v", termsResult.Buckets)
}

func TestStatusCodeDistribution(t *testing.T) {
	setupAccessLogIndex(t)
	defer teardownAccessLogIndex(t)

	logs := loadAccessLogs(t)
	indexAccessLogs(t, logs)

	ctx := context.Background()
	resp, err := StatusCodeDistribution(ctx, client, accessLogIndex)
	require.NoError(t, err)

	var aggs map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(resp.Aggregations, &aggs))

	var termsResult struct {
		Buckets []struct {
			Key      int `json:"key"`
			DocCount int `json:"doc_count"`
		} `json:"buckets"`
	}
	require.NoError(t, json.Unmarshal(aggs["status_codes"], &termsResult))
	assert.Greater(t, len(termsResult.Buckets), 0)
	t.Logf("Status code distribution: %+v", termsResult.Buckets)
}

func TestErrorRateOverTime(t *testing.T) {
	setupAccessLogIndex(t)
	defer teardownAccessLogIndex(t)

	logs := loadAccessLogs(t)
	indexAccessLogs(t, logs)

	ctx := context.Background()
	resp, err := ErrorRateOverTime(ctx, client, accessLogIndex, "hour")
	require.NoError(t, err)

	var aggs map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(resp.Aggregations, &aggs))

	var histResult struct {
		Buckets []struct {
			KeyAsString string `json:"key_as_string"`
			DocCount    int    `json:"doc_count"`
		} `json:"buckets"`
	}
	require.NoError(t, json.Unmarshal(aggs["errors_over_time"], &histResult))
	assert.Greater(t, len(histResult.Buckets), 0)
	t.Logf("Error rate over time: %+v", histResult.Buckets)
}

func TestRequestCountByEndpoint(t *testing.T) {
	setupAccessLogIndex(t)
	defer teardownAccessLogIndex(t)

	logs := loadAccessLogs(t)
	indexAccessLogs(t, logs)

	ctx := context.Background()
	resp, err := RequestCountByEndpoint(ctx, client, accessLogIndex, 10)
	require.NoError(t, err)

	var aggs map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(resp.Aggregations, &aggs))

	var termsResult struct {
		Buckets []struct {
			Key      string `json:"key"`
			DocCount int    `json:"doc_count"`
		} `json:"buckets"`
	}
	require.NoError(t, json.Unmarshal(aggs["by_endpoint"], &termsResult))
	assert.Greater(t, len(termsResult.Buckets), 0)
	t.Logf("Request count by endpoint: %+v", termsResult.Buckets)
}

func TestSlowestEndpoints(t *testing.T) {
	setupAccessLogIndex(t)
	defer teardownAccessLogIndex(t)

	logs := loadAccessLogs(t)
	indexAccessLogs(t, logs)

	ctx := context.Background()
	resp, err := SlowestEndpoints(ctx, client, accessLogIndex, 5)
	require.NoError(t, err)

	var aggs map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(resp.Aggregations, &aggs))

	var termsResult struct {
		Buckets []struct {
			Key             string `json:"key"`
			DocCount        int    `json:"doc_count"`
			AvgResponseTime struct {
				Value float64 `json:"value"`
			} `json:"avg_response_time"`
		} `json:"buckets"`
	}
	require.NoError(t, json.Unmarshal(aggs["by_endpoint"], &termsResult))
	assert.Greater(t, len(termsResult.Buckets), 0)
	t.Logf("Slowest endpoints: %+v", termsResult.Buckets)
}

func TestPercentileResponseTime(t *testing.T) {
	setupAccessLogIndex(t)
	defer teardownAccessLogIndex(t)

	logs := loadAccessLogs(t)
	indexAccessLogs(t, logs)

	ctx := context.Background()
	resp, err := PercentileResponseTime(ctx, client, accessLogIndex)
	require.NoError(t, err)

	var aggs map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(resp.Aggregations, &aggs))

	var percResult struct {
		Values map[string]float64 `json:"values"`
	}
	require.NoError(t, json.Unmarshal(aggs["response_time_percentiles"], &percResult))
	assert.Contains(t, percResult.Values, "50.0")
	assert.Contains(t, percResult.Values, "95.0")
	assert.Contains(t, percResult.Values, "99.0")
	t.Logf("Percentile response times: %+v", percResult.Values)
}
