package app

import "context"

// Storager is what the app needs from durable storage. The concrete
// *storage.Storage (constructed in main) satisfies it. Add the query methods
// your domain needs here, at the consumer, as the service grows.
type Storager interface {
	Ping(ctx context.Context) error
}
