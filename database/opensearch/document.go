package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

// IndexDocument는 단건 문서를 색인한다.
func IndexDocument(ctx context.Context, client *opensearchapi.Client, indexName, docID string, doc any) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("문서 직렬화 실패: %w", err)
	}
	_, err = client.Index(ctx, opensearchapi.IndexReq{
		Index:      indexName,
		DocumentID: docID,
		Body:       strings.NewReader(string(body)),
		Params:     opensearchapi.IndexParams{Refresh: "true"},
	})
	return err
}

// BulkIndex는 여러 문서를 벌크로 색인한다.
func BulkIndex(ctx context.Context, client *opensearchapi.Client, indexName string, docs map[string]any) error {
	var sb strings.Builder
	for id, doc := range docs {
		meta := fmt.Sprintf(`{"index":{"_index":"%s","_id":"%s"}}`, indexName, id)
		sb.WriteString(meta)
		sb.WriteString("\n")
		body, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("문서 직렬화 실패 (id=%s): %w", id, err)
		}
		sb.Write(body)
		sb.WriteString("\n")
	}

	resp, err := client.Bulk(ctx, opensearchapi.BulkReq{
		Body:   strings.NewReader(sb.String()),
		Params: opensearchapi.BulkParams{Refresh: "true"},
	})
	if err != nil {
		return err
	}
	if resp.Errors {
		return fmt.Errorf("벌크 색인 중 에러 발생")
	}
	return nil
}

// GetDocument는 문서 ID로 단건 조회한다.
func GetDocument(ctx context.Context, client *opensearchapi.Client, indexName, docID string) (*opensearchapi.DocumentGetResp, error) {
	return client.Document.Get(ctx, opensearchapi.DocumentGetReq{
		Index:      indexName,
		DocumentID: docID,
	})
}

// UpdateDocument는 문서의 일부 필드를 수정한다.
func UpdateDocument(ctx context.Context, client *opensearchapi.Client, indexName, docID string, fields map[string]any) error {
	doc := map[string]any{"doc": fields}
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("업데이트 데이터 직렬화 실패: %w", err)
	}
	_, err = client.Update(ctx, opensearchapi.UpdateReq{
		Index:      indexName,
		DocumentID: docID,
		Body:       strings.NewReader(string(body)),
		Params:     opensearchapi.UpdateParams{Refresh: "true"},
	})
	return err
}

// DeleteDocument는 문서 ID로 단건 삭제한다.
func DeleteDocument(ctx context.Context, client *opensearchapi.Client, indexName, docID string) error {
	_, err := client.Document.Delete(ctx, opensearchapi.DocumentDeleteReq{
		Index:      indexName,
		DocumentID: docID,
		Params: opensearchapi.DocumentDeleteParams{
			Refresh: "true",
		},
	})
	return err
}

// DeleteByQuery는 쿼리 조건에 맞는 문서를 삭제한다.
func DeleteByQuery(ctx context.Context, client *opensearchapi.Client, indexName, query string) (*opensearchapi.DocumentDeleteByQueryResp, error) {
	return client.Document.DeleteByQuery(ctx, opensearchapi.DocumentDeleteByQueryReq{
		Indices: []string{indexName},
		Body:    strings.NewReader(query),
		Params:  opensearchapi.DocumentDeleteByQueryParams{Refresh: opensearchapi.ToPointer(true)},
	})
}
