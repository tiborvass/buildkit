package cniprovider

import (
	"context"
	"os"
	"path/filepath"

	"github.com/containerd/go-cni"
	"github.com/gofrs/flock"
	"github.com/moby/buildkit/identity"
	"github.com/moby/buildkit/util/network"
	"github.com/pkg/errors"
)

type Opt struct {
	Root       string
	ConfigPath string
	BinaryDir  string
}

func New(opt Opt) (network.Provider, error) {
	if _, err := os.Stat(opt.ConfigPath); err != nil {
		return nil, errors.Wrapf(err, "failed to read cni config %q", opt.ConfigPath)
	}
	if _, err := os.Stat(opt.BinaryDir); err != nil {
		return nil, errors.Wrapf(err, "failed to read cni binary dir %q", opt.BinaryDir)
	}

	cniHandle, err := cni.New(
		cni.WithMinNetworkCount(2),
		cni.WithConfFile(opt.ConfigPath),
		cni.WithPluginDir([]string{opt.BinaryDir}),
		cni.WithLoNetwork,
		cni.WithInterfacePrefix(("eth")))
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	cp := &cniProvider{CNI: cniHandle, root: opt.Root}
	if err := cp.initNetwork(); err != nil {
		return nil, err
	}
	return cp, nil
}

type cniProvider struct {
	cni.CNI
	root string
}

func (c *cniProvider) initNetwork() error {
	if v := os.Getenv("BUILDKIT_CNI_INIT_LOCK_PATH"); v != "" {
		l := flock.New(v)
		if err := l.Lock(); err != nil {
			return err
		}
		defer l.Unlock()
	}
	ns, err := c.New()
	if err != nil {
		return err
	}
	return ns.Close()
}

func (c *cniProvider) New() (network.Namespace, error) {
	id := identity.NewID()
	nsPath := filepath.Join(c.root, "net/cni", id)
	if err := os.MkdirAll(filepath.Dir(nsPath), 0700); err != nil {
		return nil, err
	}

	if err := createNetNS(nsPath); err != nil {
		os.RemoveAll(filepath.Dir(nsPath))
		return nil, err
	}

	if _, err := c.CNI.Setup(context.TODO(), id, nsPath); err != nil {
		os.RemoveAll(filepath.Dir(nsPath))
		return nil, errors.Wrap(err, "CNI setup error")
	}

	return &cniNS{path: nsPath, id: id, handle: c.CNI}, nil
}

type cniNS struct {
	handle cni.CNI
	id     string
	path   string
}
