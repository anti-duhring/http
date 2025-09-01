package request

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"slices"

	"github.com/rotisserie/eris"
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

	err := isTargetValid(target, method)
	if err != nil {
		return eris.Wrap(err, "calling isTargetValid")
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
