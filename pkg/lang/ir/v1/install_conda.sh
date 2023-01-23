set -euo pipefail && \
sh /tmp/miniconda.sh -b -u -p /opt/conda && \
touch ~/.bashrc && \
echo ". /opt/conda/etc/profile.d/conda.sh" >> ~/.bashrc && \
echo "conda activate base" >> ~/.bashrc && \
echo -e "channels:\n  - defaults" > /opt/conda/.condarc && \
find /opt/conda/ -follow -type f -name '*.a' -delete && \
find /opt/conda/ -follow -type f -name '*.js.map' -delete && \
/opt/conda/bin/conda clean -afy
