// +build aix darwin linux nacl solaris

package godirwalk

import "syscall"

func inoFromDirent(de *syscall.Dirent) uint64 {
	// cast necessary on file systems that store ino as different type
	return uint64(de.Ino)
}
