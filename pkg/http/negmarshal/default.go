package negmarshal

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type responseObject struct {
	Data    any  `json:"data"`
	Error   any  `json:"error"`
	Success bool `json:"success"`
}

func MarshalerJSON() MarshalFunc {
	return func(ctx context.Context, w http.ResponseWriter, v any, errObj any) error {
		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(responseObject{
			Data:    v,
			Error:   errObj,
			Success: errObj == nil,
		})
	}
}

func MarshalerXML() MarshalFunc {
	return func(ctx context.Context, w http.ResponseWriter, v any, errObj any) error {
		w.Header().Set("Content-Type", "application/xml")
		return xml.NewEncoder(w).Encode(responseObject{
			Data:    v,
			Error:   errObj,
			Success: errObj == nil,
		})
	}
}
