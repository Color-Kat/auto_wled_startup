# Auto WLED Startup

An open-source Go project to control WLED devices, 
providing features like automatic startup.

## Features

- Turns on WLED at startup with "Blends" effect (75%) or solid color (25%).
- Menu to set IP, manage autostart, and run presets.
- Works on Windows, Linux, macOS.

## Installation

1. Go to [Releases](https://github.com/Color-Kat/auto_wled_startup/releases).
2. Download your binary:
    - Windows: `wled_startup_windows_amd64.exe` (64-bit)
    - Linux: `wled_startup_linux_amd64` (64-bit)
    - macOS: `wled_startup_darwin_amd64` (Intel) or `darwin_arm64` (M1/M2)
3. Put it in a folder (e.g., `C:\WLED` or `~/wled`).

## Usage

- **Menu**: Run the binary (e.g., `wled_startup_windows_amd64.exe` or `./wled_startup_linux_amd64`).
- **Auto-start test**: Add `--run` (e.g., `wled_startup_windows_amd64.exe --run`) to test script that runs at startup.

## Setup

1. Run the binary to open the menu.
2. Pick `1` to set your WLED IP (e.g., `192.168.1.100`).
3. It will create `wled_config.txt` in the folder.

## Autostart

- Menu option `2`: Adds to startup.
- Menu option `3`: Removes from startup.

## Troubleshooting

- **No effect**: Check IP in `wled_config.txt` and network.
- **Autostart fails**: Verify startup folder (Windows) or service (Linux/macOS).
