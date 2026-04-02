package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DateRange represents a time range for MongoDB queries
type DateRange struct {
	Start *time.Time
	End   *time.Time
}

// ToBson converts DateRange to a MongoDB query filter
func (r *DateRange) ToBson() bson.M {
	if r == nil {
		return nil
	}
	m := bson.M{}
	if r.Start != nil {
		m["$gte"] = *r.Start
	}
	if r.End != nil {
		m["$lte"] = *r.End
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

// ExecuteWithSession provides a helper to execute a function within a MongoDB session (transaction if desired)
func ExecuteWithSession(ctx context.Context, client *Client, session *mongo.Session, fn func(sessCtx mongo.SessionContext, sess *mongo.Session) error) error {
	if session != nil {
		_, err := (*session).WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
			return nil, fn(sessCtx, session)
		})
		return err
	}

	sess, err := client.Client.StartSession()
	if err != nil {
		return err
	}
	defer sess.EndSession(ctx)

	return mongo.WithSession(ctx, sess, func(sessCtx mongo.SessionContext) error {
		return fn(sessCtx, &sess)
	})
}

// ExecuteWithSessionResult provides a helper to execute a function and return a result
func ExecuteWithSessionResult[T any](ctx context.Context, client *Client, session *mongo.Session, fn func(sessCtx mongo.SessionContext, sess *mongo.Session) (T, error)) (T, error) {
	if session != nil {
		res, err := (*session).WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
			return fn(sessCtx, session)
		})
		var zero T
		if err != nil {
			return zero, err
		}
		if res != nil {
			return res.(T), nil
		}
		return zero, nil
	}

	sess, err := client.Client.StartSession()
	if err != nil {
		var zero T
		return zero, err
	}
	defer sess.EndSession(ctx)

	var result T
	err = mongo.WithSession(ctx, sess, func(sessCtx mongo.SessionContext) error {
		res, err := fn(sessCtx, &sess)
		if err != nil {
			return err
		}
		result = res
		return nil
	})
	return result, err
}
