# ResNet Serving Example

## Quick start (GPU)

Build and run your service:

```bash
envd build -t resnet:serving
docker run --rm --gpus all -p 8000:8000 resnet
```

Test with client:

```bash
python client.py
```

Check the metrics:

```bash
curl localhost:8000/metrics
```

