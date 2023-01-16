#!/bin/bash

# URL of the file to download
url="https://julialang-s3.julialang.org/bin/linux/x64/1.8/julia-1.8.3-linux-x86_64.tar.gz"

# Download the file
curl -o julia.tar.gz $url

# Get the SHA256 hash of the downloaded file
hash=$(sha256sum julia.tar.gz | awk '{print $1}')

# Check the hash against the expected value
expected_hash="33c3b09356ffaa25d3331c3646b1f2d4b09944e8f93fcb994957801b8bbf58a9"
while [ "$hash" != "$expected_hash" ]; do
    rm julia.tar.gz
    curl -o julia.tar.gz $url
    hash=$(sha256sum julia.tar.gz | awk '{print $1}')
done