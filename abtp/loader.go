package abtp

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go abtp xdp/abtp_xdp.c

import (
	"fmt"
	"net"

	"github.com/cilium/ebpf/link"
)

// Loader manages the XDP program lifecycle.
type Loader struct {
	ifaceName string
	link      link.Link
	objs      abtpObjects
}

// NewLoader creates a new XDP loader for the given interface.
func NewLoader(ifaceName string) *Loader {
	return &Loader{
		ifaceName: ifaceName,
	}
}

// Attach loads and attaches the XDP program to the network interface.
func (l *Loader) Attach() error {
	// Look up the network interface by name.
	iface, err := net.InterfaceByName(l.ifaceName)
	if err != nil {
		return fmt.Errorf("lookup network iface %q: %w", l.ifaceName, err)
	}

	// Load pre-compiled programs into the kernel.
	if err := loadAbtpObjects(&l.objs, nil); err != nil {
		return fmt.Errorf("loading objects: %w", err)
	}

	// Attach the program.
	lnk, err := link.AttachXDP(link.XDPOptions{
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
	}
	if err := l.objs.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during detach: %v", errs)
	}
	return nil
}
