#!/bin/bash
# Build script for multiple platforms

set -e

APP_NAME="pupervisor"
BUILD_DIR="build/bin"
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

echo "==> Building ${APP_NAME} v${VERSION}"

mkdir -p ${BUILD_DIR}

# Build for current platform
build_current() {
    echo "==> Building for current platform..."
    CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o ${BUILD_DIR}/${APP_NAME} ./cmd/server
    echo "    Built: ${BUILD_DIR}/${APP_NAME}"
}

# Build for all platforms
build_all() {
    platforms=(
        "linux/amd64"
        "linux/arm64"
        "darwin/amd64"
        "darwin/arm64"
        "windows/amd64"
    )

    for platform in "${platforms[@]}"; do
        GOOS=${platform%/*}
        GOARCH=${platform#*/}
        output="${BUILD_DIR}/${APP_NAME}-${GOOS}-${GOARCH}"

        if [ "$GOOS" = "windows" ]; then
            output="${output}.exe"
        fi

        echo "==> Building for ${GOOS}/${GOARCH}..."
        CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="${LDFLAGS}" -o ${output} ./cmd/server
        echo "    Built: ${output}"
    done
}

case "${1:-current}" in
    all)
        build_all
        ;;
    *)
        build_current
        ;;
esac

echo ""
echo "==> Build complete!"
ls -la ${BUILD_DIR}/
