package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/diofeher/lspexample/analysis"
	"github.com/diofeher/lspexample/lsp"
	"github.com/diofeher/lspexample/rpc"
)

func main() {
	fmt.Println("Starting")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	filePath := dir + "/log.txt"
	logger := getLogger(filePath)
	logger.Println("Started logger..")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.SplitFunc)

	state := analysis.NewState()
	writer := os.Stdout

	logger.Println("Waiting for messages...")

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s\n", err)
			continue
		}
		logger.Printf("Received method: %v\n", method)
		handleMessage(logger, writer, state, method, contents)
	}

	logger.Println("Finished. ")
	for {
		time.Sleep(1 * time.Second)
		logger.Println("Still waiting for messages...")
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state *analysis.State, method string, contents []byte) {
	logger.Printf("Received method: %v\n", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		logger.Println(string(contents))
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error unmarshalling initialize request: %v\n", err)
			return
		}
		if request.Params.ClientInfo == nil {
			logger.Printf("Client info is nil\n")
		} else {
			logger.Printf("Connected to client: %s %s\n", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		}

		msg := lsp.NewInitializeResponse(request.ID)
		writeResponse(logger, writer, msg)

		logger.Println("Sent initialize response")
	case "initialized":
		logger.Println(string(contents))
		return
	case "textDocument/didOpen":
		var request lsp.TextDocumentDidOpenNotification
		logger.Println(string(contents))
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didOpen request: %v\n", err)
			return
		}
		logger.Printf("Opened file: %s %s\n", request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		diagnostics := state.OpenDocument(logger, request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		writeResponse(logger, writer, lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		})
		writeResponse(logger, writer, diagnostics)
		return
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		logger.Println(string(contents))
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didChange request: %v\n", err)
			return
		}
		logger.Printf("Changed: %s\n", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			diagnostics := state.UpdateDocument(logger, request.Params.TextDocument.URI, change.Text)

			writeResponse(logger, writer, lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})
		}
		return
	case "textDocument/hover":
		var request lsp.TextDocumentHoverRequest
		logger.Println(string(contents))
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/hover request: %v\n", err)
			return
		}

		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position.Line)
		writeResponse(logger, writer, response)
	case "textDocument/definition":
		var request lsp.TextDocumentHoverRequest
		logger.Println(string(contents))
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/hover request: %v\n", err)
			return
		}

		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(logger, writer, response)
	case "textDocument/codeAction":
		var request lsp.TextDocumentCodeActionRequest
		logger.Println(string(contents))
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/codeAction request: %v\n", err)
			return
		}

		response := state.CodeAction(request.ID, request.Params.TextDocument.URI)
		writeResponse(logger, writer, response)
	case "textDocument/completion":
		var request lsp.TextDocumentCompletionRequest
		logger.Println(string(contents))
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/completion request: %v\n", err)
			return
		}

		response := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI)
		writeResponse(logger, writer, response)
	}
}

func writeResponse(logger *log.Logger, writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	logger.Printf("Sending response: %s\n", reply)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)

	if err != nil {
		panic(err)
	}

	return log.New(logfile, "[lsp]", log.Ldate|log.Lshortfile|log.Ltime)
}
