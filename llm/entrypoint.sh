#!/bin/sh

ollama serve &

# Wait until Ollama API is available
until curl -s http://localhost:11434/api/status >/dev/null; do
    echo "Waiting for Ollama to start..."
    sleep 2
done

# Pull the gemma3:1b model
curl -X POST -H "Content-Type: application/json" -d '{"name":"gemma3:1b"}' http://localhost:11434/api/pull

wait
