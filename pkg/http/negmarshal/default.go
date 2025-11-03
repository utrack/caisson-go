package negmarshal

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
)

type responseObject struct {
	Data    any  `json:"data"`
	Error   any  `json:"error"`
	Success bool `json:"success"`
}

func MarshalerJSON() MarshalFunc {
	return func(ctx context.Context, w io.Writer, v any, errObj any) error {
		return json.NewEncoder(w).Encode(responseObject{
			Data:    v,
			Error:   errObj,
			Success: errObj == nil,
		})
	}
}

func MarshalerXML() MarshalFunc {
	return func(ctx context.Context, w io.Writer, v any, errObj any) error {
		return xml.NewEncoder(w).Encode(responseObject{
			Data:    v,
			Error:   errObj,
			Success: errObj == nil,
		})
	}
}
