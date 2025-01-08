# Adopt from https://pytorch.org/tutorials/intermediate/tensorboard_profiler_tutorial.html
import torch
import torch.nn
import torch.optim
import torch.profiler
import torch.utils.data
import torchvision.datasets
import torchvision.models
import torchvision.transforms as T
from torchvision.models import EfficientNet_B0_Weights, efficientnet_b0

transform = T.Compose(
    [T.Resize(224), T.ToTensor(), T.Normalize((0.5, 0.5, 0.5), (0.5, 0.5, 0.5))]
)
train_set = torchvision.datasets.CIFAR10(
    root="./data", train=True, download=True, transform=transform
)
train_loader = torch.utils.data.DataLoader(train_set, batch_size=2, shuffle=True)

device = torch.device("cuda:0")
model = efficientnet_b0(weights=EfficientNet_B0_Weights.DEFAULT).to(device)
criterion = torch.nn.CrossEntropyLoss().cuda(device)
optimizer = torch.optim.SGD(model.parameters(), lr=0.001, momentum=0.9)
model.train()


def train(data):
    inputs, labels = data[0].to(device=device), data[1].to(device=device)
    outputs = model(inputs)
    loss = criterion(outputs, labels)
    optimizer.zero_grad()
    loss.backward()
    optimizer.step()


with torch.profiler.profile(
    schedule=torch.profiler.schedule(wait=1, warmup=1, active=3, repeat=2),
    on_trace_ready=torch.profiler.tensorboard_trace_handler(
        "/home/envd/log/efficientnet"
    ),
    record_shapes=True,
    profile_memory=True,
    with_stack=True,
) as prof:
    for step, batch_data in enumerate(train_loader):
        if step >= (1 + 1 + 3) * 2:
            break
        train(batch_data)
        prof.step()  # Need to call this at the end of each step to notify profiler of steps' boundary.


prof = torch.profiler.profile(
    schedule=torch.profiler.schedule(wait=1, warmup=1, active=3, repeat=2),
    on_trace_ready=torch.profiler.tensorboard_trace_handler(
        "/home/envd/log/efficientnet"
    ),
    record_shapes=True,
    with_stack=True,
)
prof.start()
for step, batch_data in enumerate(train_loader):
    if step >= (1 + 1 + 3) * 2:
        break
    train(batch_data)
    prof.step()
prof.stop()
