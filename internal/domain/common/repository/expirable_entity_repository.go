package repository

// ExpirableEntityRepository maps to ExpirableEntityRepository.java
// @MappedFrom ExpirableEntityRepository
type ExpirableEntityRepository struct {
}

// @MappedFrom isExpired(long creationDate)
func (r *ExpirableEntityRepository) IsExpired() {
}

// @MappedFrom getEntityExpirationDate()
func (r *ExpirableEntityRepository) GetEntityExpirationDate() {
}
