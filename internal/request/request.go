package request

import (
	"errors"
	"io"
	"strings"
)

var (
	ERR_BAD_REQUEST = errors.New("bad request")
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(err, ERR_BAD_REQUEST)
	}

	raw := string(b)
	rawSli := strings.Split(raw, "\r\n")
	if len(rawSli) < 1 {
		return nil, ERR_BAD_REQUEST
	}

	req := Request{}
	reqLine := rawSli[0]

	err = parseRequestLine(reqLine, &req)
	if err != nil {
		return nil, errors.Join(err, ERR_BAD_REQUEST)
	}

	return &req, nil
}
