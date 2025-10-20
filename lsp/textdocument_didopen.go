package lsp

type TextDocumentDidOpenNotification struct {
	Notification
	Params DidOpenTextDocumentParams `json:"params"`
}

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

func NewTextDocumentDidOpenNotification(params DidOpenTextDocumentParams) *TextDocumentDidOpenNotification {
	return &TextDocumentDidOpenNotification{
		Notification: Notification{
			RPC:    "2.0",
			Method: "textDocument/didOpen",
		},
	}
}
