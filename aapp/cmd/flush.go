package cmd

import (
	"fmt"
	"github.com/ayachain/go-aya/aapp"
	ARsponse "github.com/ayachain/go-aya/response"
	cmds "github.com/ipfs/go-ipfs-cmds"
)

var flushCmd = &cmds.Command{
	Helptext:cmds.HelpText{
		Tagline: "Flush AApp ALVM file system to ipfs repo",
	},
	Arguments: []cmds.Argument {
		cmds.StringArg("aappns", true, false, "Path to AApp."),
	},
	Run:func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) error {

		ap := aapp.Manager.AAppOf(req.Arguments[0])

		if ap == nil {

			return ARsponse.EmitErrorResponse(
				re,
				fmt.Errorf("%v is not a daemoned AAppServices", req.Arguments[0]),
				)

		} else {

			c, err := ap.FlushMFS()
			if err != nil {
				return ARsponse.EmitErrorResponse(re, err)
			} else {
				return ARsponse.EmitSuccessResponse(re, c.String())
			}

		}
	},
}