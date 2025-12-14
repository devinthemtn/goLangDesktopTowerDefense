# Tower Defense Game Makefile

.PHONY: build run clean install-deps demo test-menu debug-waves help

# Default target
all: build

# Build the game
build:
	@echo "Building Tower Defense Game..."
	go build -o tower-defense *.go
	@echo "✅ Build complete! Run with: make run"

# Run the game
run: build
	@echo "Starting Tower Defense Game..."
	./tower-defense

# Install dependencies
install-deps:
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo "Installing system dependencies (Ubuntu/Debian)..."
	sudo apt-get update
	sudo apt-get install -y libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev

# Run graphics demo
demo: build
	@echo "Starting enhanced graphics demo..."
	./demo.sh

# Test menu navigation
test-menu: build
	@echo "Testing menu navigation system..."
	@echo "Instructions:"
	@echo "1. Use ↑/↓ keys to navigate between all 3 options"
	@echo "2. Try mouse hover and click on menu options"
	@echo "3. Verify 'Endless Mode' is selectable"
	@echo "4. Press ESC to exit when done testing"
	@echo ""
	@echo "Starting game for menu testing..."
	./tower-defense

# Debug wave progression
debug-waves: build
	@echo "Starting wave progression debug session..."
	./debug-waves.sh

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f tower-defense
	go clean

# Test the code
test:
	@echo "Running tests..."
	go test ./...

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Check for common issues
vet:
	@echo "Running go vet..."
	go vet ./...

# Create a release build
release: clean fmt vet
	@echo "Building release version..."
	go build -ldflags="-s -w" -o tower-defense *.go
	@echo "✅ Release build complete!"

# Show help
help:
	@echo "Tower Defense Game - Available Commands:"
	@echo ""
	@echo "  make build        - Build the game"
	@echo "  make run          - Build and run the game"
	@echo "  make demo         - Run enhanced graphics demonstration"
	@echo "  make test-menu    - Test menu navigation system"
	@echo "  make debug-waves  - Debug wave progression issues"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make install-deps - Install system and Go dependencies"
	@echo "  make test         - Run tests"
	@echo "  make fmt          - Format Go code"
	@echo "  make vet          - Run go vet"
	@echo "  make release      - Create optimized release build"
	@echo "  make help         - Show this help message"
	@echo ""
	@echo "Game Controls:"
	@echo "  Mouse Click       - Place tower"
	@echo "  Key 1            - Select Basic Tower ($50)"
	@echo "  Key 2            - Select Heavy Tower ($100)"
