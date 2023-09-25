package main

import (
	"bufio"
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

func main() {
	com5 := ComPort{Name: "COM5"}
	com6 := ComPort{Name: "COM6"}

	for {
		clearScreen() // Clear the screen
		fmt.Println("")
		fmt.Println("")
		fmt.Println("  |-------------------------------------|")
		fmt.Println("  |   Repeater System Power Control App |")
		fmt.Println("  |              v1.0 Rev A             |")
		fmt.Println("  |-------------------------------------|")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("   1. AC Power On")
		fmt.Println("   2. AC Power Off")
		fmt.Println("   -----------------")
		fmt.Println("   3. DC Power On")
		fmt.Println("   4. DC Power Off")
		fmt.Println("   -----------------")
		fmt.Println("   5. AC Power Cycle")
		fmt.Println("   6. DC Power Cycle")
		fmt.Println("   -----------------")
		fmt.Println("   C. Com Port Settings")
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
			com5.SendCommand("A0 01 01 A2")
			fmt.Println("AC Power On")
		case "2":
			clearScreen()
			com5.SendCommand("A0 01 00 A1")
			fmt.Println("AC Power Off")
		case "3":
			clearScreen()
			com6.SendCommand("A0 01 01 A2")
			fmt.Println("DC Power On")
		case "4":
			clearScreen()
			com6.SendCommand("A0 01 00 A1")
			fmt.Println("DC Power Off")
		case "5":
			clearScreen()
			fmt.Println("|-------------------------------------|")
			fmt.Println("|        Power Cycling AC Power       |")
			fmt.Println("|             Please wait...          |")
			fmt.Println("|-------------------------------------|")
			com5.SendCommand("A0 01 00 A1")
			fmt.Println("Cycling AC Power (Turned off)")
			time.Sleep(2 * time.Second)
			com5.SendCommand("A0 01 01 A2")
			fmt.Println("Cycling AC Power (Turned on)")
		case "6":
			clearScreen()
			fmt.Println("  |-------------------------------------|")
			fmt.Println("  |        Power Cycling DC Power       |")
			fmt.Println("  |             Please wait...          |")
			fmt.Println("  |-------------------------------------|")
			com6.SendCommand("A0 01 00 A1")
			fmt.Println("Cycling DC Power (Turned off)")
			time.Sleep(2 * time.Second)
			com6.SendCommand("A0 01 01 A2")
			fmt.Println("Cycling DC Power (Turned on)")
		case "e":
			fmt.Println("Exiting...")
		
			return
		case "c":
			clearScreen()
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  |              COM Port Settings               |")
			fmt.Println("  | Set the COM Ports to match in Device Manager |")
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  | AC Power Relay: COM5  | DC Power Relay: COM6 |")
			fmt.Println("  |----------------------------------------------|")
			time.Sleep(6 * time.Second)
		case "i":
			clearScreen()
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  |      Repeater System Power Control App       |")
			fmt.Println("  |                  v1.0 Rev A                  |")
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  |               Written By: KD9HAE             |")
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  |               Written In: GoLang             |")
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  |             Written: Sept: 23rd 2023         |")
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  |             Updated:                         |")
			fmt.Println("  |----------------------------------------------|")
			fmt.Println("  | https://github.com/jfeyen4249/repeater_power |")
			fmt.Println("  |----------------------------------------------|")
			time.Sleep(6 * time.Second)

		default:
			fmt.Println("Invalid option. Please select a valid option (1-6, 0 to exit).")
		}

		time.Sleep(1 * time.Second) // Wait for 3 seconds before reloading the menu
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
