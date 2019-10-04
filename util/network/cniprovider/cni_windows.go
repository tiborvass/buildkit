package cniprovider

import (
	"context"
	"os"
	"path/filepath"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

func (ns *cniNS) Set(s *specs.Spec) {
	/*
		oci.WithLinuxNamespace(specs.LinuxNamespace{
			Type: specs.NetworkNamespace,
			Path: ns.path,
		})(nil, nil, nil, s)
	*/
}

func (ns *cniNS) Close() error {
	err := ns.handle.Remove(context.TODO(), ns.id, ns.path)

	/*
		if err1 := unix.Unmount(ns.path, unix.MNT_DETACH); err1 != nil {
			if err1 != syscall.EINVAL && err1 != syscall.ENOENT && err == nil {
				err = errors.Wrap(err1, "error unmounting network namespace")
			}
		}
	*/
	if err1 := os.RemoveAll(filepath.Dir(ns.path)); err1 != nil && !os.IsNotExist(err1) && err == nil {
		err = errors.Wrap(err, "error removing network namespace")
	}

	return err
}
