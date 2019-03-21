package utils

import (
	"context"
	"encoding/json"
	"github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs-files"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

//AFMS的ipfs files 相关子命令全部采用 -f=false 调用，若需要刷新到磁盘，需要手动调用flush，或者使用AFMS_FlushPath进行刷写磁盘
//禁用自动同步，因为每次都是计算一个块中到所有交易，每次修改都同步到磁盘，会产生大量Disk I/O 降低效率
//调用AFMS实现方法绝大多是情况都工作在内存中，是一些不确定都，或者尚未计算完全的交易

const AFMS_FILE  	=	"file"
const AFMS_DIR  	=	"directory"

type AFMS_Stat struct {
	Hash			string
	Size			uint64
	CumulativeSize 	uint64
	ChildBlocks		uint64
	Type 			string
}

func (as *AFMS_Stat) IsDir() bool {
	return strings.EqualFold(as.Type, AFMS_DIR)
}

func (as *AFMS_Stat) IsFile() bool {
	return strings.EqualFold(as.Type, AFMS_FILE)
}

func AFMS_PathStat(mpath string) (stat *AFMS_Stat, err error) {

	if !strings.HasPrefix(mpath,"/") {
		mpath = "/" + mpath
	}

	if r, err := shell.NewLocalShell().Request("files/stat").Arguments(mpath).Option("flush",false).Send(context.Background()); err != nil {

		return nil, err

	} else {

		if r.Error != nil {

			return nil, errors.New(r.Error.Error())

		} else {

			if bs, err := ioutil.ReadAll(r.Output); err != nil {

				return nil, err

			} else {

				stat = &AFMS_Stat{}

				if json.Unmarshal(bs, stat) != nil {

					return nil, err

				} else {

					return stat,nil

				}

			}

		}
	}

}

//检测文件是否存在
func AFMS_IsPathExist(mpath string) bool {

	if _, err := AFMS_PathStat(mpath); err != nil {
		return false
	} else {
		return true
	}

}

func AFMS_RemovePath(mpath string) bool {

	if !strings.HasPrefix(mpath,"/") {
		mpath = "/" + mpath
	}

	return shell.NewLocalShell().Request("files/rm").Arguments(mpath).Option("flush",false).Option("r", true).Exec(context.Background(), nil) == nil
}

func AFMS_DownloadPathToDir(source string, dist string) bool {
	return shell.NewLocalShell().Request("files/cp").Arguments(source, dist).Option("flush",false).Exec(context.Background(), nil) == nil
}

//装载Dapp运行文件
func AFMS_ReloadDapp(bdhash string, mfspath string) bool {

	if !strings.HasPrefix(bdhash,"/ipfs") {
		bdhash = "/ipfs/" + bdhash
	}

	if !strings.HasPrefix(mfspath, "/") {
		mfspath = "/" + mfspath
	}

	if AFMS_IsPathExist(mfspath) {
		//文件存在直接删除
		AFMS_RemovePath(mfspath)
	}

	return AFMS_DownloadPathToDir(bdhash, mfspath)

}

func AFMS_DestoryDapp(nsp string) bool {

	var mfsTpath string

	if nsp[0] == '/' {
		mfsTpath = nsp
	} else {
		mfsTpath = "/" + nsp
	}

	return AFMS_RemovePath(mfsTpath)
}

func AFMS_ReadDappCode(path string) (code string, err error) {

	mfsTpath := "/" + path + "/_dapp/main.lua"

	bs,err := shell.NewLocalShell().Request("files/read").Arguments(mfsTpath).Option("flush",false).Send(context.Background())

	if err != nil {
		return "", err
	} else {

		if bs.Error != nil {
			return "", errors.New( bs.Error.Error() )
		}

		if codebs,err := ioutil.ReadAll(bs.Output); err != nil {
			return "", err
		} else {
			return string(codebs), nil
		}
	}

}

func AFMS_RenameFile(bpath string, s string, d string) error {

	source := "/" + bpath + "/" + s
	dist := "/" + bpath + "/" + d

	return shell.NewLocalShell().Request("files/mv").Arguments(source,dist).Option("flush",false).Exec(context.Background(), nil)
}

func AFMS_CreateFile(path string, name string, data[] byte) error {

	fpath := "/" + path + "/" + name

	shell := shell.NewLocalShell()

	fr := files.NewBytesFile([]byte(data))
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", fr)})
	fileReader := files.NewMultiFileReader(slf, true)

	reqb := shell.Request("files/write").Arguments(fpath).Body(fileReader).Option("e",true).Option("t",true).Option("flush",false)

	if err := reqb.Exec(context.Background(), nil); err != nil {
		return err
	} else {
		return nil
	}
}

func AFMS_FileAppend(bpathstring string, fpath string, data[] byte) error {

	spath := "/" + bpathstring + "/" + fpath

	stat, err := AFMS_PathStat(spath)

	if err != nil {
		return err
	} else if stat.Type != AFMS_FILE {
		return errors.New("AFMS_FileAppend : " + spath + " not a file.")
	}

	fr := files.NewBytesFile(data)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", fr)})
	fileReader := files.NewMultiFileReader(slf, true)

	reqb := shell.NewLocalShell().Request("files/write").Arguments(spath).Body(fileReader).Option("o",stat.Size + 1).Option("flush",false)

	if err := reqb.Exec(context.Background(), nil); err != nil {
		return err
	} else {
		return nil
	}

}

func AFMS_FlushPath(mpath string) error {
	return shell.NewLocalShell().Request("files/flush").Arguments(mpath).Exec(context.Background(),nil)
}