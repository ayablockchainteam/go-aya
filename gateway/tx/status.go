package tx

import (
	RspFct "github.com/ayachain/go-aya/gateway/response"
	DSP "github.com/ayachain/go-aya/statepool"
	Atx "github.com/ayachain/go-aya/statepool/tx"
	"net/http"
)

//Method:Get
//URL:http:127.0.0.1:6001/tx/status?dappns=QmbLSVGXVJ4dMxBNkxneThhAnXVBWdGp7i2S42bseXh2hS&txhash=0xa44f6b24f913428a66623490a87dffaa6d7ed28b5d2e4932b453cdf34343ecc7
func TxStatusHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		RspFct.CreateError(RspFct.GATEWAY_ERROR_OnlySupport_Get).WriteToStream(&w);return
	}

	txhash := r.URL.Query().Get("txhash")

	if len(txhash) <= 0 {
		RspFct.CreateError(RspFct.GATEWAY_ERROR_MissParmas_TxHash).WriteToStream(&w);return
	}

	dappns := r.URL.Query().Get("dappns")

	if len(txhash) <= 0 {
		RspFct.CreateError(RspFct.GATEWAY_ERROR_MissParmas_DappNS).WriteToStream(&w);return
	}

	bindex,tx,s,rep := DSP.DappStatePool.GetTxStatus(dappns, txhash)

	switch s {
	case Atx.TxState_NotFound:
		RspFct.CreateError(RspFct.GATEWAY_ERROR_Tx_NotFound).WriteToStream(&w);return

	case Atx.TxState_Pending:

		r := &Atx.TxReceipt{
			BlockIndex:bindex,
			Tx:*tx,
			TxHash:tx.GetSha256Hash(),
			Status:"Pending",
			Response:nil,
		}

		RspFct.CreateSuccess(r).WriteToStream(&w)
		return

	case Atx.TxState_WaitPack:

		r := &Atx.TxReceipt{
			BlockIndex:bindex,
			Tx:*tx,
			TxHash:tx.GetSha256Hash(),
			Status:"WaitingPackage",
			Response:nil,
		}

		RspFct.CreateSuccess(r).WriteToStream(&w)
		return

	case Atx.TxState_Confirm:
		RspFct.CreateSuccess(rep).WriteToStream(&w)
		return
	}



}