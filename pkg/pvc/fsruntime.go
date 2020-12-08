package pvc

import (
	"github.com/google/uuid"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DeviceInf interface {
	Init() error
	DeviceUuid() string
	BlockSize() int64
	TotalBlocks() uint64
	FreeBlocks() uint64
}

type DeviceAllocInf interface {
	Alloc(id string, args ...interface{}) error
	Remove(id string, args ...interface{}) error
}

const dir_model = 0755
const hostpath_id = ".ez-cloud-id"

type fsInfo struct {
	uuid       string
	blockSize  int64
	blockTotal uint64 // capacity
	blockFree  uint64 // available
}

type hostPathDevice struct {
	hostDir string
	info    *fsInfo
}

func NewHostPathDevice(hostDir string) DeviceInf {
	return &hostPathDevice{
		hostDir: hostDir,
		info: &fsInfo{
			uuid:       "-",
			blockSize:  4096,
			blockTotal: 0,
			blockFree:  0,
		},
	}
}

func (f *hostPathDevice) Init() error {
	var e error = nil
	f.info.blockSize, f.info.blockTotal, f.info.blockFree, _, _, _, e = statfs(f.hostDir)
	if e != nil {
		return e
	}

	idFile := filepath.Join(f.hostDir, hostpath_id)
	f.info.uuid, e = readFsid(idFile)
	if os.IsNotExist(e) {
		f.info.uuid = uuid.New().String()
		e = writeFsId(idFile, f.info.uuid)
	}
	return e
}

func (f *hostPathDevice) DeviceUuid() string {
	return f.info.uuid
}

func (f *hostPathDevice) BlockSize() int64 {
	return f.info.blockSize
}

func (f *hostPathDevice) TotalBlocks() uint64 {
	return f.info.blockTotal
}

func (f *hostPathDevice) FreeBlocks() uint64 {
	return f.info.blockFree
}

func (f *hostPathDevice) Alloc(id string, args ...interface{}) error {
	dir := filepath.Join(f.hostDir, id)
	var mkdirErr error
	if fi, err := os.Stat(dir); os.IsNotExist(err) {
		mkdirErr = os.Mkdir(dir, dir_model)
	} else if fi != nil && fi.IsDir() {
		mkdirErr = nil // reuse exists dir
	} else {
		mkdirErr = &os.PathError{Op: "mkdir", Path: dir, Err: err}
	}
	return mkdirErr
}

func (f *hostPathDevice) Remove(id string, args ...interface{}) error {
	dir := filepath.Join(f.hostDir, id)
	recDirName := filepath.Join(f.hostDir, newDirName())
	err := os.Rename(dir, recDirName)
	return err
}

func writeFsId(idFile string, id string) error {
	return ioutil.WriteFile(idFile, []byte(id), 0666)
}

func readFsid(idFile string) (string, error) {
	bytes, e := ioutil.ReadFile(idFile)
	if e != nil {
		return "", e
	}
	return string(bytes), nil
}

// return: bsize, capacity, available, free, inodes, inodesFree
func statfs(path string) (int64, uint64, uint64, uint64, uint64, uint64, error) {
	statfs := &unix.Statfs_t{}
	err := unix.Statfs(path, statfs)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	return statfs.Bsize, statfs.Blocks, statfs.Bavail, statfs.Bfree, statfs.Files, statfs.Ffree, nil
}

func newDirName() string {
	return "." + uuid.New().String()
}
