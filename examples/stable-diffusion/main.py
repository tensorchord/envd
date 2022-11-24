import random
import sys
import os

from diffusers import StableDiffusionPipeline
import torch
from torch import autocast

device = "cuda"


def dummy(images, **kwargs):
    return images, False


# Read prompt from command line
prompt = " ".join(sys.argv[1:])

# if your are limited by GPU memory and have less than 10GB of GPU RAM available, you can use fp16 for this example just like this line below
# pipe = StableDiffusionPipeline.from_pretrained("CompVis/stable-diffusion-v1-4", torch_dtype=torch.float16, revision="fp16",  use_auth_token=os.environ['HUGGINGFACE_TOKEN'])
pipe = StableDiffusionPipeline.from_pretrained(
    "CompVis/stable-diffusion-v1-4", use_auth_token=os.environ["HUGGINGFACE_TOKEN"]
)
pipe.to(device)
pipe.safety_checker = dummy

# Run until we exit with CTRL+C
while True:
    n = random.randint(1000, 9999)
    with autocast("cuda"):
        image = pipe(prompt, guidance_scale=7.5).images[0]
    image.save(f"{n}.jpeg")
