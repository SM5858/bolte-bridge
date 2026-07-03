// Package core defines and implements the core relay component of the bridge.
// The core relay communicates with connected messaging adapters via the
// core.Adapter interface.
package core

import (
	"context"

	"bolte-bridge/internal/relay"
)

// Adapter is the medium-specific edge of the bridge. The Email and Matrix
// adapters each implement it; the Core Relay drives them without knowing which
// medium is behind the interface.
type Adapter interface {
	// Medium reports which medium this adapter serves.
	Medium() relay.Medium

	// Fetch returns every message that has arrived since the last committed
	// cursor, oldest first. It reads from the medium (an IMAP fetch, a one-shot
	// Matrix sync) and advances the adapter's in-memory cursor, but does NOT
	// durably persist that cursor — that is Commit's job. Fetch is therefore
	// idempotent across process restarts until Commit succeeds: a crash after
	// Fetch but before Commit replays the same messages on the next run.
	Fetch(ctx context.Context) ([]relay.Message, error)

	// Send delivers one routed message into this adapter's medium. It is called
	// once per message; the core records per-message success or failure, so
	// Send reports the outcome of this message alone and does not batch.
	Send(ctx context.Context, msg relay.RoutedMessage) error

	// Commit durably advances this adapter's read cursor past everything
	// returned by the preceding Fetch. The core calls it only after the relay
	// of those messages has succeeded, as part of the tick's final transaction.
	// Until Commit returns nil, the messages from Fetch remain pending.
	Commit(ctx context.Context) error
}
