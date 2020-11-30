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
	BlockSize() uint64
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
	blockSize  uint64
	blockTotal uint64 // capacity
	blockFree  uint64 // available
}

type hostpathDevice struct {
	mountDir string
	info     *fsInfo
}

func NewhostpathDevice(mountDir string) DeviceInf {
	return &hostpathDevice{
		mountDir: mountDir,
		info: &fsInfo{
			uuid:       "-",
			blockSize:  4096,
			blockTotal: 0,
			blockFree:  0,
		},
	}
}

func (f *hostpathDevice) Init() error {
	var e error = nil
	f.info.blockSize, f.info.blockTotal, f.info.blockFree, _, _, _, e = statfs(f.mountDir)
	if e != nil {
		return e
	}

	idFile := filepath.Join(f.mountDir, hostpath_id)
	f.info.uuid, e = readFsid(idFile)
	if os.IsNotExist(e) {
		f.info.uuid = uuid.New().String()
		e = writeFsId(idFile, f.info.uuid)
	}
	return e
}

func (f *hostpathDevice) DeviceUuid() string {
	return f.info.uuid
}

func (f *hostpathDevice) BlockSize() uint64 {
	return f.info.blockSize
}

func (f *hostpathDevice) TotalBlocks() uint64 {
	return f.info.blockTotal
}

func (f *hostpathDevice) FreeBlocks() uint64 {
	return f.info.blockFree
}

func (f *hostpathDevice) Alloc(id string, args ...interface{}) error {
	dir := filepath.Join(f.mountDir, id)
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

func (f *hostpathDevice) Remove(id string, args ...interface{}) error {
	dir := filepath.Join(f.mountDir, id)
	recDirName := filepath.Join(f.mountDir, "."+uuid.New().String())
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
func statfs(path string) (uint64, uint64, uint64, uint64, uint64, uint64, error) {
	statfs := &unix.Statfs_t{}
	err := unix.Statfs(path, statfs)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	//
	// f_bfree  => free blocks in fs
	// f_bavail => free blocks available to unprivileged user
	//
	return uint64(statfs.Bsize), statfs.Blocks, statfs.Bavail, statfs.Bfree, statfs.Files, statfs.Ffree, nil
}
