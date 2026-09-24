package collector

import "github.com/BugLrd/HoneyPot-DashBoard/internal/event"

// Event is an alias to the shared event.Event type.
// It keeps the collector package free of the actual struct definition
// to avoid import cycles.
//
// All code that previously referred to collector.Event will now
// refer to event.Event via this alias.
//
// This file does not introduce new imports beyond the shared package.
//
// If new fields are added to the Event struct, ensure this alias
// continues to point to the updated type.

type Event = event.Event
