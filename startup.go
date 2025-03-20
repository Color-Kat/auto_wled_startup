package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func addToStartup() {
	wledIP := getWledIP()
	if wledIP == "" {
		fmt.Println("WLED IP not set. Please set it first.")
		return
	}

	// Get path to current exe file
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error finding executable:", err)
		return
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		fmt.Println("Error resolving absolute path:", err)
		return
	}

	// Copy to auto_wled_startup.exe
	startupExe := filepath.Join(filepath.Dir(exePath), "auto_wled_startup.exe")
	if err := copyFile(exePath, startupExe); err != nil {
		fmt.Println("Error copying executable:", err)
		return
	}

	// Add to startup with flag --run (skip menu, turn on wled)
	switch runtime.GOOS {
	case "windows":
		startupDir := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
		shortcutPath := filepath.Join(startupDir, "WLEDStartup.lnk")
		startupExeEscaped := strings.ReplaceAll(startupExe, `\`, `\\`)

		cmd := exec.Command(
			"powershell", "-Command",
			fmt.Sprintf(
				`$WshShell = New-Object -ComObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut('%s'); $Shortcut.TargetPath = '%s'; $Shortcut.Arguments = '--run'; $Shortcut.WorkingDirectory = '%s'; $Shortcut.Save()`,
				shortcutPath,
				startupExeEscaped,
				strings.ReplaceAll(filepath.Dir(startupExe), `\`, `\\`),
			),
		)

		if err := cmd.Run(); err != nil {
			fmt.Println("Error adding to startup:", err)
			return
		}

		fmt.Println("Added to Windows startup.")

	case "linux":
		serviceFile := "/etc/systemd/system/wled-startup.service"
		serviceContent := fmt.Sprintf(`[Unit]
Description=WLED Startup
After=network-online.target

[Service]
ExecStart=%s --run
WorkingDirectory=%s
Type=oneshot
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target`, startupExe, filepath.Dir(startupExe))

		if err := os.WriteFile(serviceFile, []byte(serviceContent), 0644); err != nil {
			fmt.Println("Error creating service file:", err)
			return
		}

		exec.Command("systemctl", "enable", "wled-startup.service").Run()
		fmt.Println("Added to Linux startup.")

	case "darwin":
		plistFile := filepath.Join(os.Getenv("HOME"), "Library/LaunchAgents/com.user.wledstartup.plist")
		plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.user.wledstartup</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
        <string>--run</string>
    </array>
    <key>WorkingDirectory</key>
    <string>%s</string>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>`, startupExe, filepath.Dir(startupExe))

		if err := os.WriteFile(plistFile, []byte(plistContent), 0644); err != nil {
			fmt.Println("Error creating plist file:", err)
			return
		}

		exec.Command("launchctl", "load", plistFile).Run()
		fmt.Println("Added to macOS startup.")
	}
}

func removeFromStartup() {
	exePath, _ := os.Executable()
	startupExe := filepath.Join(filepath.Dir(exePath), "auto_wled_startup.exe")

	switch runtime.GOOS {
	case "windows":
		startupDir := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
		shortcutPath := filepath.Join(startupDir, "WLEDStartup.lnk")

		if err := os.Remove(shortcutPath); err != nil {
			fmt.Println("Error removing from startup:", err)
		}

		if err := os.Remove(startupExe); err != nil {
			fmt.Println("Error removing startup executable:", err)
		}

		fmt.Println("Removed from Windows startup.")

	case "linux":
		serviceFile := "/etc/systemd/system/wled-startup.service"
		exec.Command("systemctl", "disable", "wled-startup.service").Run()

		if err := os.Remove(serviceFile); err != nil {
			fmt.Println("Error removing service file:", err)
		}

		fmt.Println("Removed from Linux startup.")

	case "darwin":
		plistFile := filepath.Join(os.Getenv("HOME"), "Library/LaunchAgents/com.user.wledstartup.plist")
		exec.Command("launchctl", "unload", plistFile).Run()

		if err := os.Remove(plistFile); err != nil {
			fmt.Println("Error removing plist file:", err)
		}

		fmt.Println("Removed from macOS startup.")
	}
}
