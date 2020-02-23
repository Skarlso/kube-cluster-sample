//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package providers

// SendProvider defines capabilities to send messages to NSQ
type SendProvider interface {
	SendImage(i uint64) error
}
