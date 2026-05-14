// pbx-cluster-manager is a tool for managing FreeSWITCH clusters.
// It uses NATS for event-driven communication between components.
// The pbxevent package defines the structure of events published by pcm.ConnPool for every FreeSWITCH ESL event, and provides functions for parsing these events from JSON and retrieving header values.

// The Event struct represents the payload of a NATS message, containing metadata fields such as PBXCoreID, EventName, CallID, TimestampMs, and a map of Headers. The Parse function decodes a JSON byte slice into an Event struct, while the GetHeader method allows for case-insensitive retrieval of header values using canonical key lookup. The package also includes tests to ensure correct parsing of events and handling of edge cases, such as missing headers and malformed JSON.
package pbxevent
