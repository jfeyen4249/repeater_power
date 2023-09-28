package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// ComPort represents a serial port.
type ComPort struct {
	Name string
}

// SendCommand sends a command to the serial port.
func (c *ComPort) SendCommand(commandHex string) error {
	cleanup := func(port *serial.Port) {
		port.Flush()
		port.Close()
	}

	cPort, err := serial.OpenPort(&serial.Config{
		Name:        c.Name,
		Baud:        9600,
		ReadTimeout: time.Second,
	})
	if err != nil {
		return err
	}
	defer cleanup(cPort)

	_, err = cPort.Write(hexStringToBytes(commandHex))
	if err != nil {
		return err
	}

	return nil
}

func hexStringToBytes(s string) []byte {
	s = strings.ReplaceAll(s, " ", "")
	bytes, _ := hexStringToBytesWithError(s)
	return bytes
}

func hexStringToBytesWithError(s string) ([]byte, error) {
	bytes := make([]byte, 0, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		b, err := hexByteWithError(s[i : i+2])
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, b)
	}
	return bytes, nil
}

func hexByteWithError(s string) (byte, error) {
	val, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		return 0, err
	}
	return byte(val), nil
}

// Configuration represents the application configuration.
type Configuration struct {
	ACPortName string // COM port for AC
	DCPortName string // COM port for DC
	ACMode     bool   // AC mode (true for ON, false for NC)
	DCMode     bool   // DC mode (true for ON, false for NC)
}

const configFileName = "config.json"

