package abft

import (
	"github.com/zilionixx/zilion-base/abft/dagidx"
	"github.com/zilionixx/zilion-base/abft/election"
	"github.com/zilionixx/zilion-base/hash"
	"github.com/zilionixx/zilion-base/inter/idx"
	"github.com/zilionixx/zilion-base/inter/pos"
)

type OrdererCallbacks struct {
	ApplyAtropos func(decidedFrame idx.Frame, atropos hash.Event) (sealEpoch *pos.Validators)

	EpochDBLoaded func(idx.Epoch)
}

type OrdererDagIndex interface {
	dagidx.ForklessCause
}

// Unlike processes events to reach finality on their order.
// Unlike abft.ZilionBFT, this raw level of abstraction doesn't track cheaters detection
type Orderer struct {
	config Config
	crit   func(error)
	store  *Store
	input  EventSource

	election *election.Election
	dagIndex OrdererDagIndex

	callback OrdererCallbacks
}

// New creates Orderer instance.
// Unlike ZilionBFT, Orderer doesn't updates DAG indexes for events, and doesn't detect cheaters
// It has only one purpose - reaching consensus on events order.
func NewOrderer(store *Store, input EventSource, dagIndex OrdererDagIndex, crit func(error), config Config) *Orderer {
	p := &Orderer{
		config:   config,
		store:    store,
		input:    input,
		crit:     crit,
		dagIndex: dagIndex,
	}

	return p
}
