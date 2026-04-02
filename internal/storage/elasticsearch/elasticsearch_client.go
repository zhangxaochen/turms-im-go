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

func (c *ElasticsearchClient) Healthcheck(ctx context.Context) error {
	return nil
}

func (c *ElasticsearchClient) PutIndex(ctx context.Context, request *model.CreateIndexRequest) error {
	return nil
}

func (c *ElasticsearchClient) PutDoc(ctx context.Context) error {
	return nil
}

func (c *ElasticsearchClient) DeleteDoc(ctx context.Context) error {
	return nil
}

func (c *ElasticsearchClient) DeleteByQuery(ctx context.Context, request *model.DeleteByQueryRequest) (*model.DeleteByQueryResponse, error) {
	return nil, nil
}

func (c *ElasticsearchClient) UpdateByQuery(ctx context.Context) error {
	return nil
}

func (c *ElasticsearchClient) Search(ctx context.Context) error {
	return nil
}

func (c *ElasticsearchClient) Bulk(ctx context.Context, request *model.BulkRequest) (*model.BulkResponse, error) {
	return nil, nil
}

func (c *ElasticsearchClient) DeletePit(ctx context.Context, request *model.ClosePointInTimeRequest) error {
	return nil
}
