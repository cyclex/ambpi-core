# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build

# Main target
all: build

# Build the executable
build:
	$(GOBUILD) -o engine cmd/*.go

# Run the application
run:
	$(GOBUILD) -o engine cmd/*.go
	mv engine ../
	sudo systemctl restart ambpi-server
	sudo systemctl restart ambpi-webhook

# Default target to run the application
.DEFAULT_GOAL := run