package request

import (
	"fmt"
	"io"
	"strings"

	"github.com/rotisserie/eris"
	"slices"
)

// ref: https://www.rfc-editor.org/rfc/rfc9110#section-9
var supportedMethods = []string{
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"PATCH",
	"DELETE",
	"CONNECT",
	"OPTIONS",
	"TRACE",
}

var supportedHttpVersions = []string{
	"1.1",
}

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(raw string, req *Request) error {
	parts := strings.Split(raw, " ")
	if len(parts) != 3 {
		return fmt.Errorf("expected 3 params but got %d", len(parts))
	}

	reqLine := RequestLine{}

	method := parts[0]
	if method == "" {
		return fmt.Errorf("invalid method")
	}

	isMethodSupported := slices.Contains(supportedMethods, method)
	if !isMethodSupported {
		return fmt.Errorf("invalid method")
	}

	reqLine.Method = method

	target := parts[1]
	if target == "" {
		return fmt.Errorf("invalid request target")
	}

	reqLine.RequestTarget = target

	httpVersion := parts[2]
	if httpVersion == "" {
		return fmt.Errorf("invalid http version")
	}

	httpVersionSli := strings.Split(httpVersion, "/")

	protocol := httpVersionSli[0]
	if protocol != "HTTP" {
		return fmt.Errorf("invalid http version")
	}

	versionNum := httpVersionSli[1]
	if versionNum == "" {
		return fmt.Errorf("invalid http version")
	}

	isVersionSupported := slices.Contains(supportedHttpVersions, versionNum)
	if !isVersionSupported {
		return fmt.Errorf("http version not supported")
	}

	reqLine.HttpVersion = versionNum
	req.RequestLine = reqLine

	return nil
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
