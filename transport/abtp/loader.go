package abtp

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go abtp xdp/abtp_xdp.c

import (
	"errors"
	"fmt"
	"net"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

// XDPLinker defines the interface for interacting with the network and XDP kernel objects.
type XDPLinker interface {
	InterfaceByName(name string) (*net.Interface, error)
	LoadObjects(obj *abtpObjects, opts *ebpf.CollectionOptions) error
	AttachXDP(opts link.XDPOptions) (link.Link, error)
}

// RealXDPLinker implements XDPLinker wrapping the actual system calls.
type RealXDPLinker struct{}

func (r *RealXDPLinker) InterfaceByName(name string) (*net.Interface, error) {
	return net.InterfaceByName(name)
}
func (r *RealXDPLinker) LoadObjects(obj *abtpObjects, opts *ebpf.CollectionOptions) error {
	return loadAbtpObjects(obj, opts)
}
func (r *RealXDPLinker) AttachXDP(opts link.XDPOptions) (link.Link, error) {
	return link.AttachXDP(opts)
}

// Loader manages the XDP program lifecycle.
type Loader struct {
	ifaceName string
	link      link.Link
	objs      abtpObjects
	linker    XDPLinker
}

// NewLoader creates a new XDP loader for the given interface.
func NewLoader(ifaceName string, linker XDPLinker) *Loader {
	if linker == nil {
		linker = &RealXDPLinker{}
	}
	return &Loader{
		ifaceName: ifaceName,
		linker:    linker,
	}
}

// Attach loads and attaches the XDP program to the network interface.
func (l *Loader) Attach() error {
	// Look up the network interface by name.
	iface, err := l.linker.InterfaceByName(l.ifaceName)
	if err != nil {
		return fmt.Errorf("lookup network iface %q: %w", l.ifaceName, err)
	}

	// Load pre-compiled programs into the kernel.
	if err := l.linker.LoadObjects(&l.objs, nil); err != nil {
		return fmt.Errorf("loading objects: %w", err)
	}

	// Attach the program.
	lnk, err := l.linker.AttachXDP(link.XDPOptions{
		Program:   l.objs.AbtpXdpProg,
		Interface: iface.Index,
	})
	if err != nil {
		return fmt.Errorf("could not attach XDP program: %w", err)
	}
	l.link = lnk

	return nil
}

// Detach removes the XDP program and cleans up resources.
func (l *Loader) Detach() error {
	var errs []error
	if l.link != nil {
		if err := l.link.Close(); err != nil {
			errs = append(errs, err)
		}
		l.link = nil
	}
	if err := l.objs.Close(); err != nil {
		errs = append(errs, err)
	}
	l.objs = abtpObjects{}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
