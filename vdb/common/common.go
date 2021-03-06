package common

import (
	"encoding/binary"
	ADB "github.com/ayachain/go-aya-alvm-adb"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ipfs/go-mfs"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

//var OpenDBOpt = &opt.Options{
//	BlockCacher:opt.NoCacher,
//	BlockCacheCapacity:-1,
//	DisableBlockCache:true,
//	OpenFilesCacher:opt.NoCacher,
//	OpenFilesCacheCapacity:1,
//	WriteBuffer:0,
//}

type MessageTypePrefix uint8

var OpenDBOpt = &opt.Options{}

var StorageDBPaths = []string{"/nodes", "/blocks", "/assets", "/receipts", "/transactions"}

var WriteOpt = &opt.WriteOptions{
	NoWriteMerge:false,
	Sync:true,
}

type AMessageEncode interface {
	RawMessageEncode() []byte
	RawMessageDecode( bs []byte ) error
}

type RawSigner interface {
	RawSignEncode( account accounts.Account ) ([]byte, error)
	RawVerifyDecode( bs []byte ) error
}

func LittleEndianBytes (number uint64) []byte {
	enc := make([]byte, 8)
	binary.LittleEndian.PutUint64(enc, number)
	return enc
}

func BigEndianBytes (number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

func BigEndianBytesUint32 (n uint32 ) []byte {
	enc := make([]byte, 4)
	binary.BigEndian.PutUint32(enc, n)
	return enc
}

func BigEndianBytesUint16 ( n uint16 ) []byte {
	enc := make([]byte, 2)
	binary.BigEndian.PutUint16(enc, n)
	return enc
}


func OpenExistedDB( dir *mfs.Directory, dbkey string ) ( *leveldb.DB, *ADB.MFSStorage, error ) {

	dbstroage := ADB.NewMFSStorage(dir, dbkey)
	if dbstroage == nil {
		panic("create adb storage expected")
	}

	db, err := leveldb.Open(dbstroage, OpenDBOpt)

	if err != nil {
		return nil, nil, err
	}

	return db, dbstroage, nil
}



func CacheDel( originDB *leveldb.DB, cacheDB *leveldb.DB, key []byte ) {

	_ = originDB.Delete(key, WriteOpt)
	_ = cacheDB.Delete(key, WriteOpt)

}

func CacheHas( originDB *leveldb.DB, cacheDB *leveldb.DB, key []byte ) (bool, error) {

	exist, err := cacheDB.Has(key, nil)
	if err != nil {
		return false, err
	}

	if !exist {

		oexist, err := originDB.Has(key,nil)
		if err != nil {
			return false, err
		}

		return oexist, nil
	}

	return exist, nil
}

func CacheGet( originDB *leveldb.DB, cacheDB *leveldb.DB, key []byte ) ([]byte, error) {

	exist, err := cacheDB.Has(key, nil)
	if err != nil {
		return nil, err
	}

	if !exist {

		v, err := originDB.Get(key, nil)
		if err != nil {
			return nil, err
		}

		if err := cacheDB.Put(key, v, WriteOpt); err != nil {
			return nil, err
		}

		return v, nil

	} else {

		return cacheDB.Get(key, nil)

	}

}