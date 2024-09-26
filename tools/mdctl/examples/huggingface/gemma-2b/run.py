from transformers import AutoTokenizer, AutoModelForCausalLM

tokenizer = AutoTokenizer.from_pretrained("gemma-2b:latest")
model = AutoModelForCausalLM.from_pretrained("gemma-2b:latest")

input_text = "Who are you?"
input_ids = tokenizer(input_text, return_tensors="pt")

outputs = model.generate(**input_ids, max_length=64)
print(tokenizer.decode(outputs[0]))
