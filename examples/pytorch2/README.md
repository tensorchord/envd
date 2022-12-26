# PyTorch profiler example

This example is adopted from torch's [official tutorial](https://pytorch.org/tutorials/intermediate/tensorboard_profiler_tutorial.html). It shows how to use PyTorch profiler to analyze performance bottlenecks in a model.

## Usage

Run `envd up` at current folder. And execute `python main.py` in the envd container. Then you can see the profiling result through TensorBoard at http://localhost:8888

## torch.compile

mode="max-autotune" triggers a bug like https://github.com/pytorch/torchdynamo/issues/1593
```
    model = torch.compile(efficientnet_b0(weights=EfficientNet_B0_Weights.DEFAULT))
```
  
