package pbxevent

import (
	"encoding/json"
	"net/textproto"
)

// Event is the NATS message payload published by pcm.ConnPool for every
// FreeSWITCH ESL event. Metadata fields are prefixed with _ to distinguish
// them from FreeSWITCH headers.
type Event struct {
	PBXCoreID   string            `json:"_pbx_core_id"`
	EventName   string            `json:"_event_name"`
	CallID      string            `json:"_call_id"`
	TimestampMs int64             `json:"_timestamp_ms"`
	Headers     map[string]string `json:"headers"`
}

// Parse decodes a NATS message payload into an Event.
func Parse(data []byte) (*Event, error) {
	var e Event
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// GetHeader returns the value of a FreeSWITCH header using canonical key
// lookup. Case-insensitive for hyphen-separated keys (MIME canonical form):
// "event-name", "Event-Name", "EVENT-NAME" all resolve to the same entry.
// Underscore keys are only first-letter-capitalized:
// "variable_domain_name" → "Variable_domain_name".
func (e *Event) GetHeader(name string) string {
	if e == nil || e.Headers == nil {
		return ""
	}
	return e.Headers[textproto.CanonicalMIMEHeaderKey(name)]
}