func loadConfiguration() (Configuration, error) {
	var config Configuration
	file, err := os.Open(configFileName)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func saveConfiguration(config Configuration) error {
	file, err := os.Create(configFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	return err
}

func main() {
	config, err := loadConfiguration()
	if err != nil {
		config = Configuration{
			ACPortName: "COM2",
			DCPortName: "COM3",
			ACMode:     true, // Default to NO mode for AC
			DCMode:     true, // Default to NO mode for DC
		}
	}

	for {
		clearScreen()
		fmt.Println("")
		fmt.Println("")
		fmt.Println("  |-------------------------------------|")
		fmt.Println("  |   Repeater System Power Control App |")
		fmt.Println("  |                  v1.1               |")
		fmt.Println("  |-------------------------------------|")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("   1. AC Power Off")
		fmt.Println("   2. AC Power On")
		fmt.Println("   -----------------")
		fmt.Println("   3. DC Power Off")
		fmt.Println("   4. DC Power On")
		fmt.Println("   -----------------")
		fmt.Println("   5. AC Power Cycle")
		fmt.Println("   6. DC Power Cycle")
		fmt.Println("   -----------------")
		fmt.Println("   7. Change COM Ports")
		fmt.Println("   8. Change Mode (NO/NC) for AC")
		fmt.Println("   9. Change Mode (NO/NC) for DC")
		fmt.Println("   C. Current Configuration")
		fmt.Println("   I. Software Info")
		fmt.Println("   E. Exit")
		fmt.Println("")
		fmt.Print("  Select an option: ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			clearScreen()
			comPort := ComPort{Name: config.ACPortName}
			if config.ACMode {
				comPort.SendCommand("A0 01 01 A2")
			} else {
				comPort.SendCommand("A0 01 00 A1")
			}
			fmt.Println("AC Power Off")
		case "2":
			clearScreen()
			comPort := ComPort{Name: config.ACPortName}
			if config.ACMode {
				comPort.SendCommand("A0 01 00 A1")
			} else {
				comPort.SendCommand("A0 01 01 A2")
			}
			fmt.Println("AC Power On")
		case "3":
			clearScreen()
			comPort := ComPort{Name: config.DCPortName}
			if config.DCMode {
				comPort.SendCommand("A0 01 01 A2")
			} else {
				comPort.SendCommand("A0 01 00 A1")
			}
			fmt.Println("DC Power Off")
		case "4":
			clearScreen()
			comPort := ComPort{Name: config.DCPortName}
			if config.DCMode {
				comPort.SendCommand("A0 01 00 A1")
			} else {
				comPort.SendCommand("A0 01 01 A2")
			}
			fmt.Println("DC Power On")
		case "5":
			clearScreen()
			fmt.Println("")
		fmt.Println("")
			fmt.Println("|-------------------------------------|")
			fmt.Println("|        Power Cycling AC Power       |")
			fmt.Println("|             Please wait...          |")
			fmt.Println("|-------------------------------------|")
			comPort := ComPort{Name: config.ACPortName}
			if config.ACMode {
				comPort.SendCommand("A0 01 01 A2")
				fmt.Println("Cycling AC Power (Turned Off)")
				time.Sleep(2 * time.Second)
				comPort.SendCommand("A0 01 00 A1")
				fmt.Println("Cycling AC Power (Turned On)")
			} else {
				comPort.SendCommand("A0 01 00 A1")
				fmt.Println("Cycling AC Power (Turned Off)")
				time.Sleep(2 * time.Second)
				comPort.SendCommand("A0 01 01 A2")
				fmt.Println("Cycling AC Power (Turned On)")
			}
			time.Sleep(2 * time.Second)
		case "6":
			clearScreen()
			fmt.Println("")
		fmt.Println("")
			fmt.Println("  |-------------------------------------|")
			fmt.Println("  |        Power Cycling DC Power       |")
			fmt.Println("  |             Please wait...          |")
			fmt.Println("  |-------------------------------------|")
			comPort := ComPort{Name: config.DCPortName}
			if config.DCMode {
				comPort.SendCommand("A0 01 01 A2")
				fmt.Println("Cycling DC Power (Turned Off)")
				time.Sleep(2 * time.Second)
				comPort.SendCommand("A0 01 00 A1")
				fmt.Println("Cycling DC Power (Turned On)")
			} else {
				comPort.SendCommand("A0 01 00 A1")
				fmt.Println("Cycling DC Power (Turned Off)")
				time.Sleep(2 * time.Second)
				comPort.SendCommand("A0 01 01 A2")
				fmt.Println("Cycling DC Power (Turned On)")
			}
			time.Sleep(2 * time.Second)


		case "7":
			clearScreen()
			fmt.Println("")
		fmt.Println("")
			fmt.Print("Enter new AC COM Port (e.g., COM4, or press Enter to keep it unchanged (", config.ACPortName,")): ")
			acPort, _ := reader.ReadString('\n')
			acPort = strings.TrimSpace(acPort)
			if acPort != "" {
				config.ACPortName = acPort
			}
			fmt.Println("")
			fmt.Println("")
			fmt.Print("Enter new DC COM Port (e.g., COM6, or press Enter to keep it unchanged (", config.DCPortName,")): ")
			dcPort, _ := reader.ReadString('\n')
			dcPort = strings.TrimSpace(dcPort)
			if dcPort != "" {
				config.DCPortName = dcPort
			}

			err := saveConfiguration(config)
			if err != nil {
				fmt.Println("Error saving configuration:", err)
			}
			fmt.Println("")

			fmt.Println("COM Ports changed and saved.")

		case "8":
			clearScreen()
			fmt.Println("")
			fmt.Println("AC Relay Connections ")
			fmt.Println("--------------------")
			fmt.Println("NC if the relay is wired on NC connector  | NC if the relay is wired on NO connector ")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("| --------- |")
			fmt.Println("| |       | |")
			fmt.Println("| | Relay | |")
			fmt.Println("| |       | |")
			fmt.Println("| --------- |")
			fmt.Println("| |X||X||X| |")
			fmt.Println("             ")
			fmt.Println("   N  C  N   ")
			fmt.Println("   O  O  C   ")
			fmt.Println("      M      ")
			fmt.Println("------------------------")
			fmt.Print("Select mode for AC (no/nc): ")

			modeInput, _ := reader.ReadString('\n')
			modeInput = strings.TrimSpace(modeInput)
			if modeInput == "no" {
				config.ACMode = true
			} else if modeInput == "nc" {
				config.ACMode = false
			} else {
				fmt.Println("Invalid mode. Please enter 'no' or 'nc'.")
				time.Sleep(2 * time.Second)
				continue
			}

			err := saveConfiguration(config)
			if err != nil {
				fmt.Println("Error saving configuration:", err)
			}
			fmt.Println("")
			fmt.Println("AC Mode changed and saved.")
		case "9":
			clearScreen()
			fmt.Println("")
			fmt.Println("DC Relay Connections")
			fmt.Println("------------------- ")
			fmt.Println("NC if the relay is wired on NC connector  | NO if the relay is wired on NO connector ")
			fmt.Println("")
			fmt.Println("| --------- |")
			fmt.Println("| |       | |")
			fmt.Println("| | Relay | |")
			fmt.Println("| |       | |")
			fmt.Println("| --------- |")
			fmt.Println("| |X||X||X| |")
			fmt.Println("             ")
			fmt.Println("   N  C  N   ")
			fmt.Println("   O  O  C   ")
			fmt.Println("      M      ")
			fmt.Println("------------------------")
			fmt.Print("Select mode for DC (NO/NC): ")
			modeInput, _ := reader.ReadString('\n')
			modeInput = strings.TrimSpace(modeInput)
			if modeInput == "no" {
				config.DCMode = true
			} else if modeInput == "nc" {
				config.DCMode = false
			} else {
				fmt.Println("Invalid mode. Please enter 'no' or 'nc'.")
				time.Sleep(2 * time.Second)
				continue
			}

			err := saveConfiguration(config)
			if err != nil {
				fmt.Println("Error saving configuration:", err)
			}
			fmt.Println("")
			fmt.Println("DC Mode changed and saved.")
		case "c":
			clearScreen()
			fmt.Println("")
			fmt.Println("")
			fmt.Println("   |---------------|-----------------|---------------|")
			fmt.Println("   |---------------|Relay Config Mode|    Com Port   |")
			fmt.Println("   |---------------|-----------------|---------------|")
			fmt.Println("   |   AC Relay    |", getModeText(config.ACMode), " |",config.ACPortName,"         |")
			fmt.Println("   |---------------|-----------------|---------------|")
			fmt.Println("   |   DC Relay    |", getModeText(config.DCMode), " |",config.DCPortName,"         |")
			fmt.Println("   |---------------|-----------------|---------------|")

			fmt.Println(" ")
			fmt.Println(" ")
			fmt.Println("Press Enter to continue...")
			reader.ReadString('\n')
		case "i":
			clearScreen()
			fmt.Println("")
			fmt.Println("")
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  |      Repeater System Power Control App       |")
			fmt.Println("  |                      v1.1                    |")
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  |               Written By: KD9HAE             |")
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  |               Written In: GoLang             |")
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  |             Written: Sept: 23rd 2023         |")
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  |             Updated: Sept: 27th 2023         |")
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  | https://github.com/jfeyen4249/repeater_power |")
			fmt.Println("  |----------------------------------------------|")
			time.Sleep(6 * time.Second)
		case "e":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("")
			fmt.Println("")
			fmt.Println("Invalid option. Please select a valid option (1-9, e to exit).")
		}

		time.Sleep(1 * time.Second)
	}
}

func clearScreen() {
	var cmd *exec.Cmd

	if isWindows() {
		cmd = exec.Command("cmd", "/c", "cls") // Windows
	} else {
		cmd = exec.Command("clear") // Unix-based systems (including macOS and Linux)
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func isWindows() bool {
	return os.Getenv("OS") == "Windows_NT"
}

func getModeText(mode bool) string {
	if mode {
		return "Normally Open "
	}
	return "Normally Close"
}