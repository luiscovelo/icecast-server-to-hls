package parser

import (
	"regexp"
	"strings"
)

var regex *regexp.Regexp

type Request struct {
	Method string
	Param  string
	Header map[string]string
	Body   []byte
}

func init() {
	pattern := `^(PUT|GET) /([^ ]+) HTTP/1\.1$`
	regex = regexp.MustCompile(pattern)
}

func Parse(data []byte) *Request {
	req := &Request{
		Method: "UNKNOWN",
		Param:  "",
		Header: make(map[string]string),
		Body:   data,
	}

	lines := strings.Split(string(data), "\r\n")
	match := regex.FindStringSubmatch(lines[0])

	if match == nil {
		return req
	}

	method := strings.TrimSpace(match[1])

	switch method {
	case "PUT":
		req.Method = method
		req.Param = match[2]

		for i, line := range lines {
			if i > 0 && len(line) > 0 {
				content := strings.Split(line, ":")

				key := strings.TrimSpace(content[0])
				value := strings.TrimSpace(content[1])

				req.Header[key] = value
			}
		}
	case "GET":
		req.Method = method
		req.Param = match[2]
	}

	return req
}
