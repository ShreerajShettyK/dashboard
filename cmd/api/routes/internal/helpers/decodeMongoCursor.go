// package helpers

// import (
// 	"context"

// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // DecodeCursor decodes the documents from the cursor into the provided slice.
// func DecodeCursor[T any](ctx context.Context, cursor *mongo.Cursor) ([]T, error) {
// 	var results []T
// 	for cursor.Next(ctx) {
// 		var result T
// 		if err := cursor.Decode(&result); err != nil {
// 			return nil, err
// 		}
// 		results = append(results, result)
// 	}
// 	return results, nil
// }

package helpers

import (
	"context"
)

// Cursor is an interface that matches the methods used by DecodeCursor.
type Cursor interface {
	Next(context.Context) bool
	Decode(interface{}) error
}

// DecodeCursor decodes the documents from the cursor into the provided slice.
func DecodeCursor[T any](ctx context.Context, cursor Cursor) ([]T, error) {
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
