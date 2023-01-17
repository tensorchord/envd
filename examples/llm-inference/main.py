from transformers import AutoModelForCausalLM, AutoTokenizer, set_seed
from transformers import pipeline

name = "bigscience/bloom-3b"
text = "Hello my name is"
max_new_tokens = 20


def generate_from_model(model, tokenizer):
    encoded_input = tokenizer(text, return_tensors="pt")
    output_sequences = model.generate(input_ids=encoded_input["input_ids"].cuda())
    return tokenizer.decode(output_sequences[0], skip_special_tokens=True)


pipe = pipeline(
    model=name,
    model_kwargs={"device_map": "auto", "load_in_8bit": True},
    max_new_tokens=max_new_tokens,
)
print(pipe(text))
