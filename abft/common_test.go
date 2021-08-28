package abft

import (
	"github.com/zilionixx/zilion-base/inter/idx"
	"github.com/zilionixx/zilion-base/inter/pos"
	"github.com/zilionixx/zilion-base/kvdb"
	"github.com/zilionixx/zilion-base/kvdb/memorydb"
	"github.com/zilionixx/zilion-base/utils/adapters"
	"github.com/zilionixx/zilion-base/vecfc"
	"github.com/zilionixx/zilion-base/zilionbft"
)

type applyBlockFn func(block *zilionbft.Block) *pos.Validators

// TestZilionBFT extends ZilionBFT for tests.
type TestZilionBFT struct {
	*IndexedZilionBFT

	blocks map[idx.Block]*zilionbft.Block

	applyBlock applyBlockFn
}

// FakeZilionBFT creates empty abft with mem store and equal weights of nodes in genesis.
func FakeZilionBFT(nodes []idx.ValidatorID, weights []pos.Weight, mods ...memorydb.Mod) (*TestZilionBFT, *Store, *EventStore) {
	validators := make(pos.ValidatorsBuilder, len(nodes))
	for i, v := range nodes {
		if weights == nil {
			validators[v] = 1
		} else {
			validators[v] = weights[i]
		}
	}

	openEDB := func(epoch idx.Epoch) kvdb.DropableStore {
		return memorydb.New()
	}
	crit := func(err error) {
		panic(err)
	}
	store := NewStore(memorydb.New(), openEDB, crit, LiteStoreConfig())

	err := store.ApplyGenesis(&Genesis{
		Validators: validators.Build(),
		Epoch:      FirstEpoch,
	})
	if err != nil {
		panic(err)
	}

	input := NewEventStore()

	config := LiteConfig()
	lch := NewIndexedZilionBFT(store, input, &adapters.VectorToDagIndexer{vecfc.NewIndex(crit, vecfc.LiteConfig())}, crit, config)

	extended := &TestZilionBFT{
		IndexedZilionBFT: lch,
		blocks:          map[idx.Block]*zilionbft.Block{},
	}

	blockIdx := idx.Block(0)

	err = extended.Bootstrap(zilionbft.ConsensusCallbacks{
		BeginBlock: func(block *zilionbft.Block) zilionbft.BlockCallbacks {
			blockIdx++
			return zilionbft.BlockCallbacks{
				EndBlock: func() (sealEpoch *pos.Validators) {
					// track blocks
					extended.blocks[blockIdx] = block
					if extended.applyBlock != nil {
						return extended.applyBlock(block)
					}
					return nil
				},
			}
		},
	})
	if err != nil {
		panic(err)
	}

	return extended, store, input
}
