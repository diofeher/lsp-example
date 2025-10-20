package rpc

// package rpc implements the RPC protocol for the LSP.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

type BaseMessage struct {
	Method string `json:"method"`
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte("\r\n\r\n"))
	if !found {
		return "", nil, fmt.Errorf("invalid message")
	}

	headerLen := len("Content-Length: ")
	contentLengthBytes := header[headerLen:]

	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, fmt.Errorf("error when parsing content length")
	}

	if len(content) != contentLength {
		return "", nil, fmt.Errorf("invalid content length")
	}

	var baseMessage BaseMessage
	err = json.Unmarshal(content, &baseMessage)
	if err != nil {
		return "", nil, fmt.Errorf("error when unmarshalling base message")
	}

	return baseMessage.Method, content[:contentLength], nil
}

func SplitFunc(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte("\r\n\r\n"))
	if !found {
		return 0, nil, nil
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, fmt.Errorf("error when parsing content length")
	}

	if len(content) < contentLength {
		return 0, nil, nil
	}

	totalLength := len(header) + 4 + contentLength
	return totalLength, data[:totalLength], nil
}
