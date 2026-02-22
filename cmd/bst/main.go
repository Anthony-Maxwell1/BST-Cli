package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Anthony-Maxwell1/BST-Cli/internal/daemon"
	"github.com/Anthony-Maxwell1/BST-Cli/internal/fetch"
	"github.com/Anthony-Maxwell1/BST-Cli/internal/ws"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bst <command>")
		return
	}

	switch os.Args[1] {

	case "fetch":
		if err := fetch.FetchLatest(); err != nil {
			fmt.Println("Error:", err)
		}

	case "run":
		if err := daemon.Run(); err != nil {
			fmt.Println("Error:", err)
		}

	case "stop":
		if err := daemon.Stop(); err != nil {
			fmt.Println("Error:", err)
		}

	case "project":
		handleProject(os.Args[2:])

	case "git":
		handleGit(os.Args[2:])

	default:
		fmt.Println("Unknown command")
	}
}

func handleProject(args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {

	case "list":
		ws.SendPacket(map[string]any{
			"type": "project_list",
		})

	case "open":
		if len(args) < 2 {
			fmt.Println("project open <name>")
			return
		}
		ws.SendPacket(map[string]any{
			"type": "project_open",
			"name": args[1],
		})

	case "close":
		ws.SendPacket(map[string]any{
			"type": "project_close",
		})
	}
}

func handleGit(args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {

	case "add":
		ws.SendPacket(map[string]any{
			"type": "git_add",
		})

	case "commit":
		if len(args) < 3 || args[1] != "-m" {
			fmt.Println(`git commit -m "message"`)
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
	}
}
