package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const configFile = "wled_config.txt"

func main() {
	rand.Seed(time.Now().UnixNano())

	// Проверяем аргументы командной строки
	if len(os.Args) > 1 && os.Args[1] == "--run" {
		runRandomGradient()
		return
	}

	// Если аргументов нет, показываем меню
	for {
		showMenu()
	}
}

func showMenu() {
	fmt.Println("\n=== WLED Controller ===")
	fmt.Println("1. Set WLED IP Address")
	fmt.Println("2. Add to Startup (or Update)")
	fmt.Println("3. Remove from Startup")
	fmt.Println("4. Run Random Gradient with Blends")
	fmt.Println("5. Run Custom Preset")
	fmt.Println("6. Exit")
	fmt.Print("Choose an option: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := scanner.Text()

	switch choice {
	case "1":
		setWledIP(scanner)
	case "2":
		addToStartup()
	case "3":
		removeFromStartup()
	case "4":
		runRandomGradient()
	case "5":
		runPreset(scanner)
	case "6":
		os.Exit(0)
	default:
		fmt.Println("Invalid option, try again.")
	}
}
