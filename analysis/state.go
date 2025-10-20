package analysis

import (
	"fmt"
	"log"
	"strings"

	"github.com/diofeher/lspexample/lsp"
)

type State struct {
	// map of filenames to their content
	Documents map[string]string
}

func NewState() *State {
	return &State{
		Documents: make(map[string]string),
	}
}

func (s *State) OpenDocument(logger *log.Logger, document, text string) []lsp.Diagnostic {
	s.Documents[document] = text

	return getDiagnosticsForFile(logger, text)
}

func (s *State) UpdateDocument(logger *log.Logger, document, text string) []lsp.Diagnostic {
	s.Documents[document] = text

	return getDiagnosticsForFile(logger, text)
}

func (s *State) Hover(id int, document string, position int) lsp.TextDocumentHoverResponse {

	documentContent, ok := s.Documents[document]
	if !ok {
		return lsp.TextDocumentHoverResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &id,
			},
		}
	}
	return lsp.TextDocumentHoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.TextDocumentHoverResult{
			Contents: fmt.Sprintf("Hover: %d\n", len(documentContent)),
		},
	}
}

func (s *State) Definition(id int, document string, position lsp.TextDocumentPosition) lsp.TextDocumentDefinitionResponse {

	return lsp.TextDocumentDefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.TextDocumentLocation{
			URI: document,
			Range: lsp.TextDocumentRange{
				Start: lsp.TextDocumentPosition{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.TextDocumentPosition{
					Line:      position.Line - 1,
					Character: 3,
				},
			},
		},
	}
}

func (s *State) CodeAction(id int, uri string) lsp.TextDocumentCodeActionResponse {
	text := s.Documents[uri]

	actions := []lsp.CodeAction{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx > 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "Kate Editor",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS Code with Kate Editor",
				Edit: &lsp.WorkspaceEdit{
					Changes: replaceChange,
				},
			})

			censorChange := map[string][]lsp.TextEdit{}
			censorChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS C**de",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Censor VS Code",
				Edit: &lsp.WorkspaceEdit{
					Changes: censorChange,
				},
			})
		}
	}

	return lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.TextDocumentCompletionResponse {
	// text := s.Documents[uri]

	items := []lsp.TextDocumentCompletionItem{}
	items = append(items, lsp.TextDocumentCompletionItem{
		Label:         "VS Code",
		Detail:        "VS Code is a code editor",
		Documentation: "VS Code is a code editor",
	})

	items = append(items, lsp.TextDocumentCompletionItem{
		Label:         "Kate Editor",
		Detail:        "Kate Editor is a code editor",
		Documentation: "Kate Editor is a code editor",
	})

	return lsp.TextDocumentCompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}
}
func LineRange(line, start, end int) lsp.TextDocumentRange {
	return lsp.TextDocumentRange{
		Start: lsp.TextDocumentPosition{
			Line:      line,
			Character: start,
		},
	}
}

func getDiagnosticsForFile(logger *log.Logger, text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}

	for row, line := range strings.Split(text, "\n") {
		logger.Println(line)
		idx := strings.Index(line, "VS Code")
		if idx > -1 {
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("VS Code")),
				Message:  "VS Code is a code editor",
				Source:   "VS Code",
				Severity: 1,
			})
		}

		idx = strings.Index(line, "Kate Editor")
		if idx > -1 {
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("Kate Editor")),
				Message:  "Kate is a code editor",
				Source:   "Kate Editor",
				Severity: 2,
			})
		}
	}
	return diagnostics
}
