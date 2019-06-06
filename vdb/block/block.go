package block

import (
	"encoding/json"
	AVdbComm "github.com/ayachain/go-aya/vdb/common"
	"github.com/ipfs/go-cid"
)

type Block struct {

	AVdbComm.RawSigner

	/// block index
	Index uint64 `json:"index"`

	/// chain id
	ChainID string `json:"chainid"`

	/// prev block hash is a ipfs block CID
	Parent  string `json:"parent"`

	/// full block data CID
	ExtraData string `json:"extradata"`

	/// broadcasting time super master node package this block times.
	Timestamp uint64 `json:"timestamp"`

	/// append data
	AppendData []byte `json:"append"`

	/// block sub transactions, is a ipfs block cid
	Txc uint16	`json:"txc"`
	Txs string	`json:"txs"`

}

/// only in create a new chain then use
type GenBlock struct {
	Block
	Consensus	string	`json:"consensus"`
	Award map[string]uint64 `json:"award"`
}

//var (
//	Genesis = &Block{Index: -4}
//	Curr 	= &Block{Index: -3}
//	Latest 	= &Block{Index: -2}
//	Pending = &Block{Index: -1}
//)

func (b *Block) GetExtraDataCid() cid.Cid {

	c, err := cid.Decode( b.ExtraData )
	if err != nil {
		return cid.Undef
	}

	return c
}


func (b *Block) Encode() []byte {

	bs, err := json.Marshal(b)

	if err != nil {
		return nil
	}

	return bs
}

func (b *Block) Decode(bs []byte) error {

	//if bs[0] != 'b' {
	//	return errors.New("this raw bytes not a block.")
	//}
	return json.Unmarshal(bs, b)
}

func (b *GenBlock) Encode() []byte {

	bs, err := json.Marshal(b)

	if err != nil {
		return nil
	}

	return bs
}

func (b *GenBlock) Decode(bs []byte) error {
	return json.Unmarshal(bs, b)
}