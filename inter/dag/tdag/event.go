package tdag

import (
	"github.com/zilionixx/zilion-base/hash"
	"github.com/zilionixx/zilion-base/inter/dag"
)

type TestEvent struct {
	dag.MutableBaseEvent
	Name string
}

func (e *TestEvent) AddParent(id hash.Event) {
	parents := e.Parents()
	parents.Add(id)
	e.SetParents(parents)
}
