package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "kube-context-switcher",
		Short: "A tool to switch between Kubernetes contexts",
		Run:   switchContext,
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}

func switchContext(cmd *cobra.Command, args []string) {
	contexts, err := getKubeContexts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing contexts: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Available contexts:")
	for i, context := range contexts {
		fmt.Printf("%d. %s\n", i+1, context)
	}

	// Wait for user input
	fmt.Println("Use arrow keys to select a context and press enter to confirm:")
	selectedIndex := 0
	err = keyboard.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening keyboard: %v\n", err)
		os.Exit(1)
	}
	defer keyboard.Close()

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading keyboard input: %v\n", err)
			os.Exit(1)
		}

		if key == keyboard.KeyArrowUp {
			if selectedIndex > 0 {
				selectedIndex--
			}
		} else if key == keyboard.KeyArrowDown {
			if selectedIndex < len(contexts)-1 {
				selectedIndex++
			}
		} else if key == keyboard.KeyEnter {
			break
		}

		// Clear console and print available contexts with selected context highlighted
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		fmt.Println("Available contexts:")
		for i, context := range contexts {
			if i == selectedIndex {
				fmt.Printf("\033[1m%d. %s\033[0m\n", i+1, context)
			} else {
				fmt.Printf("%d. %s\n", i+1, context)
			}
		}
	}

	selectedContext := contexts[selectedIndex]
	err = setKubeContext(selectedContext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting context: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Current context set to %q\n", selectedContext)
}

func getKubeContexts() ([]string, error) {
	cmd := exec.Command("kubectl", "config", "get-contexts", "-o", "name")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get contexts: %v", err)
	}
	contexts := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(contexts) == 0 {
		return nil, fmt.Errorf("no contexts found")
	}
	return contexts, nil
}

func setKubeContext(context string) error {
	cmd := exec.Command("kubectl", "config", "use-context", context)
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 1 {
					return fmt.Errorf("context %q not found", context)
				}
			}
		}
		return fmt.Errorf("failed to set context: %v", err)
	}
	return nil
}
