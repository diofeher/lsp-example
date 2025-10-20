package lsp

type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     TextDocumentPosition   `json:"position"`
}

type TextDocumentPosition struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type TextDocumentLocation struct {
	URI   string            `json:"uri"`
	Range TextDocumentRange `json:"range"`
}

type TextDocumentRange struct {
	Start TextDocumentPosition `json:"start"`
	End   TextDocumentPosition `json:"end"`
}
