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

func DecodeMessage(msg []byte) (int, error) {
	header, content, found := bytes.Cut(msg, []byte("\r\n\r\n"))
	if !found {
		return 0, fmt.Errorf("invalid message")
	}

	headerLen := len("Content-Length: ")
	contentLengthBytes := header[headerLen:]

	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, fmt.Errorf("error when parsing content length")
	}

	if len(content) != contentLength {
		return 0, fmt.Errorf("invalid content length")
	}

	return contentLength, nil
}
