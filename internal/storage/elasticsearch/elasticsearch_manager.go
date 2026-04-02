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

// @MappedFrom putUserDoc(Long userId, String name)
func (m *ElasticsearchManager) PutUserDoc(ctx context.Context) error {
	return nil
}

// @MappedFrom putUserDocs(Collection<Long> userIds, String name)
func (m *ElasticsearchManager) PutUserDocs(ctx context.Context) error {
	return nil
}

// @MappedFrom deleteUserDoc(Long userId)
func (m *ElasticsearchManager) DeleteUserDoc(ctx context.Context) error {
	return nil
}

// @MappedFrom deleteUserDocs(Collection<Long> userIds)
func (m *ElasticsearchManager) DeleteUserDocs(ctx context.Context) error {
	return nil
}

// @MappedFrom searchUserDocs(@Nullable Integer from, @Nullable Integer size, String name, @Nullable Collection<Long> ids, boolean highlight, @Nullable String scrollId, @Nullable String keepAlive)
func (m *ElasticsearchManager) SearchUserDocs(ctx context.Context) error {
	return nil
}

// @MappedFrom putGroupDoc(Long groupId, String name)
func (m *ElasticsearchManager) PutGroupDoc(ctx context.Context) error {
	return nil
}

// @MappedFrom putGroupDocs(Collection<Long> groupIds, String name)
func (m *ElasticsearchManager) PutGroupDocs(ctx context.Context) error {
	return nil
}

// @MappedFrom deleteGroupDocs(Collection<Long> groupIds)
func (m *ElasticsearchManager) DeleteGroupDocs(ctx context.Context) error {
	return nil
}

// @MappedFrom deleteAllGroupDocs()
func (m *ElasticsearchManager) DeleteAllGroupDocs(ctx context.Context) error {
	return nil
}

// @MappedFrom searchGroupDocs(@Nullable Integer from, @Nullable Integer size, String name, @Nullable Collection<Long> ids, boolean highlight, @Nullable String scrollId, @Nullable String keepAlive)
func (m *ElasticsearchManager) SearchGroupDocs(ctx context.Context) error {
	return nil
}

// @MappedFrom deletePitForUserDocs(String scrollId)
func (m *ElasticsearchManager) DeletePitForUserDocs(ctx context.Context) error {
	return nil
}
