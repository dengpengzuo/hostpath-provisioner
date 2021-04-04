package hostfscsi

import (
	"k8s.io/utils/mount"
)

// Mount mounts the source to target path.
func UtilMount(source, target, fstype string, options []string) error {
	dummyMount := mount.New("")
	return dummyMount.Mount(source, target, fstype, options)
}

// UnMount mounts the source to target path.
func UtilUnMount(target string) error {
	dummyMount := mount.New("")
	return dummyMount.Unmount(target)
}
