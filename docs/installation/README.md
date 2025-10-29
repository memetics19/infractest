# Installation Guide

This guide covers different ways to install Infractest on your system.

## 📋 Prerequisites

Before installing Infractest, ensure you have the following prerequisites:

### Required
- **Go 1.25.3 or later** - [Download Go](https://golang.org/dl/)
- **Terraform** - [Install Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)

### Optional (for live mode testing)
- **Cloud Provider CLI tools** (AWS CLI, Azure CLI, etc.)
- **Cloud Provider credentials** configured

## 🚀 Installation Methods

### Method 1: Go Install (Recommended)

The easiest way to install Infractest is using `go install`:

```bash
go install github.com/memetics19/infractest/cmd/infractest@latest
```

This will install the latest version to your `$GOPATH/bin` directory.

### Method 2: Build from Source

1. **Clone the repository**:
   ```bash
   git clone https://github.com/memetics19/infractest.git
   cd infractest
   ```

2. **Build the binary**:
   ```bash
   go build -o bin/infractest cmd/infractest/main.go
   ```

3. **Make it executable**:
   ```bash
   chmod +x bin/infractest
   ```

4. **Add to PATH** (optional):
   ```bash
   export PATH=$PATH:$(pwd)/bin
   ```

### Method 3: Download Pre-built Binaries

Download pre-built binaries from [GitHub Releases](https://github.com/memetics19/infractest/releases):

```bash
# For Linux (AMD64)
wget https://github.com/memetics19/infractest/releases/latest/download/infractest-linux-amd64.tar.gz
tar -xzf infractest-linux-amd64.tar.gz
sudo mv infractest /usr/local/bin/

# For macOS (ARM64)
wget https://github.com/memetics19/infractest/releases/latest/download/infractest-darwin-arm64.tar.gz
tar -xzf infractest-darwin-arm64.tar.gz
sudo mv infractest /usr/local/bin/

# For Windows (AMD64)
# Download infractest-windows-amd64.zip and extract to your PATH
```

## ✅ Verify Installation

After installation, verify that Infractest is working correctly:

```bash
infractest --help
```

You should see output similar to:

```
Usage of infractest:
  -dir string
        directory containing .tfunittest.hcl test files (default "tests")
  -json string
        path to write JSON report (optional)
  -mode string
        test mode: mock | live (default "mock")
```

## 🔧 Configuration

### Environment Variables

Set these environment variables for optimal performance:

```bash
# Terraform configuration
export TF_LOG=INFO                    # Enable Terraform logging
export TF_LOG_PATH=terraform.log      # Log file location
export TF_IN_AUTOMATION=true          # Disable interactive prompts

# For live mode testing
export AWS_ACCESS_KEY_ID=your_key     # AWS credentials
export AWS_SECRET_ACCESS_KEY=your_secret
export AWS_DEFAULT_REGION=us-east-1   # Default AWS region

# For Azure live mode testing
export AZURE_CLIENT_ID=your_client_id
export AZURE_CLIENT_SECRET=your_secret
export AZURE_TENANT_ID=your_tenant_id
export AZURE_SUBSCRIPTION_ID=your_subscription_id
```

### Shell Configuration

Add Infractest to your shell configuration:

**Bash** (`~/.bashrc` or `~/.bash_profile`):
```bash
export PATH=$PATH:$HOME/go/bin
```

**Zsh** (`~/.zshrc`):
```bash
export PATH=$PATH:$HOME/go/bin
```

**Fish** (`~/.config/fish/config.fish`):
```fish
set -gx PATH $PATH $HOME/go/bin
```

## 🐳 Docker Installation

For containerized environments, you can use Docker:

### Using Dockerfile

Create a `Dockerfile`:

```dockerfile
FROM golang:1.25.3-alpine AS builder

# Install Terraform
RUN apk add --no-cache wget unzip
RUN wget https://releases.hashicorp.com/terraform/1.5.0/terraform_1.5.0_linux_amd64.zip
RUN unzip terraform_1.5.0_linux_amd64.zip
RUN mv terraform /usr/local/bin/

# Install Infractest
RUN go install github.com/memetics19/infractest/cmd/infractest@latest

FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy Terraform and Infractest
COPY --from=builder /usr/local/bin/terraform /usr/local/bin/
COPY --from=builder /root/go/bin/infractest /usr/local/bin/

WORKDIR /workspace
ENTRYPOINT ["infractest"]
```

Build and run:

```bash
# Build the image
docker build -t infractest .

# Run tests
docker run -v $(pwd):/workspace infractest -dir tests
```

### Using Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'
services:
  infractest:
    build: .
    volumes:
      - .:/workspace
    working_dir: /workspace
    environment:
      - TF_LOG=INFO
      - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
      - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
    command: ["-dir", "tests", "-mode", "live"]
```

Run with:

```bash
docker-compose run infractest
```

## 🔄 Updating Infractest

### Using Go Install

```bash
go install github.com/memetics19/infractest/cmd/infractest@latest
```

### From Source

```bash
cd infractest
git pull origin main
go build -o bin/infractest cmd/infractest/main.go
```

### Check Version

```bash
infractest --version
```

## 🛠 Development Installation

For contributing to Infractest:

1. **Fork and clone**:
   ```bash
   git clone https://github.com/your-username/infractest.git
   cd infractest
   git remote add upstream https://github.com/memetics19/infractest.git
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Install development tools**:
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   go install github.com/air-verse/air@latest
   ```

4. **Build and test**:
   ```bash
   go build -o bin/infractest cmd/infractest/main.go
   go test ./...
   ```

## 🚨 Troubleshooting

### Common Installation Issues

#### Go Not Found
```bash
# Error: go: command not found
# Solution: Install Go from https://golang.org/dl/
```

#### Permission Denied
```bash
# Error: permission denied
# Solution: Use sudo or install to user directory
go install github.com/memetics19/infractest/cmd/infractest@latest
# or
sudo go install github.com/memetics19/infractest/cmd/infractest@latest
```

#### Module Not Found
```bash
# Error: module not found
# Solution: Ensure Go modules are enabled
export GO111MODULE=on
go install github.com/memetics19/infractest/cmd/infractest@latest
```

#### Terraform Not Found
```bash
# Error: terraform: command not found
# Solution: Install Terraform
# See: https://learn.hashicorp.com/tutorials/terraform/install-cli
```

### Verification Steps

1. **Check Go version**:
   ```bash
   go version
   # Should be 1.25.3 or later
   ```

2. **Check Terraform version**:
   ```bash
   terraform version
   # Should be installed and working
   ```

3. **Check Infractest installation**:
   ```bash
   infractest --help
   # Should show help output
   ```

4. **Test with sample**:
   ```bash
   # Create a simple test
   mkdir -p test-sample
   echo 'test "sample" { module = "." }' > test-sample/sample.tfunittest.hcl
   
   # Run the test
   infractest -dir test-sample
   ```

## 📚 Next Steps

After successful installation:

1. Read the [Quick Start Guide](../usage/quick-start.md)
2. Learn about [Basic Concepts](../usage/concepts.md)
3. Try the [Basic Examples](../examples/basic.md)
4. Explore [Writing Tests](../usage/writing-tests.md)

## 🆘 Getting Help

If you encounter issues during installation:

1. Check the [Troubleshooting Guide](../../troubleshooting/common-issues.md)
2. Search [GitHub Issues](https://github.com/memetics19/infractest/issues)
3. Ask questions in [GitHub Discussions](https://github.com/memetics19/infractest/discussions)
4. Review the [FAQ](../../troubleshooting/faq.md)
