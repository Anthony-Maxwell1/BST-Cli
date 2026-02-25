package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Anthony-Maxwell1/BST-Cli/internal/daemon"
	"github.com/Anthony-Maxwell1/BST-Cli/internal/fetch"
	"github.com/Anthony-Maxwell1/BST-Cli/internal/ws"
	"github.com/google/uuid"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bst <command>")
		return
	}

	switch os.Args[1] {

	case "fetch":
		if err := fetch.FetchLatest(); err != nil {
			fmt.Println("Error fetching latest release:", err)
		} else {
			fmt.Println("Fetched and extracted latest release successfully.")
		}

	case "run":
		if err := daemon.Run(); err != nil {
			fmt.Println("Error starting daemon:", err)
		}

	case "stop":
		if err := daemon.Stop(); err != nil {
			fmt.Println("Error stopping daemon:", err)
		}

	case "project":
		handleProject(os.Args[2:])

	case "git":
		handleGit(os.Args[2:])

	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}

func handleProject(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: bst project <list|open|close>")
		return
	}

	switch args[0] {
	case "list":
		id := uuid.New().String()
		ws.SendPacket(map[string]any{
			"type":    "cli",
			"command": "list-projects",
			"id":      id,
		})

	case "open":
		if len(args) < 2 {
			fmt.Println("Usage: bst project open <name>")
			return
		}
		ws.SendPacket(map[string]any{
			"type":    "cli",
			"command": "open-project",
			"args": map[string]any{
				"name": args[1],
			},
		})
	case "close":
		ws.SendPacket(map[string]any{
			"type":    "cli",
			"command": "close-project",
		})

	default:
		fmt.Println("Unknown project command:", args[0])
	}
}

func handleGit(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: bst git <add|commit|push|pull>")
		return
	}

	switch args[0] {
	case "add":
		ws.SendPacket(map[string]any{
			"type": "git_add",
		})

	case "commit":
		if len(args) < 3 || args[1] != "-m" {
			fmt.Println(`Usage: bst git commit -m "message"`)
			return
		}
		message := strings.Join(args[2:], " ")
		ws.SendPacket(map[string]any{
			"type":    "git_commit",
			"message": message,
		})

	case "push":
		ws.SendPacket(map[string]any{
			"type": "git_push",
		})

	case "pull":
		ws.SendPacket(map[string]any{
			"type": "git_pull",
		})

	default:
		fmt.Println("Unknown git command:", args[0])
	}
}
