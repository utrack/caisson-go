package negmarshal

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
)

func MarshalerJSON() MarshalFunc {
	return func(ctx context.Context, w io.Writer, v any) error {
		return json.NewEncoder(w).Encode(v)
	}
}

func MarshalerXML() MarshalFunc {
	return func(ctx context.Context, w io.Writer, v any) error {
		return xml.NewEncoder(w).Encode(v)
	}
}
