package server

import (
	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/open-telemetry/opamp-go/server/types"
)

// AgentCore captures the minimal behavior an agent must expose to be managed
// by an AgentStore. Both implementations provide these methods on their Agent types.
type AgentCore interface {
	// OfferConnectionSettings offers connection settings (e.g., certificate rotation)
	// to the agent.
	OfferConnectionSettings(offers *protobufs.ConnectionSettingsOffers)
}

// CustomConfig exposes the ability to set a custom config on an agent.
// Implementations that support custom config can opt-in to this capability
// in addition to AgentCore.
type CustomConfig interface {
	// SetCustomConfig applies a custom config to the agent and optionally notifies
	// the caller on the next status update.
	SetCustomConfig(cfg *protobufs.AgentConfigMap, notifyNextStatusUpdate chan<- struct{})
}

// ReadonlyAgent represents a safe, read-only snapshot of an agent's state.
// Keep it intentionally abstract so each implementation can return its own
// cloned snapshot without leaking internal types.
type ReadonlyAgent interface {
	// Optionally add read-only accessors if a common surface is needed later.
}

// AgentStore is the common interface for managing agents and their connections.
//
// ID must be comparable (so it can be a map key); this supports uuid.UUID and
// any named alias such as data.InstanceId.
// A is the concrete agent type that satisfies AgentCore in the implementation.
type AgentStore[ID comparable, A AgentCore] interface {
	// RemoveConnection removes the connection and all agents associated with it.
	RemoveConnection(conn types.Connection)

	// FindAgent returns the agent by ID, or nil if not found.
	FindAgent(id ID) A

	// FindOrCreateAgent returns the existing agent or creates a new one associated with conn.
	FindOrCreateAgent(id ID, conn types.Connection) A

	// GetAgentReadonlyClone returns a read-only clone of the specified agent, or nil if not found.
	GetAgentReadonlyClone(id ID) ReadonlyAgent

	// GetAllAgentsReadonlyClone returns read-only clones of all agents keyed by ID.
	GetAllAgentsReadonlyClone() map[ID]ReadonlyAgent

	// OfferAgentConnectionSettings offers connection settings to the specified agent.
	OfferAgentConnectionSettings(id ID, offers *protobufs.ConnectionSettingsOffers)

	// SetCustomConfigForAgent applies a custom config to a specific agent.
	SetCustomConfigForAgent(id ID, cfg *protobufs.AgentConfigMap, notifyNextStatusUpdate chan<- struct{})
}

// OfferHashing is an optional capability interface for stores that compute and/or
// ensure a stable hash on connection settings offers before sending them. The
// opamp-go example store can implement this; other stores may ignore it.
type OfferHashing interface {
	EnsureOffersHash(offers *protobufs.ConnectionSettingsOffers)
}
