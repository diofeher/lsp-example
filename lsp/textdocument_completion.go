package lsp

type TextDocumentCompletionRequest struct {
	Request
	Params TextDocumentCompletionParams `json:"params"`
}

type TextDocumentCompletionParams struct {
	TextDocumentPositionParams
}

type TextDocumentCompletionResponse struct {
	Response
	Result []TextDocumentCompletionItem `json:"result"`
}

type TextDocumentCompletionItem struct {
	Label         string `json:"label"`
	Detail        string `json:"detail"`
	Documentation string `json:"documentation"`
}
