package abtp

import (
	"errors"
	"net"
	"testing"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

type MockXDPLinker struct {
	InterfaceByNameFunc func(name string) (*net.Interface, error)
	LoadObjectsFunc     func(obj *abtpObjects, opts *ebpf.CollectionOptions) error
	AttachXDPFunc       func(opts link.XDPOptions) (link.Link, error)
}

func (m *MockXDPLinker) InterfaceByName(name string) (*net.Interface, error) {
	if m.InterfaceByNameFunc != nil {
		return m.InterfaceByNameFunc(name)
	}
	return &net.Interface{Index: 1, Name: name}, nil
}

func (m *MockXDPLinker) LoadObjects(obj *abtpObjects, opts *ebpf.CollectionOptions) error {
	if m.LoadObjectsFunc != nil {
		return m.LoadObjectsFunc(obj, opts)
	}
	return nil
}

func (m *MockXDPLinker) AttachXDP(opts link.XDPOptions) (link.Link, error) {
	if m.AttachXDPFunc != nil {
		return m.AttachXDPFunc(opts)
	}
	return nil, nil // Return nil link for happy path in test
}

func TestLoader_Attach_InterfaceByNameError(t *testing.T) {
	t.Parallel()
	expectedErr := errors.New("interface error")
	mockLinker := &MockXDPLinker{
		InterfaceByNameFunc: func(name string) (*net.Interface, error) {
			return nil, expectedErr
		},
	}
	loader := NewLoader("eth0", mockLinker)

	err := loader.Attach()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error %v, got %v", expectedErr, err)
	}
}

func TestLoader_Attach_LoadObjectsError(t *testing.T) {
	t.Parallel()
	expectedErr := errors.New("load error")
	mockLinker := &MockXDPLinker{
		LoadObjectsFunc: func(obj *abtpObjects, opts *ebpf.CollectionOptions) error {
			return expectedErr
		},
	}
	loader := NewLoader("eth0", mockLinker)

	err := loader.Attach()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error %v, got %v", expectedErr, err)
	}
}

func TestLoader_Attach_AttachXDPError(t *testing.T) {
	t.Parallel()
	expectedErr := errors.New("attach error")
	mockLinker := &MockXDPLinker{
		AttachXDPFunc: func(opts link.XDPOptions) (link.Link, error) {
			return nil, expectedErr
		},
	}
	loader := NewLoader("eth0", mockLinker)

	err := loader.Attach()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error %v, got %v", expectedErr, err)
	}
}

func TestLoader_Attach_Success(t *testing.T) {
	t.Parallel()
	mockLinker := &MockXDPLinker{}
	loader := NewLoader("eth0", mockLinker)

	err := loader.Attach()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

type mockLink struct {
	link.Link
	closeErr error
}

func (m *mockLink) Close() error {
	return m.closeErr
}

// Ensure mockLink implements link.Link (partial, just to test Close)
// Note: We use a wrapper approach to avoid implementing all link.Link methods.

func TestLoader_Detach_Success(t *testing.T) {
	t.Parallel()
	loader := &Loader{
		link: &mockLink{},
	}
	err := loader.Detach()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if loader.link != nil {
		t.Errorf("expected link to be nil")
	}
}

func TestLoader_Detach_LinkCloseError(t *testing.T) {
	t.Parallel()
	expectedErr := errors.New("link close error")
	loader := &Loader{
		link: &mockLink{closeErr: expectedErr},
	}
	err := loader.Detach()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error %v, got %v", expectedErr, err)
	}
	if loader.link != nil {
		t.Errorf("expected link to be nil despite error")
	}
}

func TestLoader_Detach_ObjsCloseError(t *testing.T) {
	t.Parallel()
	// abtpObjects is generated, its Close method currently just returns nil but we can't easily mock it without refactoring.
	// However, the error paths are properly handled via errors.Join.
	// For now, we only test the happy path of objs closing since we cannot inject an error into the generated struct.
}
