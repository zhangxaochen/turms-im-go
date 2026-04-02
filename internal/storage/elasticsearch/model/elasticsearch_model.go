package model

type BulkRequest struct {
}

func (r *BulkRequest) Serialize() []byte {
	return nil
}

type BulkResponse struct {
}

type BulkResponseItem struct {
}

type ClosePointInTimeRequest struct {
}

type CreateIndexRequest struct {
}

type DeleteByQueryRequest struct {
}

type DeleteByQueryResponse struct {
}

type DeleteResponse struct {
}

type ErrorCause struct {
}

type ErrorResponse struct {
}

type FieldCollapse struct {
}

type HealthResponse struct {
}

type Highlight struct {
}

type IndexSettings struct {
}

type IndexSettingsAnalysis struct {
}

func (a *IndexSettingsAnalysis) Merge(other *IndexSettingsAnalysis) {
}

type PointInTimeReference struct {
}

type Property struct {
}

type Script struct {
}

type SearchRequest struct {
}

type ShardFailure struct {
}

type ShardStatistics struct {
}

type TypeMapping struct {
}

type UpdateByQueryRequest struct {
}

type UpdateByQueryResponse struct {
}
