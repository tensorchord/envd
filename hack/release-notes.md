# Quick Start

## What's New?

Find out on [changelog](https://github.com/tensorchord/envd/blob/main/CHANGELOG.md).

## Breaking Changes and Known Issues

Can be found in the [Getting Started](https://envd.tensorchord.ai/guide/getting-started.html).

## Installation

### CLI

#### Mac

Available via `curl`

```bash
# Download the binary
curl -sLO https://github.com/tensorchord/envd/releases/download/${version}/envd_${version}_Darwin_x86_64.gz

# Unzip
gunzip envd_${version}_Darwin_x86_64.gz

# Make binary executable
chmod +x envd_${version}_Darwin_x86_64.gz

# Move binary to path
mv ./envd_${version}_Darwin_x86_64 /usr/local/bin/envd

# Test installation
envd version
```

#### Linux

Available via `curl`

```bash
# Download the binary
curl -sLO https://github.com/tensorchord/envd/releases/download/${version}/envd_${version}_Linux_x86_64.gz

# Unzip
gunzip envd_${version}_Linux_x86_64.gz

# Make binary executable
chmod +x envd_${version}_Linux_x86_64

# Move binary to path
mv ./envd_${version}_Linux_x86_64 /usr/local/bin/envd

# Test installation
envd version
```