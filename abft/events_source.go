package abft

import (
	"github.com/zilionixx/zilion-base/hash"
	"github.com/zilionixx/zilion-base/inter/dag"
)

// EventSource is a callback for getting events from an external storage.
type EventSource interface {
	HasEvent(hash.Event) bool
	GetEvent(hash.Event) dag.Event
}
