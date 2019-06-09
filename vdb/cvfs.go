package vdb

import (
	"context"
	"errors"
	"fmt"
	AWrok "github.com/ayachain/go-aya/consensus/core/worker"
	AAssetses "github.com/ayachain/go-aya/vdb/assets"
	ABlock "github.com/ayachain/go-aya/vdb/block"
	AVdbComm "github.com/ayachain/go-aya/vdb/common"
	AHeader "github.com/ayachain/go-aya/vdb/headers"
	AReceipts "github.com/ayachain/go-aya/vdb/receipt"
	ATx "github.com/ayachain/go-aya/vdb/transaction"
	EComm "github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-mfs"
	"github.com/ipfs/go-unixfs"
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

var (
	aVFSDAGNodeConversionError = errors.New("conversion proto node expected")
)

type CVFS interface {
	Close()
	Blocks() 					ABlock.BlocksAPI		/// Body
	Headers() 					AHeader.HeadersAPI		/// Indexes
	Receipts() 					AReceipts.ReceiptsAPI	/// Receipt
	Assetses() 					AAssetses.AssetsAPI		/// Asset
	Transactions() 				ATx.TransactionAPI		/// Transaction
	OpenTransaction() 			(*Transaction, error)	/// Open Transaction to commit writing
	Flush(context.Context ) 	(cid.Cid, error)		/// Flush root cid
	SeekToBlock( bcid cid.Cid )  error
}

type aCVFS struct {

	CVFS
	*mfs.Root
	inode *core.IpfsNode
	ctx context.Context
	ctxCancel context.CancelFunc

	servies map[string]AVdbComm.VDBSerices
}


func CreateVFS( block *ABlock.GenBlock, ind *core.IpfsNode ) (cid.Cid, error) {

	var err error

	cvfs, err := LinkVFS( cid.Undef, ind)
	defer cvfs.Close()

	if err != nil {
		return cid.Undef, err
	}

	genBatchGroup := AWrok.NewGroup()

	// Award
	for addr, amount := range block.Award {

		assetBn := AAssetses.NewAssets( amount, amount, 0 ).Encode()

		genBatchGroup.Put( cvfs.Assetses().DBKey(), EComm.HexToAddress(addr).Bytes(), assetBn )

	}

	// Header
	genBlockHash := block.GetHash()
	err = cvfs.Headers().AppendHeaders(genBatchGroup, &AHeader.Header{BlockIndex:0, Hash:genBlockHash})
	if err != nil {
		return cid.Undef, err
	}


	// Block
	err = cvfs.Blocks().WriteGenBlock( genBatchGroup, block )
	if err != nil {
		return cid.Undef, err
	}

	txcommiter, err := cvfs.OpenTransaction()
	if err != nil {
		return cid.Undef, err
	}

	if err := txcommiter.Write( genBatchGroup ); err != nil {
		return cid.Undef, err
	}

	if err := txcommiter.Commit(); err != nil {
		return cid.Undef, err
	}

	return cvfs.Flush(context.TODO())
}


//ctx context.Context, aappns string, pnode *dag.ProtoNode, ind *core.IpfsNode
func LinkVFS( baseCid cid.Cid, ind *core.IpfsNode ) (CVFS, error) {

	ctx, cancel := context.WithCancel( context.Background() )
	root, err := newMFSRoot( ctx, baseCid, ind )

	if err != nil {
		return nil, err
	}

	assetDir, err := AVdbComm.LookupDBPath(root, AAssetses.DBPATH)
	if err != nil {
		return nil, err
	}
	headerDir, err := AVdbComm.LookupDBPath( root,  AHeader.DBPATH )
	if err != nil {
		return nil, err
	}


	blockDir, err := AVdbComm.LookupDBPath(root, ABlock.DBPath)
	if err != nil {
		return nil, err
	}

	txsDir, err := AVdbComm.LookupDBPath(root, ATx.DBPath)
	if err != nil {
		return nil, err
	}

	receiptDir, err := AVdbComm.LookupDBPath(root, AReceipts.DBPath)
	if err != nil {
		return nil, err
	}

	var (
			assetsServices 	= AAssetses.CreateServices( assetDir )
			headerServices	= AHeader.CreateServices( headerDir )
			blockServices	= ABlock.CreateServices( blockDir, headerServices )
			txServices		= ATx.CreateServices( txsDir )
			receiptServices	= AReceipts.CreateServices(receiptDir)
	)

	vfs := &aCVFS{
		ctx:ctx,
		ctxCancel:cancel,
		Root:root,
		inode:ind,
		servies: map[string]AVdbComm.VDBSerices{
			AAssetses.DBPATH 	: assetsServices,
			AHeader.DBPATH 		: headerServices,
			ABlock.DBPath		: blockServices,
			ATx.DBPath			: txServices,
			AReceipts.DBPath	: receiptServices,
		},
	}

	return vfs, nil
}

func newMFSRoot( ctx context.Context, c cid.Cid, ind *core.IpfsNode ) ( *mfs.Root, error ) {

	var pbnd *merkledag.ProtoNode

	if c == cid.Undef {

		pbnd = unixfs.EmptyDirNode()

	} else {

		rnd, err := ind.DAG.Get(ctx, c)
		if err != nil {
			return nil, err
		}

		var ok bool
		pbnd, ok = rnd.(*merkledag.ProtoNode)
		if !ok {
			return nil, aVFSDAGNodeConversionError
		}

	}

	mroot, err := mfs.NewRoot(ctx, ind.DAG, pbnd, func(i context.Context, i2 cid.Cid) error {
		fmt.Println("CVFS Published New CID : " + i2.String())
		return nil
	})

	if err != nil {
		return nil, err
	}

	return mroot, nil
}

func ( vfs *aCVFS ) SeekToBlock( bcid cid.Cid ) error {

	nd, err := mfs.FlushPath( context.TODO(), vfs.Root, "/" )
	if err != nil {
		return err
	}

	if nd.Cid().Equals(bcid) {
		return nil
	}

	vfs.ctxCancel()

	vfs.ctx, vfs.ctxCancel = context.WithCancel(context.Background())
	if err = vfs.Root.Close(); err != nil {
		return err
	}

	root, err := newMFSRoot(vfs.ctx, nd.Cid(), vfs.inode)
	if err != nil {
		return err
	}

	vfs.Root = root

	return nil
}

func ( vfs *aCVFS ) Close() {

	vfs.ctxCancel()

	for k, v := range vfs.servies {

		if v != nil {
			v.Close()
			fmt.Printf("%v services closed.\n", k)
		}

	}

}

func ( vfs *aCVFS ) Assetses() AAssetses.AssetsAPI {

	absapi, exist := vfs.servies[ AAssetses.DBPATH ]
	if !exist {
		return nil
	}

	api, ok := absapi.(AAssetses.AssetsAPI)
	if !ok {
		return nil
	}

	return api
}

func ( vfs *aCVFS ) Headers() AHeader.HeadersAPI {

	absapi, exist := vfs.servies[ AHeader.DBPATH ]
	if !exist {
		return nil
	}

	api, ok := absapi.(AHeader.HeadersAPI)
	if !ok {
		return nil
	}

	return api

}

func ( vfs *aCVFS ) Blocks() ABlock.BlocksAPI {

	absapi, exist := vfs.servies[ ABlock.DBPath ]
	if !exist {
		return nil
	}

	api, ok := absapi.(ABlock.BlocksAPI)
	if !ok {
		return nil
	}

	return api

}

func ( vfs *aCVFS ) Transactions() ATx.TransactionAPI {

	absapi, exist := vfs.servies[ ATx.DBPath ]
	if !exist {
		return nil
	}

	api, ok := absapi.(ATx.TransactionAPI)
	if !ok {
		return nil
	}

	return api

}

func ( vfs *aCVFS ) Receipts() AReceipts.ReceiptsAPI {

	absapi, exist := vfs.servies[ AReceipts.DBPath ]
	if !exist {
		return nil
	}

	api, ok := absapi.(AReceipts.ReceiptsAPI)
	if !ok {
		return nil
	}

	return api

}

func ( vfs *aCVFS ) OpenTransaction() (*Transaction, error) {

	var err error
	tx := &Transaction{
		transactions:make(map[string]*leveldb.Transaction),
		lockers: make(map[string]*sync.RWMutex),
	}

	for k, v := range vfs.servies {

		tx.transactions[k], tx.lockers[k], err = v.OpenVDBTransaction()
		if err != nil {
			return nil, err
		}

	}

	return tx, nil
}

func ( vfs *aCVFS ) WriteTaskGroup( group *AWrok.TaskBatchGroup) error {

	tx, err := vfs.OpenTransaction()
	if err != nil {
		return nil
	}

	if err := tx.Write(group); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func ( vfs *aCVFS ) Flush( ctx context.Context ) (cid.Cid, error) {

	nd, err := mfs.FlushPath( ctx, vfs.Root, "/" )
	if err != nil {
		return cid.Undef, err
	}

	return nd.Cid(), nil
}