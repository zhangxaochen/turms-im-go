package elasticsearch

import (
	"context"
)

type ElasticsearchManager struct {
	client *ElasticsearchClient
}

func NewElasticsearchManager(client *ElasticsearchClient) *ElasticsearchManager {
	return &ElasticsearchManager{
		client: client,
	}
}

func (m *ElasticsearchManager) PutUserDoc(ctx context.Context) error {
	return nil
}

func (m *ElasticsearchManager) PutUserDocs(ctx context.Context) error {
	return nil
}

func (m *ElasticsearchManager) DeleteUserDoc(ctx context.Context) error {
	return nil
}

func (m *ElasticsearchManager) DeleteUserDocs(ctx context.Context) error {
	return nil
}

func (m *ElasticsearchManager) SearchUserDocs(ctx context.Context) error {
	return nil
}

func (m *ElasticsearchManager) PutGroupDoc(ctx context.Context) error {
	return nil
}

func (m *ElasticsearchManager) PutGroupDocs(ctx context.Context) error {
	return nil
}

func (m *ElasticsearchManager) DeleteGroupDocs(ctx context.Context) error {
	return nil
}

func (m *ElasticsearchManager) DeleteAllGroupDocs(ctx context.Context) error {
	return nil
}

func (m *ElasticsearchManager) SearchGroupDocs(ctx context.Context) error {
	return nil
}

func (m *ElasticsearchManager) DeletePitForUserDocs(ctx context.Context) error {
	return nil
}
