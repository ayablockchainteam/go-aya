package tx

import (
	"context"
	"fmt"
	ACStep "github.com/ayachain/go-aya/consensus/core/step"
	ADog "github.com/ayachain/go-aya/consensus/core/watchdog"
	AWork "github.com/ayachain/go-aya/consensus/core/worker"
	AMsgTx "github.com/ayachain/go-aya/vdb/transaction"
	"github.com/ipfs/go-cid"
	"time"
)

func (s *TxSender) Identifier( ) string {
	return identifier
}

func (s *TxSender) SetNextStep( ns ACStep.ConsensusStep ) {

}

func (s *TxSender) ChannelAccept() chan *ADog.MsgFromDogs {
	return acceptChan
}

func (s *TxSender) NextStep() ACStep.ConsensusStep {
	return nil
}

func (s *TxSender) Consensued( msg *ADog.MsgFromDogs ) {

	mcid, ok := msg.ExtraData.(cid.Cid)
	if !ok {

		err, ok := msg.ExtraData.(error)
		if ok {

			fmt.Printf("Consensued by warrning : %v", err)
			return

		} else {

			fmt.Println("Consensued by UnKown states" )
			return

		}

	}

	fmt.Printf("Consensued latest cid : %v", mcid.String())

	return
}

func (s *TxSender) StartListenAccept( ctx context.Context )() {

	go func() {

		fmt.Printf("%v Online\n", s.Identifier() )

		select {
		case dmsg := <-acceptChan:

			tx := &AMsgTx.Transaction{}

			tx.Decode(dmsg.Data)

			switch dmsg.Data[0] {
			case 'b':

				// check extradata type
				fbg, ok := dmsg.ExtraData.(*AWork.TaskBatchGroup)
				if !ok {
					dmsg.ResultDefer( ADog.FinalResult_ELV0 )
					break
				}

				// begin commit to main vdb, use transaction
				exeTx, err := s.cvfs.OpenTransaction()
				if err != nil {
					fmt.Print(err)
					dmsg.ResultDefer( ADog.FinalResult_ELV0 )
					break
				}

				// write batch to transaction
				if err := exeTx.Write(fbg); err != nil {
					fmt.Print(err)
					dmsg.ResultDefer( ADog.FinalResult_ELV0 )
					break
				}

				// try commit, if commit has any error, it will roll back automatically
				if err := exeTx.Commit(); err != nil {
					fmt.Print(err)
					dmsg.ResultDefer( ADog.FinalResult_ELV0 )
					break
				}

				// transaction success, try get latest vdb root path cid
				ctx, cancel := context.WithCancel(context.TODO())
				defer cancel()

				if mcid, err := s.cvfs.Flush(ctx); err != nil {
					dmsg.ExtraData = err
					dmsg.ResultDefer( ADog.FinalResult_ELV0 )
					break
				} else {
					dmsg.ExtraData = mcid
				}

				s.Consensued( dmsg )
				dmsg.ResultDefer( ADog.FinalResult_Success )
			}

		case <- ctx.Done() :
			break

		default:
			time.Sleep( time.Microsecond  * 100 )
		}

		fmt.Printf("%v consensus step close accept.", s.Identifier() )
	}()

}