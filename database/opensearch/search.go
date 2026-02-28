package opensearch

import (
	"context"
	"strings"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

// Search는 JSON 쿼리 DSL로 검색을 수행한다.
func Search(ctx context.Context, client *opensearchapi.Client, indexName, query string) (*opensearchapi.SearchResp, error) {
	return client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{indexName},
		Body:    strings.NewReader(query),
	})
}

// MatchQuery는 match 쿼리로 전문 검색을 수행한다.
func MatchQuery(ctx context.Context, client *opensearchapi.Client, indexName, field, text string) (*opensearchapi.SearchResp, error) {
	query := `{"query":{"match":{"` + field + `":"` + text + `"}}}`
	return Search(ctx, client, indexName, query)
}

// TermQuery는 정확한 값 매칭을 수행한다.
func TermQuery(ctx context.Context, client *opensearchapi.Client, indexName, field, value string) (*opensearchapi.SearchResp, error) {
	query := `{"query":{"term":{"` + field + `":"` + value + `"}}}`
	return Search(ctx, client, indexName, query)
}

// MatchPhraseQuery는 구문(phrase) 검색을 수행한다.
func MatchPhraseQuery(ctx context.Context, client *opensearchapi.Client, indexName, field, phrase string) (*opensearchapi.SearchResp, error) {
	query := `{"query":{"match_phrase":{"` + field + `":"` + phrase + `"}}}`
	return Search(ctx, client, indexName, query)
}

// MultiMatchQuery는 여러 필드에서 검색을 수행한다.
func MultiMatchQuery(ctx context.Context, client *opensearchapi.Client, indexName, text string, fields []string) (*opensearchapi.SearchResp, error) {
	fieldsJSON := `"` + strings.Join(fields, `","`) + `"`
	query := `{"query":{"multi_match":{"query":"` + text + `","fields":[` + fieldsJSON + `]}}}`
	return Search(ctx, client, indexName, query)
}
