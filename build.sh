#!/bin/bash
targets=(
    "windows/386"
    "windows/amd64"
    "linux/386"
    "linux/amd64"
    "darwin/amd64"
    "darwin/arm64"
)
output_dir="bin"
mkdir -p $output_dir
for target in "${targets[@]}"; do
    os=$(echo $target | cut -d '/' -f1)
    arch=$(echo $target | cut -d '/' -f2)
    output_name="wled_startup_${os}_${arch}"
    if [ "$os" = "windows" ]; then
        output_name+=".exe"
    fi
    echo "Building for $os/$arch: $output_name"
    GOOS=$os GOARCH=$arch go build -o $output_dir/$output_name
done