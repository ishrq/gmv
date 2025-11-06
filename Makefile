.PHONY: all build man clean install test

# Variables
BINARY_NAME=gmv
INSTALL_PATH=/usr/local/bin
MAN_PATH=/usr/local/share/man/man1
MAN_SOURCE=gmv.1
MAN_GZ=gmv.1.gz

# Default target
all: build man

# Build the binary
build:
	go build -o $(BINARY_NAME) main.go

# Generate and compress man page
man:
	go run tools/gen-man/*.go > $(MAN_SOURCE)
	gzip -f $(MAN_SOURCE)
	@echo "Generated $(MAN_GZ)"

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f $(MAN_SOURCE) $(MAN_GZ)

# Install binary and man page to system
install: build man
	install -d $(INSTALL_PATH)
	install -m 755 $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	install -d $(MAN_PATH)
	install -m 644 $(MAN_GZ) $(MAN_PATH)/$(MAN_GZ)
	@echo "Installed $(BINARY_NAME) to $(INSTALL_PATH)"
	@echo "Installed man page to $(MAN_PATH)"

# Uninstall from system
uninstall:
	rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	rm -f $(MAN_PATH)/$(MAN_GZ)
	@echo "Uninstalled $(BINARY_NAME)"

# Run tests
test:
	cd test && go test -v

# Display help
help:
	@echo "Makefile targets:"
	@echo "  all       - Build binary and man page (default)"
	@echo "  build     - Build the gmv binary"
	@echo "  man       - Generate and compress man page"
	@echo "  clean     - Remove build artifacts"
	@echo "  install   - Install binary and man page to system"
	@echo "  uninstall - Remove binary and man page from system"
	@echo "  test      - Run integration tests"
	@echo "  help      - Display this help message"
