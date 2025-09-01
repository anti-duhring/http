package request

import (
	"fmt"
	"io"
	"strings"

	"github.com/rotisserie/eris"
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, eris.Wrap(err, "error calling io.ReadAll")
	}

	raw := string(b)
	rawSli := strings.Split(raw, "\r\n")
	if len(rawSli) < 1 {
		return nil, fmt.Errorf("invalid request")
	}

	req := Request{}
	reqLine := rawSli[0]

	err = parseRequestLine(reqLine, &req)
	if err != nil {
		return nil, eris.Wrap(err, "error calling parseRequestLine")
	}

	return &req, nil
}
