package http

import (
	"net/http"

	"github.com/diegodesousas/go-devkit/pkg/encoding"
)

var serializer = encoding.NewJSONSerializer()

func WriteJson(w http.ResponseWriter, serializable any) error {
	body, err := serializer.Serialize(serializable)
	if err != nil {
		return err
	}

	if _, err = w.Write(body); err != nil {
		return err
	}

	return nil
}
