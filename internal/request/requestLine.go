package request

import (
	"errors"
	"regexp"
	"strings"

	"slices"
)

var (
	// ref: https://www.rfc-editor.org/rfc/rfc9110#section-9
	supportedMethods = []string{
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

	supportedHttpVersions = []string{
		"1.1",
	}

	// Ref: https://datatracker.ietf.org/doc/html/rfc9112#name-request-target
	// Origin Form: absolute-path [ "?" query ]
	originFormRegex = regexp.MustCompile(`^\/[^\s]*(?:\?[^\s]*)?$`)

	// Absolute Form: absolute-URI
	absoluteFormRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*:\/\/[^\s]+$`)

	// Authority Form: host:port
	authorityFormRegex = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*|(?:[0-9]{1,3}\.){3}[0-9]{1,3}):[0-9]+$`)

	ERR_MALFORMED_REQUEST_LINE = errors.New("malformed request line")
	ERR_INVALID_METHOD         = errors.New("invalid method")
	ERR_INVALID_HTTP_VERSION   = errors.New("invalid http version")
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func isTargetValid(target string, method string) error {
	ok := originFormRegex.MatchString(target)
	if ok {
		return nil
	}

	ok = absoluteFormRegex.MatchString(target)
	if ok {
		return nil
	}

	ok = authorityFormRegex.MatchString(target)
	if ok {
		return nil
	}

	if method == "OPTIONS" && target == "*" {
		return nil
	}

	return errors.New("invalid request target")
}

func parseRequestLine(raw string, req *Request) error {
	parts := strings.Split(raw, " ")
	if len(parts) != 3 {
		return ERR_MALFORMED_REQUEST_LINE
	}

	reqLine := RequestLine{}

	method := parts[0]
	if method == "" {
		return ERR_INVALID_METHOD
	}

	isMethodSupported := slices.Contains(supportedMethods, method)
	if !isMethodSupported {
		return ERR_INVALID_METHOD
	}

	reqLine.Method = method

	target := parts[1]
	if target == "" {
		return ERR_MALFORMED_REQUEST_LINE
	}

	err := isTargetValid(target, method)
	if err != nil {
		return errors.Join(err, ERR_MALFORMED_REQUEST_LINE)
	}

	reqLine.RequestTarget = target

	httpVersion := parts[2]
	if httpVersion == "" {
		return ERR_INVALID_HTTP_VERSION
	}

	httpVersionSli := strings.Split(httpVersion, "/")

	protocol := httpVersionSli[0]
	if protocol != "HTTP" {
		return ERR_INVALID_HTTP_VERSION
	}

	versionNum := httpVersionSli[1]
	if versionNum == "" {
		return ERR_INVALID_HTTP_VERSION
	}

	isVersionSupported := slices.Contains(supportedHttpVersions, versionNum)
	if !isVersionSupported {
		return ERR_INVALID_HTTP_VERSION
	}

	reqLine.HttpVersion = versionNum
	req.RequestLine = reqLine

	return nil
}
