package pbxevent

import (
	"testing"
)

func TestParse_Complete(t *testing.T) {
	data := []byte(`{
		"_pbx_core_id": "core-1",
		"_event_name": "CHANNEL_ANSWER",
		"_call_id": "abc-123",
		"_timestamp_ms": 1715000000000,
		"headers": {
			"Event-Name": "CHANNEL_ANSWER",
			"Unique-Id": "abc-123",
			"Variable_domain_name": "namitech.io"
		}
	}`)

	e, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if e.PBXCoreID != "core-1" {
		t.Errorf("PBXCoreID = %q, want core-1", e.PBXCoreID)
	}
	if e.EventName != "CHANNEL_ANSWER" {
		t.Errorf("EventName = %q, want CHANNEL_ANSWER", e.EventName)
	}
	if e.CallID != "abc-123" {
		t.Errorf("CallID = %q, want abc-123", e.CallID)
	}
	if e.TimestampMs != 1715000000000 {
		t.Errorf("TimestampMs = %d, want 1715000000000", e.TimestampMs)
	}
	if len(e.Headers) != 3 {
		t.Errorf("len(Headers) = %d, want 3", len(e.Headers))
	}
}

func TestParse_MissingHeaders(t *testing.T) {
	data := []byte(`{"_pbx_core_id":"core-1","_event_name":"CHANNEL_CREATE","_call_id":"x","_timestamp_ms":0}`)
	e, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if e.Headers != nil {
		t.Errorf("Headers should be nil when absent from JSON")
	}
	// GetHeader must not panic on nil map.
	if got := e.GetHeader("Event-Name"); got != "" {
		t.Errorf("GetHeader on nil Headers = %q, want empty", got)
	}
}

func TestParse_MalformedJSON(t *testing.T) {
	_, err := Parse([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestGetHeader_Found(t *testing.T) {
	e := &Event{Headers: map[string]string{"Event-Name": "CHANNEL_ANSWER"}}
	if got := e.GetHeader("Event-Name"); got != "CHANNEL_ANSWER" {
		t.Errorf("GetHeader(Event-Name) = %q, want CHANNEL_ANSWER", got)
	}
}

func TestGetHeader_CaseVariant(t *testing.T) {
	e := &Event{Headers: map[string]string{"Event-Name": "CHANNEL_ANSWER"}}
	// "event-name" canonicalizes to "Event-Name"
	if got := e.GetHeader("event-name"); got != "CHANNEL_ANSWER" {
		t.Errorf("GetHeader(event-name) = %q, want CHANNEL_ANSWER", got)
	}
}

func TestGetHeader_Variable(t *testing.T) {
	// textproto: "variable_domain_name" → "Variable_domain_name"
	e := &Event{Headers: map[string]string{"Variable_domain_name": "namitech.io"}}
	if got := e.GetHeader("variable_domain_name"); got != "namitech.io" {
		t.Errorf("GetHeader(variable_domain_name) = %q, want namitech.io", got)
	}
}

func TestGetHeader_Missing(t *testing.T) {
	e := &Event{Headers: map[string]string{"Event-Name": "CHANNEL_CREATE"}}
	if got := e.GetHeader("Nonexistent-Header"); got != "" {
		t.Errorf("GetHeader(missing) = %q, want empty", got)
	}
}

func TestGetHeader_NilEvent(t *testing.T) {
	var e *Event
	if got := e.GetHeader("Event-Name"); got != "" {
		t.Errorf("GetHeader on nil Event = %q, want empty", got)
	}
}
