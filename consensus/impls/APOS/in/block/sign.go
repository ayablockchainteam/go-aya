package block

import (
	Avdb "github.com/ayachain/go-aya/vdb"
	"github.com/ipfs/go-ipfs/core"
)

func SignaturerStep ( msg interface{}, vdb Avdb.CacheCVFS, ind *core.IpfsNode ) (interface{}, error) {
	return msg ,nil
}