package maple

import "encoding/json"

// WithSchemaVersion stamps the plugin's version into its schema. The committed schema
// has none, and `pulumi package add` refuses to generate an SDK from a schema without one.
func WithSchemaVersion(schema []byte, version string) ([]byte, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(schema, &fields); err != nil {
		return nil, err
	}
	v, err := json.Marshal(version)
	if err != nil {
		return nil, err
	}
	fields["version"] = v
	return json.Marshal(fields)
}
