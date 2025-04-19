package errmarshalhttp

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/longkai/rfc7807"
	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/caisson-go/pkg/errorbag"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Marshal(rspErr error, w http.ResponseWriter, r *http.Request) {

	span := trace.SpanFromContext(r.Context())
	if span != nil {
		span.RecordError(rspErr)
	}

	code := errors.Code(rspErr)

	var rsp rfc7807.ProblemDetail
	if code == nil {
		w.WriteHeader(http.StatusInternalServerError)
		rsp.Status = http.StatusInternalServerError
		rsp.Detail = rspErr.Error()
	} else {
		w.WriteHeader(code.HTTPCode())

		rsp.Status = code.HTTPCode()
		rsp.Type = code.Type()
		rsp.Title = code.Message()
		rsp.Detail = rspErr.Error()

	}

	pairs := errorbag.ListPairs(rspErr)
	rsp.Extensions = pairs

	if span != nil {
		for k, v := range pairs {
			// TODO: convert to attribute.Value
			span.SetAttributes(attribute.String(k, fmt.Sprintf("%v", v)))
		}
	}

	buf, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.ErrorContext(r.Context(), "failed to marshal error", slog.String("err", err.Error()), slog.String("stack", fmt.Sprintf("%+v", err)))
		return
	}


	w.Write(buf)
}
