package errmarshalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/longkai/rfc7807"
	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/caisson-go/levels/level3/errorbag"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TODO move context away, write span somewhere else
func ToRFC7807(ctx context.Context, rspErr error) any {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(rspErr)
	}

	code := errors.Code(rspErr)

	var rsp rfc7807.ProblemDetail
	if code == nil {
		rsp.Status = http.StatusInternalServerError
		rsp.Detail = rspErr.Error()
	} else {

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

	return rsp
}
