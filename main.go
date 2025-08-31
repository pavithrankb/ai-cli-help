package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// App metadata
const (
	AppName    = "cli-help"
	AppVersion = "0.5.2"
	AppAuthor  = "Pavithran KB"
)

var unsafeMode bool = false
var aiMode bool = false
var verbose bool = false

// getFriendlyOS returns a human-readable OS string
func getFriendlyOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "linux":
		// Try to get distribution from /etc/os-release
		if data, err := os.ReadFile("/etc/os-release"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					return strings.Trim(line[len("PRETTY_NAME="):], `"`)
				}
			}
		}
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return runtime.GOOS
	}
}

// Ask AI using Bearer token + Bedrock HTTP API
func askAI(prompt string) (string, error) {
	apiKey := os.Getenv("AWS_BEARER_TOKEN_BEDROCK")
	if apiKey == "" {
		return "", fmt.Errorf("missing AWS_BEARER_TOKEN_BEDROCK env var")
	}

	modelID := os.Getenv("BEDROCK_MODEL_ID")
	if modelID == "" {
		modelID = "anthropic.claude-3-haiku-20240307-v1:0"
	}

	url := fmt.Sprintf("https://bedrock-runtime.us-east-1.amazonaws.com/model/%s/invoke", modelID)

	arch := runtime.GOARCH
	osFriendly := getFriendlyOS()
	systemInfo := fmt.Sprintf("System Info: OS=%s Arch=%s\n", osFriendly, arch)

	fullPrompt := systemInfo + prompt + `
Respond only with the shell command(s).
If multiple commands are needed, return them as a multi-line script.
Do NOT include explanations.`

	payload := map[string]interface{}{
		"anthropic_version": "bedrock-2023-05-31",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": fullPrompt,
			},
		},
		"max_tokens": 300,
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(context.TODO(), "POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error: %s", string(b))
	}

	if verbose {
		fmt.Println("📜 Raw AI response:\n", string(b))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		return "", err
	}

	// Claude returns output inside "content" array
	if contentArr, ok := result["content"].([]interface{}); ok && len(contentArr) > 0 {
		if first, ok := contentArr[0].(map[string]interface{}); ok {
			if text, ok := first["text"].(string); ok {
				return strings.TrimSpace(text), nil
			}
		}
	}

	return "", fmt.Errorf("unexpected AI response format: %+v", result)
}

func printHelp() {
	fmt.Printf(`
%s - AI-powered CLI helper tool
Author: %s
Version: %s

Usage:
  %s [options]

Options:
  --help       Show this help message
  --version    Show version info
  --unsafe     Enable UNSAFE MODE (allows 'rm' commands ⚠️ USE WITH CARE)
  --ai         Enable AI mode (translate natural language into CLI commands)
  --verbose    Print raw AI response JSON (for debugging)

Description:
  %s helps you translate natural language into CLI commands.
  Safe mode (default) blocks dangerous commands like 'rm'.
  Unsafe mode explicitly allows them.
  AI mode uses Amazon Bedrock Anthropic Claude to suggest commands.
  The AI now automatically receives your OS and architecture to generate compatible commands.
`, AppName, AppAuthor, AppVersion, AppName, AppName)
}

func main() {
	args := os.Args[1:]

	// Handle flags
	for _, arg := range args {
		switch arg {
		case "--help", "-h":
			printHelp()
			return
		case "--version", "-v":
			fmt.Printf("%s v%s by %s\n", AppName, AppVersion, AppAuthor)
			return
		case "--unsafe":
			unsafeMode = true
			fmt.Println("⚠️ UNSAFE MODE enabled! 'rm' commands are now allowed.")
		case "--ai":
			aiMode = true
			fmt.Println("🤖 AI mode enabled! Using Bedrock Anthropic Claude (Bearer token).")
		case "--verbose":
			verbose = true
			fmt.Println("🔍 Verbose mode: Raw AI responses will be printed.")
		}
	}

	fmt.Println("Welcome to CLI Helper 🚀")
	if unsafeMode {
		fmt.Println("⚠️ Running in UNSAFE MODE. Dangerous commands are allowed.")
	} else {
		fmt.Println("✅ Running in SAFE MODE. 'rm' commands are blocked.")
	}
	if aiMode {
		fmt.Println("🤖 AI Mode is ON: Type natural language and I'll translate to commands.")
	}
	fmt.Println("Type your command (or natural language request). Type 'exit' to quit.")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Goodbye 👋")
			break
		}

		var command string
		if aiMode {
			cmdOut, err := askAI(input)
			if err != nil {
				fmt.Printf("❌ AI error: %v\n", err)
				continue
			}
			fmt.Printf("🤖 Suggested command:\n%s\n", cmdOut)
			fmt.Print("Run this? (y/n): ")
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(confirm)
			if confirm != "y" {
				fmt.Println("❌ Skipped.")
				continue
			}
			command = cmdOut
		} else {
			command = input
		}

		// Safety filter
		if !unsafeMode && (strings.HasPrefix(command, "rm ") || strings.Contains(command, " rm ") || command == "rm") {
			fmt.Println("❌ Error: 'rm' commands are not allowed in SAFE MODE.")
			continue
		}

		cmd := exec.Command("bash", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err := cmd.Run()
		if err != nil {
			fmt.Printf("⚠️ Command failed: %v\n", err)
		}
	}
}
