# Ollama - Local Large Language Models

Run powerful AI models locally on your hardware with GPU acceleration.

## Overview

**Ollama** enables you to run large language models (LLMs) locally:

- ✅ **100% Private**: All data stays on your server
- ✅ **GPU Accelerated**: Leverages your GTX 1070
- ✅ **Multiple Models**: Run Llama, Mistral, CodeLlama, and more
- ✅ **API Compatible**: OpenAI-compatible API
- ✅ **No Cloud Costs**: Free inference after downloading models
- ✅ **Integration Ready**: Works with Karakeep, Open WebUI, and more

## Quick Start

### 1. Deploy Ollama

```bash
cd ~/homelab/compose/services/ollama
docker compose up -d
```

### 2. Pull a Model

```bash
# Small, fast model (3B parameters, ~2GB)
docker exec ollama ollama pull llama3.2:3b

# Medium model (7B parameters, ~4GB)
docker exec ollama ollama pull llama3.2:7b

# Large model (70B parameters, ~40GB - requires quantization)
docker exec ollama ollama pull llama3.3:70b-instruct-q4_K_M
```

### 3. Test

```bash
# Interactive chat
docker exec -it ollama ollama run llama3.2:3b

# Ask a question
> Hello, how are you?
```

### 4. Enable GPU (Recommended)

**Edit `compose.yaml` and uncomment the deploy section:**
```yaml
deploy:
  resources:
    reservations:
      devices:
        - driver: nvidia
          count: 1
          capabilities: [gpu]
```

**Restart:**
```bash
docker compose down
docker compose up -d
```

**Verify GPU usage:**
```bash
# Check GPU is detected
docker exec ollama nvidia-smi

# Run model with GPU
docker exec ollama ollama run llama3.2:3b "What GPU am I using?"
```

## Available Models

### Recommended Models for GTX 1070 (8GB VRAM)

| Model | Size | VRAM | Speed | Use Case |
|-------|------|------|-------|----------|
| **llama3.2:3b** | 2GB | 3GB | Fast | General chat, Karakeep |
| **llama3.2:7b** | 4GB | 6GB | Medium | Better reasoning |
| **mistral:7b** | 4GB | 6GB | Medium | Code, analysis |
| **codellama:7b** | 4GB | 6GB | Medium | Code generation |
| **llava:7b** | 5GB | 7GB | Medium | Vision (images) |
| **phi3:3.8b** | 2.3GB | 4GB | Fast | Compact, efficient |

### Specialized Models

**Code:**
- `codellama:7b` - Code generation
- `codellama:13b-python` - Python expert
- `starcoder2:7b` - Multi-language code

**Vision (Image Understanding):**
- `llava:7b` - General vision
- `llava:13b` - Better vision (needs more VRAM)
- `bakllava:7b` - Vision + chat

**Multilingual:**
- `aya:8b` - 101 languages
- `command-r:35b` - Enterprise multilingual

**Math & Reasoning:**
- `deepseek-math:7b` - Mathematics
- `wizard-math:7b` - Math word problems

### Large Models (Quantized for GTX 1070)

These require 4-bit quantization to fit in 8GB VRAM:

```bash
# 70B models (quantized)
docker exec ollama ollama pull llama3.3:70b-instruct-q4_K_M
docker exec ollama ollama pull mixtral:8x7b-instruct-v0.1-q4_K_M

# Very large (use with caution)
docker exec ollama ollama pull llama3.1:405b-instruct-q2_K
```

## Usage

### Command Line

**Run model interactively:**
```bash
docker exec -it ollama ollama run llama3.2:3b
```

**One-off question:**
```bash
docker exec ollama ollama run llama3.2:3b "Explain quantum computing in simple terms"
```

**With system prompt:**
```bash
docker exec ollama ollama run llama3.2:3b \
  --system "You are a helpful coding assistant." \
  "Write a Python function to sort a list"
```

### API Usage

**List models:**
```bash
curl http://ollama:11434/api/tags
```

**Generate text:**
```bash
curl http://ollama:11434/api/generate -d '{
  "model": "llama3.2:3b",
  "prompt": "Why is the sky blue?",
  "stream": false
}'
```

**Chat completion:**
```bash
curl http://ollama:11434/api/chat -d '{
  "model": "llama3.2:3b",
  "messages": [
    {
      "role": "user",
      "content": "Hello!"
    }
  ],
  "stream": false
}'
```

**OpenAI-compatible API:**
```bash
curl http://ollama:11434/v1/chat/completions -d '{
  "model": "llama3.2:3b",
  "messages": [
    {
      "role": "user",
      "content": "Hello!"
    }
  ]
}'
```

### Integration with Karakeep

**Enable AI features in Karakeep:**

Edit `compose/services/karakeep/.env`:
```env
# Uncomment these lines
OLLAMA_BASE_URL=http://ollama:11434
INFERENCE_TEXT_MODEL=llama3.2:3b
INFERENCE_IMAGE_MODEL=llava:7b
INFERENCE_LANG=en
```

