package world

import (
	"github.com/itiky/goPixelWorld/world/types"
)

const (
	// tileWorkersNum defines the number of Tile processing workers.
	tileWorkersNum = 32
	// procTileJobChSize defines the procTileJobCh buffer size.
	procTileJobChSize = 50000
)

// initProcessing inits the processing engine.
func (m *Map) initProcessing() {
	// Common channels
	m.procTileJobCh = make(chan *types.Tile, procTileJobChSize)
	m.procRequestCh = make(chan struct{})
	m.procAckCh = make(chan struct{})

	// Start workers and init output queue for each
	for i := 0; i < tileWorkersNum; i++ {
		m.procActions = append(m.procActions, make([]types.Action, 0, procTileJobChSize))
		go m.tileWorker(i)
	}

	// Start the main process ahead worker
	go m.processingWorker()
	m.processingStart()
}

// processingStart sends a request to start the next processing round.
func (m *Map) processingStart() {
	m.procRequestCh <- struct{}{}
}

// processingDone waits until the processing is done and output is ready to be collected.
func (m *Map) processingDone() {
	<-m.procAckCh
}
