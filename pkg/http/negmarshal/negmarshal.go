package negmarshal

import (
	"context"
	"net/http"

	"github.com/longkai/rfc7807"
	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/caisson-go/pkg/http/errmarshalhttp"
	contentnegotiation "gitlab.com/jamietanna/content-negotiation-go"
)

// MarshalFunc marshals the value in some single format (like json.Marshal or xml.Marshal).
type MarshalFunc func(ctx context.Context, w http.ResponseWriter, rsp any, errObj *rfc7807.ProblemDetail) error

// NegotiatedMarshalFunc marshals the value in the negotiated format,
// based on the request's Accept header.
type NegotiatedMarshalFunc func(r *http.Request, w http.ResponseWriter, rsp any, errObj error) error

// Default returns a NegotiatedMarshalFunc that supports JSON and XML outputs.
func Default() NegotiatedMarshalFunc {
	return New(map[string]MarshalFunc{
		"application/json": MarshalerJSON(),
		"application/xml":  MarshalerXML(),
	}, MarshalerJSON()).Marshal
}

type negotiator struct {
	mm       map[string]MarshalFunc
	defaultm MarshalFunc
	known    []string

	neg contentnegotiation.Negotiator
}

func New(mm map[string]MarshalFunc, defaultMarshaler MarshalFunc) *negotiator {
	keys := make([]string, 0, len(mm))
	for k := range mm {
		keys = append(keys, k)
	}
	return &negotiator{
		mm:       mm,
		defaultm: defaultMarshaler,
		known:    keys,
		neg:      contentnegotiation.NewNegotiator(keys...),
	}
}

func (n *negotiator) Marshal(r *http.Request, w http.ResponseWriter, v any, errObj error) error {
	var errRFC *rfc7807.ProblemDetail

	if errObj != nil {
		var ok bool
		errRFC, ok = errObj.(*rfc7807.ProblemDetail)
		if !ok {
			errRFC = errmarshalhttp.ToRFC7807(r.Context(), errObj)
		}
	}
	accepts := r.Header.Get("Accept")
	if accepts == "" || accepts == "*/*" {
		m := n.defaultm
		if m == nil {
			return errors.New("no default marshaler provided")
		}
		return m(r.Context(), w, v, errRFC)
	}
	mType, _, err := n.neg.Negotiate(accepts)
	if err != nil {
		return errors.Wrapd(err, "failed to negotiate content type", "accepts", accepts, "supported", n.known)
	}
	m, ok := n.mm[mType.String()]
	if !ok {
		panic("something's really wrong, no marshaler for known-negotiated type " + mType.String())
	}
	return m(r.Context(), w, v, errRFC)
}
