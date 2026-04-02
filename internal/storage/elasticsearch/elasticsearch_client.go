package elasticsearch

import (
	"context"

	"im.turms/server/internal/storage/elasticsearch/model"
)

type ElasticsearchClient struct {
}

func NewElasticsearchClient() *ElasticsearchClient {
	return &ElasticsearchClient{}
}

// @MappedFrom healthcheck()
func (c *ElasticsearchClient) Healthcheck(ctx context.Context) error {
	return nil
}

// @MappedFrom putIndex(String index, CreateIndexRequest request)
func (c *ElasticsearchClient) PutIndex(ctx context.Context, request *model.CreateIndexRequest) error {
	return nil
}

// @MappedFrom putDoc(String index, String id, Supplier<ByteBuf> payloadSupplier)
func (c *ElasticsearchClient) PutDoc(ctx context.Context) error {
	return nil
}

// @MappedFrom deleteDoc(String index, String id)
func (c *ElasticsearchClient) DeleteDoc(ctx context.Context) error {
	return nil
}

// @MappedFrom deleteByQuery(String index, DeleteByQueryRequest request)
func (c *ElasticsearchClient) DeleteByQuery(ctx context.Context, request *model.DeleteByQueryRequest) (*model.DeleteByQueryResponse, error) {
	return nil, nil
}

// @MappedFrom updateByQuery(String index, UpdateByQueryRequest request)
func (c *ElasticsearchClient) UpdateByQuery(ctx context.Context) error {
	return nil
}

// @MappedFrom search(String baseDn, Scope scope, DerefAliases derefAliases, int sizeLimit, int timeLimit, boolean typeOnly, List<String> attributes, String filter)
// @MappedFrom search(String index, SearchRequest request, ObjectReader reader)
func (c *ElasticsearchClient) Search(ctx context.Context) error {
	return nil
}

// @MappedFrom bulk(BulkRequest request)
func (c *ElasticsearchClient) Bulk(ctx context.Context, request *model.BulkRequest) (*model.BulkResponse, error) {
	return nil, nil
}

// @MappedFrom deletePit(String scrollId)
func (c *ElasticsearchClient) DeletePit(ctx context.Context, request *model.ClosePointInTimeRequest) error {
	return nil
}
