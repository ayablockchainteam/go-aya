package assets

import (
	"bytes"
	"encoding/binary"
	"github.com/ayachain/go-aya/vdb/common"
	EComm "github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/go-mfs"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
	"sync"
)

type aAssetes struct {
	AssetsAPI
	*mfs.Directory
	RWLocker sync.RWMutex
	rawdb *leveldb.DB
	mfsstorage storage.Storage
}

func CreateServices( mdir *mfs.Directory, rdOnly bool ) AssetsAPI {

	api := &aAssetes{
		Directory:mdir,
	}

	api.rawdb, api.mfsstorage = common.OpenExistedDB( mdir, DBPATH, rdOnly )

	return api
}

func (api *aAssetes) DBKey() string {
	return DBPATH
}

func (api *aAssetes) VotingCountOf( key []byte ) ( uint64, error ) {

	ast, err := api.AssetsOf(key)
	if err != nil {
		return 0, err
	}

	return ast.Vote, nil

}

func (api *aAssetes) AssetsOf( key []byte ) ( *Assets, error ) {

	bnc, err := api.rawdb.Get(key, nil)
	if err != nil {
		return nil, err
	}

	rcd := &Assets{}
	if err := rcd.Decode(bnc); err != nil {
		return nil, err
	}

	return rcd, nil
}


func (api *aAssetes) OpenVDBTransaction() (*leveldb.Transaction, *sync.RWMutex, error) {

	tx, err := api.rawdb.OpenTransaction()
	if err != nil {
		return nil, nil, err
	}

	return tx, &api.RWLocker, nil
}


func (api *aAssetes) Close() {

	api.RWLocker.Lock()
	defer api.RWLocker.Unlock()

	_ = api.rawdb.Close()
	_ = api.mfsstorage.Close()
	_ = api.Flush()
}

func (api *aAssetes) GetLockedTop100() ( []*SortAssets, error ) {

	list := make([]*SortAssets, 100)

	sbs := bytes.NewBuffer([]byte(DBTopIndexPrefix))
	sbs.Write( []byte{ 0x00, 0x00 } )

	ebs := bytes.NewBuffer([]byte(DBTopIndexPrefix))
	ebs.Write( []byte{ 0xff, 0xff })

	topIt := api.rawdb.NewIterator( &util.Range{Start:sbs.Bytes(), Limit:ebs.Bytes()}, nil )
	defer topIt.Release()

	for topIt.Next() {

		rcd, err := api.AssetsOf( topIt.Value() )
		if err != nil {
			return nil, err
		}

		assets := &SortAssets{
			Address : EComm.BytesToAddress(topIt.Value()),
			Assets : rcd,
		}

		nbs := topIt.Key()[ len([]byte(DBTopIndexPrefix)) : ]

		index := int(binary.BigEndian.Uint16(nbs))

		if index < 0 || index >= 100 {
			return nil, errors.New("array index bound")
		}

		list[index] = assets
	}

	return list, nil
}