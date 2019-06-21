package node

import (
	AvdbComm "github.com/ayachain/go-aya/vdb/common"
	AIndexes "github.com/ayachain/go-aya/vdb/indexes"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

type aCache struct {

	writer

	headAPI AIndexes.IndexesServices
	source *leveldb.DB
	cdb *leveldb.DB

	delKeys [][]byte
}

func newCache( sourceDB *leveldb.DB ) (Caches, error) {

	memsto := storage.NewMemStorage()

	mdb, err := leveldb.Open(memsto, nil)
	if err != nil {
		return nil, err
	}

	c := &aCache{
		source:sourceDB,
		cdb:mdb,
	}

	return c, nil
}

func (cache *aCache) GetNodeByPeerId( peerId string ) (*Node, error) {

	bs, err := AvdbComm.CacheGet( cache.source, cache.cdb, []byte(peerId) )
	if err != nil {
		return nil, err
	}

	nd := &Node{}

	if err := nd.Decode(bs); err != nil {
		return nil, err
	}

	return nd, nil
}


func (cache *aCache) Close() {
	_ = cache.cdb.Close()
}


func (cache *aCache) MergerBatch() *leveldb.Batch {

	batch := &leveldb.Batch{}

	it := cache.cdb.NewIterator(nil, nil)

	for _, delk := range cache.delKeys {
		batch.Delete(delk)
	}

	for it.Next() {

		batch.Put( it.Key(), it.Value() )

	}


	return batch
}


func (cache *aCache) Update( peerId string, node *Node ) error {

	exist, err := AvdbComm.CacheHas( cache.source, cache.cdb, []byte(peerId) )

	if err != nil {
		return err
	}

	if !exist {
		return leveldb.ErrNotFound
	}

	return cache.cdb.Put([]byte(peerId), node.Encode(), nil)

}


func (cache *aCache) Insert( peerId string, node *Node ) error {

	return cache.cdb.Put([]byte(peerId), node.Encode(), nil)

}


func (cache *aCache) Del( peerId string ) {

	_ = cache.cdb.Delete([]byte(peerId), nil)

	cache.delKeys = append(cache.delKeys, []byte(peerId))

}