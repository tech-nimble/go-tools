// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

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

func (j JSONAPI) Render(w http.ResponseWriter) error {
	j.WriteContentType(w)

	payload, err := jsonapi.Marshal(j.Data)
	if err != nil {
		return err
	}

	if p, ok := payload.(*jsonapi.ManyPayload); ok {
		if j.Links != nil {
			p.Links = j.Links
		}

		if j.Meta != nil {
			p.Meta = j.Meta
		}

		payload = p
	}

	return json.NewEncoder(w).Encode(payload)
}

func (j JSONAPI) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{jsonapi.MediaType}
	}
}
