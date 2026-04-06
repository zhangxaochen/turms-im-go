package response

// DeleteResultDTO represents the response for a delete operation.
type DeleteResultDTO struct {
	DeletedCount int64 `json:"deletedCount"`
}

// UpdateResultDTO represents the response for an update operation.
type UpdateResultDTO struct {
	UpdatedCount int64 `json:"updatedCount"`
}
