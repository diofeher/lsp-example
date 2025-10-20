package rpc_test

import (
	"testing"

	"github.com/diofeher/lspexample/rpc"

	"github.com/google/go-cmp/cmp"
)

func TestEncodeMessage(t *testing.T) {
	tests := []struct {
		Name    string
		Message any
		WantMsg string
	}{
		{
			Name: "test",
			Message: struct {
				ID int
			}{
				ID: 1,
			},
			WantMsg: "Content-Length: 8\r\n\r\n{\"ID\":1}",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotMsg := rpc.EncodeMessage(test.Message)
			if diff := cmp.Diff(test.WantMsg, gotMsg); diff != "" {
				t.Errorf("EncodeMessage() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDecodeMessage(t *testing.T) {
	tests := []struct {
		Name       string
		Msg        string
		WantSize   int
		WantMethod string
	}{
		{
			Name:       "test",
			Msg:        "Content-Length: 17\r\n\r\n{\"Method\":\"test\"}",
			WantSize:   17,
			WantMethod: "test",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			gotMethod, gotSize, err := rpc.DecodeMessage([]byte(test.Msg))
			if err != nil {
				t.Fatalf("DecodeMessage() error: %d: %v", gotSize, err)
			}

			if diff := cmp.Diff(test.WantSize, gotSize); diff != "" {
				t.Errorf("DecodeMessage() mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(test.WantMethod, gotMethod); diff != "" {
				t.Errorf("DecodeMessage() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
