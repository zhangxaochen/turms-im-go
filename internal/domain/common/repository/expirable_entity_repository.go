package repository

import "time"

// ExpirableEntityRepository maps to ExpirableEntityRepository.java
// @MappedFrom ExpirableEntityRepository
type ExpirableEntityRepository struct {
	GetEntityExpireAfterSecondsFunc func() int
}

// GetEntityExpireAfterSeconds maps an abstract method hook
func (r *ExpirableEntityRepository) GetEntityExpireAfterSeconds() int {
	if r.GetEntityExpireAfterSecondsFunc != nil {
		return r.GetEntityExpireAfterSecondsFunc()
	}
	return 0
}

// @MappedFrom isExpired(long creationDate)
func (r *ExpirableEntityRepository) IsExpired(creationDate int64) bool {
	expireAfterSeconds := r.GetEntityExpireAfterSeconds()
	return expireAfterSeconds > 0 && creationDate < time.Now().UnixMilli()-int64(expireAfterSeconds)*1000
}

// @MappedFrom getEntityExpirationDate()
func (r *ExpirableEntityRepository) GetEntityExpirationDate() *time.Time {
	expireAfterSeconds := r.GetEntityExpireAfterSeconds()
	if expireAfterSeconds <= 0 {
		return nil
	}
	t := time.Now().Add(-time.Duration(expireAfterSeconds) * time.Second)
	return &t
}
