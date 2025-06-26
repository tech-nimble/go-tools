package render

import (
	"encoding/json"
	"net/http"

	"github.com/google/jsonapi"
)

type JSONAPI struct {
	Data  any
	Meta  *jsonapi.Meta  `json:"-"`
	Links *jsonapi.Links `json:"-"`
}

func (J JSONAPI) Render(w http.ResponseWriter) error {
	J.WriteContentType(w)

	payload, err := jsonapi.Marshal(J.Data)
	if err != nil {
		return err
	}

	if p, ok := payload.(*jsonapi.ManyPayload); ok {
		if J.Links != nil {
			p.Links = J.Links
		}

		if J.Meta != nil {
			p.Meta = J.Meta
		}

		payload = p
	}

	return json.NewEncoder(w).Encode(payload)
}

func (J JSONAPI) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{jsonapi.MediaType}
	}
}
