package generate

import (
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/common/horn"
	"github.com/Futuremine-chain/futuremine/service/pool"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"time"
)

const (
	module           = "generate"
	maxPackedTxCount = 1999
)

type Generate struct {
	horn        *horn.Horn
	pool        *pool.Pool
	dPos        dpos.IDPos
	chain       blockchain.IChain
	minerWorkCh chan bool
	stop        chan bool
	stopped     chan bool
}

func NewGenerate(chain blockchain.IChain, dPos dpos.IDPos, pool *pool.Pool, horn *horn.Horn) *Generate {
	return &Generate{
		pool:    pool,
		horn:    horn,
		chain:   chain,
		dPos:    dPos,
		stop:    make(chan bool),
		stopped: make(chan bool),
	}
}

func (g *Generate) Name() string {
	return module
}

func (g *Generate) Start() error {
	go g.generate()
	log.Info("Generate started successfully", "module", module)
	return nil
}

func (g *Generate) Stop() error {
	g.stop <- true
	return nil
}

func (g *Generate) Info() map[string]interface{} {
	return make(map[string]interface{}, 0)
}

func (g *Generate) generate() {
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case _, _ = <-g.stop:
			log.Info("Stop generate block")
			return
		case t := <-ticker:
			g.generateBlock(t)
		}
	}
}

func (g *Generate) generateBlock(now time.Time) {
	header, err := g.chain.NextHeader(uint64(now.Unix()))
	if err != nil {
		log.Error("Failed to generate next header", "module", module, "error", err)
		return
	}
	if err := g.dPos.CheckTime(header, g.chain); err != nil {
		return
	}

	err = g.dPos.CheckSigner(header, g.chain)
	if err != nil {
		//.Warn("check winner failed!", "height", header.Height, "error", err)
		return
	}
	txs := g.pool.NeedPackaged(maxPackedTxCount)
	nextBlock, err := g.chain.NextBlock(txs, uint64(now.Unix()))
	if err != nil {
		log.Error("Failed to generate block", "module", module, "error", err)
	}
	// Check if it is your turn to make blocks

	log.Info("Block generation successful", "module", module,
		"height", nextBlock.GetHeight(),
		"msgcount", nextBlock.BlockBody().Count(),
		"hash", nextBlock.GetHash().String(),
		"signer", nextBlock.GetSigner().String(),
	)
	g.horn.BroadcastBlock(nextBlock)
}