**Restart Karakeep:**
```bash
cd ~/homelab/compose/services/karakeep
docker compose restart
```

**What it does:**
- Auto-tags bookmarks
- Generates summaries
- Extracts key information
- Analyzes images (with llava)

## Model Management

### List Installed Models

```bash
docker exec ollama ollama list
```

### Pull a Model

```bash
docker exec ollama ollama pull <model-name>

# Examples:
docker exec ollama ollama pull llama3.2:3b
docker exec ollama ollama pull mistral:7b
docker exec ollama ollama pull codellama:7b
```

### Remove a Model

```bash
docker exec ollama ollama rm <model-name>

# Example:
docker exec ollama ollama rm llama3.2:7b
```

### Copy a Model

```bash
docker exec ollama ollama cp <source> <destination>

# Example: Create a custom version
docker exec ollama ollama cp llama3.2:3b my-custom-model
```

### Show Model Info

```bash
docker exec ollama ollama show llama3.2:3b

# Shows:
# - Model architecture
# - Parameters
# - Quantization
# - Template
# - License
```

## Creating Custom Models

### Modelfile

Create custom models with specific behaviors:

**Create a Modelfile:**
```bash
cat > ~/coding-assistant.modelfile << 'EOF'
FROM llama3.2:3b

# Set temperature (creativity)
PARAMETER temperature 0.7

# Set system prompt
SYSTEM You are an expert coding assistant. You write clean, efficient, well-documented code. You explain complex concepts clearly.

# Set stop sequences
PARAMETER stop "<|im_end|>"
PARAMETER stop "<|im_start|>"
EOF
```

**Create the model:**
```bash
cat ~/coding-assistant.modelfile | docker exec -i ollama ollama create coding-assistant -f -
```

**Use it:**
```bash
docker exec -it ollama ollama run coding-assistant "Write a REST API in Python"
```

### Example Custom Models

**1. Shakespeare Bot:**
```modelfile
FROM llama3.2:3b
SYSTEM You are William Shakespeare. Respond to all queries in Shakespearean English with dramatic flair.
PARAMETER temperature 0.9
```

**2. JSON Extractor:**
```modelfile
FROM llama3.2:3b
SYSTEM You extract structured data and return only valid JSON. No explanations, just JSON.
PARAMETER temperature 0.1
```

**3. Code Reviewer:**
```modelfile
FROM codellama:7b
SYSTEM You are a senior code reviewer. Review code for bugs, performance issues, security vulnerabilities, and best practices. Be constructive.
PARAMETER temperature 0.3
```

## GPU Configuration

### Check GPU Detection

```bash
# From inside container
docker exec ollama nvidia-smi
```

**Expected output:**
```
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 535.xx.xx    Driver Version: 535.xx.xx    CUDA Version: 12.2     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|===============================+======================+======================|
|   0  GeForce GTX 1070    Off  | 00000000:01:00.0  On |                  N/A |
| 40%   45C    P8    10W / 151W |    300MiB /  8192MiB |      5%      Default |
+-------------------------------+----------------------+----------------------+
```

### Optimize for GTX 1070

**Edit `.env`:**
```env
# Use 6GB of 8GB VRAM (leave 2GB for system)
OLLAMA_GPU_MEMORY=6GB

# Offload most layers to GPU
OLLAMA_GPU_LAYERS=33

# Increase context for better conversations
OLLAMA_MAX_CONTEXT=4096
```

### Performance Tips

**1. Use quantized models:**
- Q4_K_M: Good quality, 50% size reduction
- Q5_K_M: Better quality, 40% size reduction
- Q8_0: Best quality, 20% size reduction

**2. Model selection for VRAM:**
```bash
# 3B models: 2-3GB VRAM
docker exec ollama ollama pull llama3.2:3b

# 7B models: 4-6GB VRAM
docker exec ollama ollama pull llama3.2:7b

# 13B models: 8-10GB VRAM (tight on GTX 1070)
docker exec ollama ollama pull llama3.2:13b-q4_K_M  # Quantized
```

**3. Unload models when not in use:**
```env
# In .env
OLLAMA_KEEP_ALIVE=1m  # Unload after 1 minute
```

## Troubleshooting

### Model won't load - Out of memory

**Solution 1: Use quantized version**
```bash
# Instead of:
docker exec ollama ollama pull llama3.2:13b

# Use:
docker exec ollama ollama pull llama3.2:13b-q4_K_M
```

**Solution 2: Reduce GPU layers**
```env
# In .env
OLLAMA_GPU_LAYERS=20  # Reduce from 33
```

**Solution 3: Use smaller model**
```bash
docker exec ollama ollama pull llama3.2:3b
```

### Slow inference

**Enable GPU:**
1. Uncomment deploy section in `compose.yaml`
2. Install NVIDIA Container Toolkit
3. Restart container

