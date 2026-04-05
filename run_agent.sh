#!/bin/bash
MAX_RETRIES=5
DELAY=15
count=0

echo "Starting Claude Sub-agent..."
while [ $count -lt $MAX_RETRIES ]; do
    echo "[$(date)] Attempt $((count+1)) / $MAX_RETRIES"
    
    # 执行 claude 命令
    claude -p --dangerously-skip-permissions "$(cat claude_prompt.txt)"
    EXIT_CODE=$?
    
    if [ $EXIT_CODE -eq 0 ]; then
        echo "[$(date)] Agent finished successfully."
        exit 0
    fi
    
    echo "[$(date)] Agent failed with exit code $EXIT_CODE. Retrying in $DELAY seconds..."
    count=$((count+1))
    sleep $DELAY
done

echo "[$(date)] Agent failed after $MAX_RETRIES attempts."
exit 1
