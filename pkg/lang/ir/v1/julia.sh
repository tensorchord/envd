set -euo pipefail && \
JULIA_URL="https://julialang-s3.julialang.org/bin/linux/x64/1.8/julia-1.8.5-linux-x86_64.tar.gz"; \
SHA256SUM="e71a24816e8fe9d5f4807664cbbb42738f5aa9fe05397d35c81d4c5d649b9d05"; \

wget "${JULIA_URL}" -O /tmp/julia.tar.gz && \
echo "${SHA256SUM}  /tmp/julia.tar.gz" > /tmp/sha256sum && \
sha256sum -c -s /tmp/sha256sum && \

res=$?

if [ $res=="OK" ]; then  
    echo "PASSED"
else
    echo "FAILED" && \
    rm /tmp/julia.tar.gz && \
    wget "${JULIA_URL}" -O /tmp/julia.tar.gz && \
    sha256sum -c -s /tmp/sha256sum
fi

