package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	exprPtr := flag.String("expr", "", "Mathematical expression to evaluate (e.g., \"2 + 3 * sqrt(4)\")")
	serverPtr := flag.Bool("server", false, "Start the API server")
	portPtr := flag.String("port", "8080", "Port for the API server")

	flag.Parse()

	if *serverPtr {
		// Server mode
		startServer(*portPtr)
	} else if *exprPtr != "" {
		// CLI mode
		result, err := Calculate(*exprPtr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result)
	} else {
		// Interactive mode
		interactiveMode()
	}
}

func interactiveMode() {
	fmt.Println("Calculator Interactive Mode")
	fmt.Println("Enter an expression (or 'exit' to quit):")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}
		if input == "" {
			continue
		}

		res, err := Calculate(input)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(res)
		}
	}
}

// API Request/Response structures
type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func startServer(port string) {
	http.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CalcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		result, err := Calculate(req.Expression)
		resp := CalcResponse{Result: result}
		if err != nil {
			resp.Error = err.Error()
			// Depending on requirements, we might want to send 400 Bad Request if formula is wrong
			// but keeping 200 with error field is also common for simple JSON RPC style
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
		os.Exit(1)
	}
}
