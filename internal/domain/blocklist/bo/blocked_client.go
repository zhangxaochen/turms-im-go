package bo

type BlockedClient[T any] struct {
	ID                 T
	BlockEndTimeMillis int64
}
