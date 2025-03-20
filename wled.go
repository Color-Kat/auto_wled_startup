package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func setWledIP(scanner *bufio.Scanner) {
	fmt.Print("Enter WLED IP Address (e.g., 192.168.1.100): ")
	scanner.Scan()
	ip := scanner.Text()
	if ip == "" {
		fmt.Println("IP cannot be empty.")
		return
	}
	err := os.WriteFile(configFile, []byte(ip), 0644)
	if err != nil {
		fmt.Println("Error saving IP:", err)
		return
	}
	fmt.Println("WLED IP set to", ip)
}

func getWledIP() string {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func sendWledRequest(jsonData string) error {
	wledIP := getWledIP()
	if wledIP == "" {
		return fmt.Errorf("WLED IP not set. Please set it in the menu.")
	}
	url := fmt.Sprintf("http://%s/json/state", wledIP)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func runRandomGradient() {
	randVal := rand.Float32()
	var jsonData string

	if randVal < 0.25 {
		r, g, b := rand.Intn(256), rand.Intn(256), rand.Intn(256)
		jsonData = fmt.Sprintf(`{
            "on": true,
            "bri": 250,
            "seg": [{
                "col": [[%d, %d, %d]],
                "fx": 0,
                "sx": 0
            }]
        }`, r, g, b)
	} else {
		paletteID := rand.Intn(56) // Палитры от 0 до 55
		jsonData = fmt.Sprintf(`{
            "on": true,
            "bri": 250,
            "seg": [{
                "fx": 115,
                "sx": 128,
                "pal": %d
            }]
        }`, paletteID)
	}

	err := sendWledRequest(jsonData)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func runPreset(scanner *bufio.Scanner) {
	fmt.Print("Enter Preset ID: ")
	scanner.Scan()
	presetID, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid Preset ID. Must be a number.")
		return
	}

	jsonData := fmt.Sprintf(`{"on": true, "ps": %d}`, presetID)
	err = sendWledRequest(jsonData)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
