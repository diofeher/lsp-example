package lsp

type TextDocumentCodeActionRequest struct {
	Request
	Params TextDocumentCodeActionParams `json:"params"`
}

type TextDocumentCodeActionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        TextDocumentRange      `json:"range"`
	Context      CodeActionContext      `json:"context"`
}

type CodeActionContext struct {
}

type TextDocumentCodeActionResponse struct {
	Response
	Result []CodeAction `json:"result"`
}

type CodeAction struct {
	Title   string         `json:"title"`
	Edit    *WorkspaceEdit `json:"edit,omitempty"`
	Command *Command       `json:"command,omitempty"`
}

type Command struct {
	Title     string `json:"title"`
	Command   string `json:"command"`
	Arguments []any  `json:"arguments,omitempty"`
}

type WorkspaceEdit struct {
	Changes map[string][]TextEdit `json:"changes"`
}

type TextEdit struct {
	Range   TextDocumentRange `json:"range"`
	NewText string            `json:"newText"`
}
