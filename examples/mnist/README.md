# MNIST Example

## Quick start (CPU)

```bash
$ envd up
```

By default, `envd up` will use the `build.envd` 

Then you can open jupyter notebook at [`http://localhost:8888`](http://localhost:8888), or open vscode remote to attach to the container.


## Quick start (GPU)

```bash
$ envd up -f build_gpu.envd
```

Also you can use `-f` option to specify the file to build

Then you can open jupyter notebook at [`http://localhost:8888`](http://localhost:8888), or open vscode remote to attach to the container.