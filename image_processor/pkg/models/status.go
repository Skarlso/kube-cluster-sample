package models

// Status is an Image status representation
type Status int

const (
	// PENDING -- not yet send to face recognition service
	PENDING Status = iota
	// PROCESSED -- processed by face recognition service; even if no person was found for the image
	PROCESSED
	// FAILEDPROCESSING -- for whatever reason the processing failed and this image is flagged for a retry
	FAILEDPROCESSING
	// PROCESSING -- used to "claim" an image by an instance and mark as being processed by face-recog
	PROCESSING
)
