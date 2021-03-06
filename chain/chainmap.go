package chain

import (
	"context"
	"encoding/json"
	"errors"
	AMinerPool "github.com/ayachain/go-aya/chain/minerpool"
	AMsgCenter "github.com/ayachain/go-aya/chain/msgcenter"
	AStatDaemon "github.com/ayachain/go-aya/chain/sdaemon"
	"github.com/ayachain/go-aya/chain/txpool"
	"github.com/ayachain/go-aya/vdb"
	"github.com/ayachain/go-aya/vdb/im"
	"github.com/ayachain/go-aya/vdb/indexes"
	EAccount "github.com/ethereum/go-ethereum/accounts"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ipfs/core"
	"github.com/prometheus/common/log"
	"sync"
)

var chains = sync.Map{}

func Conn( ctx context.Context, chainId string, ind *core.IpfsNode, acc EAccount.Account ) error {

	maplist, err := GenBlocksMap(ind)
	if err != nil {
		return err
	}

	genBlock, exist := maplist[chainId]
	if !exist {
		return errors.New(`can't find the corresponding chain, did you not execute "add"`)
	}

	_, exist = chains.Load(genBlock.ChainID)
	if exist {
		return ErrAlreadyExistConnected
	}

	idxs := indexes.CreateServices(ind, genBlock.ChainID, false)
	if idxs == nil {
		return errors.New(`can't find the corresponding chain, did you not execute "add"`)
	}
	defer func() {
		if err := idxs.Close(); err != nil {
			log.Error(err)
		}
	}()

	lidx, err := idxs.GetLatest()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("Read Block: %012d CID: %v", lidx.BlockIndex, lidx.FullCID.String())

	vdbfs, err := vdb.LinkVFS(genBlock.ChainID, ind, idxs)
	if err != nil {
		return ErrCantLinkToChainExpected
	}

	/// core services
	asd := AStatDaemon.NewDaemon(AStatDaemon.DefaultConfig)

	amp := AMinerPool.NewPool( ind, chainId, idxs, asd )

	amc := AMsgCenter.NewCenter( ind, vdbfs, AMsgCenter.DefaultTrustedConfig, asd )

	txp := txpool.NewTxPool( ind, chainId, vdbfs, acc, AMsgCenter.GetChannelTopics(chainId, AMsgCenter.MessageChannelMiningBlock), asd)

	ac := &aChain{
		ChainId:chainId,
		INode:ind,
		CVFS:vdbfs,
		TXP:txp,
		IDX:idxs,
		AMC:amc,
		AMP:amp,
		ASD:asd,
		CancelCh:make(chan struct{}),
	}

	chains.Store(genBlock.ChainID, ac)
	defer chains.Delete(genBlock.ChainID)

	sctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return ac.LinkStart(sctx)
}

func GenBlocksMap( ind *core.IpfsNode ) (map[string]*im.GenBlock, error) {

	dsk := datastore.NewKey(AChainMapKey)
	val, err := ind.Repo.Datastore().Get(dsk)

	if err != nil {
		if err != datastore.ErrNotFound {
			return nil, err
		}
	}

	rmap := make(map[string]*im.GenBlock)
	if val != nil {

		if err := json.Unmarshal( val, &rmap ); err != nil {
			return rmap, nil
		}

	}

	return rmap, nil
}

func AddChain( genBlock *im.GenBlock, ind *core.IpfsNode, r bool ) error {

	maplist, err := GenBlocksMap(ind)
	if err != nil {
		return err
	}

	_, exist := maplist[genBlock.ChainID]
	if exist && !r {
		return errors.New("chain are already exist")
	}

	// Create indexes
	idxServer := indexes.CreateServices(ind, genBlock.ChainID, r)
	if idxServer == nil {
		return errors.New("create chain indexes services failed")
	}
	defer func() {
		if err := idxServer.Close(); err != nil {
			log.Error(err)
		}
	}()

	// Create CVFS and write genBlock
	if _, err := vdb.CreateVFS(genBlock, ind, idxServer); err != nil {
		return err
	}

	maplist[genBlock.ChainID] = genBlock
	bs, err := json.Marshal(maplist)
	if err != nil {
		return err
	}

	dsk := datastore.NewKey(AChainMapKey)
	if err := ind.Repo.Datastore().Put( dsk, bs ); err != nil {
		return err
	}

	return nil
}

func GetChainByIdentifier(chainId string) Chain {

	v, ok := chains.Load(chainId)
	if !ok {
		return nil
	}

	c, ok := v.(Chain)
	if !ok {
		return nil
	}

	return c
}

func DisconnectionAll() {

	chains.Range(func(key, value interface{}) bool {

		c, ok := value.(Chain)
		if !ok {
			return true
		}

		c.Disconnect()
		return true
	})
}