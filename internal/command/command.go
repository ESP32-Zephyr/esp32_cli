package command

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/peterh/liner"
)

type shellState struct {
	prompt string
	greeting string 
}

var shell = shellState{
		prompt: "> ",
		greeting: `
███████╗███████╗██████╗ ██████╗ ██████╗         ███████╗██╗  ██╗███████╗██╗     ██╗     
██╔════╝██╔════╝██╔══██╗╚════██╗╚════██╗        ██╔════╝██║  ██║██╔════╝██║     ██║     
█████╗  ███████╗██████╔╝ █████╔╝ █████╔╝        ███████╗███████║█████╗  ██║     ██║     
██╔══╝  ╚════██║██╔═══╝  ╚═══██╗██╔═══╝         ╚════██║██╔══██║██╔══╝  ██║     ██║     
███████╗███████║██║     ██████╔╝███████╗        ███████║██║  ██║███████╗███████╗███████╗
╚══════╝╚══════╝╚═╝     ╚═════╝ ╚══════╝        ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝
`,
	}

func (s *shellState) greetingPrint () {
	fmt.Println(s.greeting)
}

func (s *shellState) setPrompt (state string) {
	s.prompt = state
}

func (s *shellState) getPrompt () string {
	return s.prompt
}

var historyPath string

var rootCmd = &cobra.Command{
		Use:   "ESP32",
		Short: "Interactive ESP32 shell",
		Long:  `Interactive ESP32 shell`,
		Run: func(cmd *cobra.Command, args []string) {
			shellLoop(cmd)
		},
	}

func Init() {
	homeDir, _ := os.UserHomeDir()
	historyPath = filepath.Join(homeDir, ".esp32_shell_history")
		
	rootCmd.AddCommand(testCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func shellLoop(root *cobra.Command) {
	lineReader := liner.NewLiner()
	defer lineReader.Close()

	// Load history
	if f, err := os.Open(historyPath); err == nil {
		lineReader.ReadHistory(f)
		f.Close()
	}

	shell.greetingPrint()
	for {
		input, _ := lineReader.Prompt(shell.getPrompt())
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		// Save to history
		lineReader.AppendHistory(input)
		if input == "exit" {
			break
		}

		// Split first word as command, rest as args
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}
		
		// Execute through Cobra
		cmdArgs := append([]string{os.Args[0]}, args...)
		root.SetArgs(cmdArgs[1:]) // Skip binary name
		if err := root.Execute(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
	}

	// Save history
	if f, err := os.Create(historyPath); err == nil {
		lineReader.WriteHistory(f)
		f.Close()
	}	
}

func testCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test command",
		Short: "test command",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			arg := ""
			if len(args) > 0 {
				arg = args[0]
			}
			fmt.Printf("Argument: %s\n", arg)
		},
	}
}

