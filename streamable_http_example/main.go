package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type HiParams struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type ArithmeticParams struct {
	A float64 `json:"a" jsonschema:"first number"`
	B float64 `json:"b" jsonschema:"second number"`
}

func SayHi(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[HiParams]) (*mcp.CallToolResultFor[any], error) {
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: "Hi " + params.Arguments.Name}},
	}, nil
}

func Add(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ArithmeticParams]) (*mcp.CallToolResultFor[any], error) {
	result := params.Arguments.A + params.Arguments.B
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: "Result: " + fmt.Sprintf("%.2f", result)}},
	}, nil
}

func Subtract(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ArithmeticParams]) (*mcp.CallToolResultFor[any], error) {
	result := params.Arguments.A - params.Arguments.B
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: "Result: " + fmt.Sprintf("%.2f", result)}},
	}, nil
}

func Multiply(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ArithmeticParams]) (*mcp.CallToolResultFor[any], error) {
	result := params.Arguments.A * params.Arguments.B
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: "Result: " + fmt.Sprintf("%.2f", result)}},
	}, nil
}

func Divide(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ArithmeticParams]) (*mcp.CallToolResultFor[any], error) {
	if params.Arguments.B == 0 {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Error: Division by zero"}},
		}, nil
	}
	result := params.Arguments.A / params.Arguments.B
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: "Result: " + fmt.Sprintf("%.2f", result)}},
	}, nil
}

func main() {
	// Create a server with multiple tools.
	server := mcp.NewServer(&mcp.Implementation{Name: "calculator", Version: "v1.0.0"}, nil)

	// Add greeting tool
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi to someone"}, SayHi)

	// Add arithmetic tools
	mcp.AddTool(server, &mcp.Tool{Name: "add", Description: "add two numbers"}, Add)
	mcp.AddTool(server, &mcp.Tool{Name: "subtract", Description: "subtract second number from first"}, Subtract)
	mcp.AddTool(server, &mcp.Tool{Name: "multiply", Description: "multiply two numbers"}, Multiply)
	mcp.AddTool(server, &mcp.Tool{Name: "divide", Description: "divide first number by second"}, Divide)

	// Create a streamable HTTP handler
	handler := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{},
	)

	// Set up HTTP server
	http.HandleFunc("/mcp", handler.ServeHTTP)

	log.Println("Starting MCP calculator server on :8080")
	log.Println("Connect to http://localhost:8080/mcp")
	log.Println("Available tools: greet, add, subtract, multiply, divide")

	// Run the HTTP server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
