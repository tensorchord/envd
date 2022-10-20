import random
import sys
import os

from diffusers import StableDiffusionPipeline

device = "cuda"


def dummy(images, **kwargs):
    return images, False


# Read prompt from command line
prompt = " ".join(sys.argv[1:])

pipe = StableDiffusionPipeline.from_pretrained(
    "CompVis/stable-diffusion-v1-4", use_auth_token=os.environ["HUGGINGFACE_TOKEN"]
)
pipe.to(device)
pipe.safety_checker = dummy

# Run until we exit with CTRL+C
while True:
    n = random.randint(1000, 9999)
    image = pipe(prompt, guidance_scale=7.5)["sample"][0]
    image.save(f"{n}.jpeg")
