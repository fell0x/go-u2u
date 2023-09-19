package gossip

import (
	"sync/atomic"

	"github.com/unicornultrafoundation/go-hashgraph/hash"
	"github.com/unicornultrafoundation/go-hashgraph/native/idx"
	"github.com/unicornultrafoundation/go-u2u/libs/core/state"
	"github.com/unicornultrafoundation/go-u2u/libs/core/types"

	"github.com/unicornultrafoundation/go-u2u/gossip/emitter"
	"github.com/unicornultrafoundation/go-u2u/native"
	"github.com/unicornultrafoundation/go-u2u/utils/wgmutex"
	"github.com/unicornultrafoundation/go-u2u/valkeystore"
	"github.com/unicornultrafoundation/go-u2u/vecmt"
)

type emitterWorldProc struct {
	s *Service
}

type emitterWorldRead struct {
	*Store
}

// emitterWorld implements emitter.World interface
type emitterWorld struct {
	emitterWorldProc
	emitterWorldRead
	*wgmutex.WgMutex
	emitter.TxPool
	valkeystore.SignerI
	types.Signer
}

func (ew *emitterWorldProc) Check(emitted *native.EventPayload, parents native.Events) error {
	// sanity check
	return ew.s.checkers.Validate(emitted, parents.Interfaces())
}

func (ew *emitterWorldProc) Process(emitted *native.EventPayload) error {
	done := ew.s.procLogger.EventConnectionStarted(emitted, true)
	defer done()
	return ew.s.processEvent(emitted)
}

func (ew *emitterWorldProc) Broadcast(emitted *native.EventPayload) {
	// PM listens and will broadcast it
	ew.s.feed.newEmittedEvent.Send(emitted)
}

func (ew *emitterWorldProc) Build(e *native.MutableEventPayload, onIndexed func()) error {
	return ew.s.buildEvent(e, onIndexed)
}

func (ew *emitterWorldProc) DagIndex() *vecmt.Index {
	return ew.s.dagIndexer
}

func (ew *emitterWorldProc) IsBusy() bool {
	return atomic.LoadUint32(&ew.s.eventBusyFlag) != 0 || atomic.LoadUint32(&ew.s.blockBusyFlag) != 0
}

func (ew *emitterWorldProc) StateDB() *state.StateDB {
	statedb, err := ew.s.store.evm.StateDB(ew.s.store.GetBlockState().FinalizedStateRoot)
	if err != nil {
		return nil
	}
	return statedb
}

func (ew *emitterWorldProc) IsSynced() bool {
	return ew.s.handler.syncStatus.AcceptEvents()
}

func (ew *emitterWorldProc) PeersNum() int {
	return ew.s.handler.peers.Len()
}

func (ew *emitterWorldRead) GetHeads(epoch idx.Epoch) hash.Events {
	return ew.Store.GetHeadsSlice(epoch)
}

func (ew *emitterWorldRead) GetLastEvent(epoch idx.Epoch, from idx.ValidatorID) *hash.Event {
	return ew.Store.GetLastEvent(epoch, from)
}

func (ew *emitterWorldRead) GetLowestBlockToDecide() idx.Block {
	return ew.Store.GetLlrState().LowestBlockToDecide
}

func (ew *emitterWorldRead) GetBlockRecordHash(n idx.Block) *hash.Hash {
	return ew.Store.GetBlockRecordHash(n)
}

func (ew *emitterWorldRead) GetBlockEpoch(block idx.Block) idx.Epoch {
	return ew.Store.FindBlockEpoch(block)
}

func (ew *emitterWorldRead) GetLowestEpochToDecide() idx.Epoch {
	return ew.Store.GetLlrState().LowestEpochToDecide
}

func (ew *emitterWorldRead) GetEpochRecordHash(epoch idx.Epoch) *hash.Hash {
	record := ew.Store.GetFullEpochRecord(epoch)
	if record == nil {
		return nil
	}
	h := record.Hash()
	return &h
}
