package indexes

import (
	"encoding/json"
	"github.com/ayachain/go-aya/vdb/common"
	EComm "github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/go-cid"
)

type Index struct {
	common.RawDBCoder		`json:"-"`

	BlockIndex 	uint64 		`json:"Index, omitempty"`
	Hash		EComm.Hash 	`json:"Hash, omitempty"`
	FullCID 	cid.Cid 	`json:"FullCID, omitempty"`
}


func(i *Index) Encode() []byte {

	bs, err := json.Marshal(i)

	if err != nil{
		return nil
	}

	return bs
}

func(i *Index) Decode(bs []byte) error {
	return json.Unmarshal(bs, i)
}