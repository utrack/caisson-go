package negmarshal

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/longkai/rfc7807"
)

type responseObject struct {
	Data    any  `json:"data"`
	Error   any  `json:"error"`
	Success bool `json:"success"`
}

func MarshalerJSON() MarshalFunc {
	return func(ctx context.Context, w http.ResponseWriter, v any, errObj *rfc7807.ProblemDetail) error {
		w.Header().Set("Content-Type", "application/json")
		if errObj != nil && errObj.Status != 0 {
			w.WriteHeader(errObj.Status)
		}
		return json.NewEncoder(w).Encode(responseObject{
			Data:    v,
			Error:   errObj,
			Success: errObj == nil,
		})
	}
}

func MarshalerXML() MarshalFunc {
	return func(ctx context.Context, w http.ResponseWriter, v any, errObj *rfc7807.ProblemDetail) error {
		w.Header().Set("Content-Type", "application/xml")
		if errObj != nil && errObj.Status != 0 {
			w.WriteHeader(errObj.Status)
		}
		return xml.NewEncoder(w).Encode(responseObject{
			Data:    v,
			Error:   errObj,
			Success: errObj == nil,
		})
	}
}
