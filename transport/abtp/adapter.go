package abtp

import (
	"context"
	"net"
	"time"

	"github.com/Orolar-CNR/IntentCore/contracts"
)

// DefaultPort is the standard UDP port for ABTP listening.
const DefaultPort = 10000

// Adapter implements contracts.Transport using ABTP over UDP sockets.
// It bridges the gap between raw UDP packet listening and the IntentCore runtime pipeline.
type Adapter struct {
	port    int
	conn    *net.UDPConn
	handler contracts.EnvelopeHandler
	quit    chan struct{}
}

// NewAdapter creates a new ABTP transport adapter listening on the specified port.
func NewAdapter(port int) *Adapter {
	if port <= 0 {
		port = DefaultPort
	}
	return &Adapter{
		port: port,
		quit: make(chan struct{}),
	}
}

// Start begins listening for UDP packets and feeding them to the handler.
func (a *Adapter) Start(ctx context.Context, handler contracts.EnvelopeHandler) error {
	a.handler = handler

	addr := &net.UDPAddr{Port: a.port}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	a.conn = conn

	// Start background listener
	go a.listen(ctx)

	return nil
}

// listen reads packets from the UDP socket and dispatches them to the handler.
func (a *Adapter) listen(ctx context.Context) {
	// 64KB max UDP packet size
	buf := make([]byte, 65535)

	for {
		select {
		case <-a.quit:
			return
		case <-ctx.Done():
			return
		default:
			// Set a short read deadline so we can check quit/context periodically
			_ = a.conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

			n, _, err := a.conn.ReadFromUDP(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue // Normal timeout, loop again
				}
				// In a real system we would log this. For now, continue on non-fatal errors.
				// If the connection is closed, ReadFromUDP returns an error that breaks the loop.
				if errorsIsClosed(err) {
					return
				}
				continue
			}

			// Copy the payload before passing it to avoid data races with the next read
			payload := make([]byte, n)
			copy(payload, buf[:n])

			// Execute handler synchronously. Backpressure is handled by the pipeline/dispatcher.
			// If we wanted higher throughput, we could dispatch to a worker pool here, but for now
			// keeping it simple aligns with Phase 2 constraints.
			if err := a.handler(ctx, payload); err != nil {
				// We don't crash on handler errors (e.g. rejection), just log (if we had a logger) and continue
			}
		}
	}
}

// Stop halts the listener and cleans up resources.
func (a *Adapter) Stop(ctx context.Context) error {
	close(a.quit)
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

// errorsIsClosed tries to detect if a network error is due to a closed connection
func errorsIsClosed(err error) bool {
	// Crude check, normally you'd use errors.Is(err, net.ErrClosed) in Go 1.16+
	return err != nil && err.Error() == "use of closed network connection"
}
