package providers

import "context"

// Processor defines a set of functions which a processor needs.
// Todo: add errorGroups for error handling.
type ProcessorProvider interface {
	ProcessImages(ctx context.Context, in chan int)
}
