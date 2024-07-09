package helpers

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// DecodeCursor decodes the documents from the cursor into the provided slice.
func DecodeCursor[T any](ctx context.Context, cursor *mongo.Cursor) ([]T, error) {
	var results []T
	for cursor.Next(ctx) {
		var result T
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
