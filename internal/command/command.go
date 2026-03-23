package command

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/peterh/liner"

	"github.com/ESP32-Zephyr/esp32_zephyr_goapi/api"
)

const appPort = 4242
var historyPath string

var rootCmd = &cobra.Command{
	Use:   "esp32_shell",
	Short: "Interactive ESP32 shell",
	Long:  `Interactive ESP32 shell`,
	Run: func(cmd *cobra.Command, args []string) {
		shellLoop(cmd)
	},
}

type shellState struct {
	prompt string
	state string
	greeting string 
	es32client *api.Esp32Client
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

func (s *shellState) setState (state string) {
	s.state = state
	if state == "connected" {
		s.prompt = s.es32client.Ipv4 + " > "
	} else {
		s.prompt = "> "
	}
}

func (s *shellState) getPrompt () string {
	return s.prompt
}

func (s *shellState) sendPing() bool {
	var success = false

	pong, err := s.es32client.Ping()
	if err != nil {
		fmt.Println("Error:", err)
		s.setState("disconnected")
	} else {
		pong := pong.GetPong()
		if pong == "pong" {
			s.setState("connected")
			success = true
		}
	}

	return success
}

func (s *shellState) connect (transport, ipv4 string, destPort uint16) {
	s.es32client, _ = api.NewEsp32Client(transport, ipv4, destPort)
	s.sendPing()
}

func (s *shellState) adc_chs_get () uint32 {
	if s.sendPing() {
		res, err := s.es32client.AdcChsGet()
		if err != nil {
			fmt.Println("Error:", err)
		}
		return res.GetAdcChs()
	}

	return 0
}

func Init() {
	homeDir, _ := os.UserHomeDir()
	historyPath = filepath.Join(homeDir, ".esp32_shell_history")
		
	rootCmd.AddCommand(testCmd())
	rootCmd.AddCommand(connectCmd())
	rootCmd.AddCommand(adcCmd())	
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

func connectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connect",
		Short: "connect to esp32",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				shell.connect(args[0], args[1], appPort)
			}
		},
	}
}

func adcCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "adc",
		Short: "ADC commands",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				channels := shell.adc_chs_get()
				fmt.Printf("%d\n", channels)
			} 
		},
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

