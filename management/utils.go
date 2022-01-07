package management

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func serializeInput(input interface{}) (io.Reader, error) {
	m, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("Unable to serialize content: %w", err)
	}
	data := bytes.NewReader(m)
	return data, nil
}