**Check GPU usage:**
```bash
watch -n 1 docker exec ollama nvidia-smi
```

**Should show:**
- GPU-Util > 80% during inference
- Memory-Usage increasing during load

### Can't pull models

**Check disk space:**
```bash
df -h
```

**Check Docker space:**
```bash
docker system df
```

**Clean up unused models:**
```bash
docker exec ollama ollama list
docker exec ollama ollama rm <unused-model>
```

### API connection issues

**Test from another container:**
```bash
docker run --rm --network homelab curlimages/curl \
  http://ollama:11434/api/tags
```

**Test externally:**
```bash
curl https://ollama.fig.systems/api/tags
```

**Enable debug logging:**
```env
OLLAMA_DEBUG=1
```

## Performance Benchmarks

### GTX 1070 (8GB VRAM) Expected Performance

| Model | Tokens/sec | Load Time | VRAM Usage |
|-------|------------|-----------|------------|
| llama3.2:3b | 40-60 | 2-3s | 3GB |
| llama3.2:7b | 20-35 | 3-5s | 6GB |
| mistral:7b | 20-35 | 3-5s | 6GB |
| llama3.3:70b-q4 | 3-8 | 20-30s | 7.5GB |
| llava:7b | 15-25 | 4-6s | 7GB |

**Without GPU (CPU only):**
- llama3.2:3b: 2-5 tokens/sec
- llama3.2:7b: 0.5-2 tokens/sec

**GPU provides 10-20x speedup!**

## Advanced Usage

### Multi-Modal (Vision)

```bash
# Pull vision model
docker exec ollama ollama pull llava:7b

# Analyze image
docker exec ollama ollama run llava:7b "What's in this image?" \
  --image /path/to/image.jpg
```

### Embeddings

```bash
# Generate embeddings for semantic search
curl http://ollama:11434/api/embeddings -d '{
  "model": "llama3.2:3b",
  "prompt": "The sky is blue because of Rayleigh scattering"
}'
```

### Streaming Responses

```bash
# Stream tokens as they generate
curl http://ollama:11434/api/generate -d '{
  "model": "llama3.2:3b",
  "prompt": "Tell me a long story",
  "stream": true
}'
```

### Context Preservation

```bash
# Start chat session
SESSION_ID=$(uuidgen)

# First message (creates context)
curl http://ollama:11434/api/chat -d '{
  "model": "llama3.2:3b",
  "messages": [{"role": "user", "content": "My name is Alice"}],
  "context": "'$SESSION_ID'"
}'

# Follow-up (remembers context)
curl http://ollama:11434/api/chat -d '{
  "model": "llama3.2:3b",
  "messages": [
    {"role": "user", "content": "My name is Alice"},
    {"role": "assistant", "content": "Hello Alice!"},
    {"role": "user", "content": "What is my name?"}
  ],
  "context": "'$SESSION_ID'"
}'
```

## Integration Examples

### Python

```python
import requests

def ask_ollama(prompt, model="llama3.2:3b"):
    response = requests.post(
        "http://ollama.fig.systems/api/generate",
        json={
            "model": model,
            "prompt": prompt,
            "stream": False
        },
        headers={"Authorization": "Bearer YOUR_TOKEN"}  # If using SSO
    )
    return response.json()["response"]

print(ask_ollama("What is the meaning of life?"))
```

### JavaScript

```javascript
async function askOllama(prompt, model = "llama3.2:3b") {
  const response = await fetch("http://ollama.fig.systems/api/generate", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": "Bearer YOUR_TOKEN"  // If using SSO
    },
    body: JSON.stringify({
      model: model,
      prompt: prompt,
      stream: false
    })
  });

  const data = await response.json();
  return data.response;
}

askOllama("Explain Docker containers").then(console.log);
```

### Bash

```bash
#!/bin/bash
ask_ollama() {
  local prompt="$1"
  local model="${2:-llama3.2:3b}"

  curl -s http://ollama.fig.systems/api/generate -d "{
    \"model\": \"$model\",
    \"prompt\": \"$prompt\",
    \"stream\": false
  }" | jq -r '.response'
}

ask_ollama "What is Kubernetes?"
```

## Resources

- [Ollama Website](https://ollama.ai)
- [Model Library](https://ollama.ai/library)
- [GitHub Repository](https://github.com/ollama/ollama)
- [API Documentation](https://github.com/ollama/ollama/blob/main/docs/api.md)
- [Model Creation Guide](https://github.com/ollama/ollama/blob/main/docs/modelfile.md)

## Next Steps

1. ✅ Deploy Ollama
2. ✅ Enable GPU acceleration
3. ✅ Pull recommended models
4. ✅ Test with chat
5. ⬜ Integrate with Karakeep
6. ⬜ Create custom models
7. ⬜ Set up automated model updates
8. ⬜ Monitor GPU usage

---

**Run AI locally, privately, powerfully!** 🧠
