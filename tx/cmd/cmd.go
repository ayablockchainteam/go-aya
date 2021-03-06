package cmd

import cmds "github.com/ipfs/go-ipfs-cmds"

var TxCMDS = &cmds.Command{

	Helptext:cmds.HelpText{
		Tagline: "AyaChain tx commands.",
	},
	Subcommands: map[string]*cmds.Command{
		"get"	 		: getCMD,
		"receipt"		: receiptCMD,
		"pool"			: poolCMD,
		"list"			: listCMD,
		"count"			: countCMD,
		"publish"		: publishCMD,
	},

}