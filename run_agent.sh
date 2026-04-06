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
        echo "[$(date)] Agent finished successfully. Attempting to commit and merge..."
        # 1. 兜底提交（防止 Claude 忘了主动跑 Commit 命令）
        git add .
        git commit -m "fix(automation): resolve parity bugs for batch 3" || true

        # 2. 回到主分支执行安全的合并。遇到冲突则让 AI 进行 Rebase + Self-Heal 自动修复
        while true; do
            cd "/Users/11176728/gemini-cli/dev-turms-im-refactor/turms-go"

            # 使用 flock 排队尝试 Merge
            (
                flock -x 200
                echo "[$(date)] Attempting merge for feature/fix-batch-3 into main..."
                if git merge "feature/fix-batch-3" --no-edit -m "Merge auto-fix batch 3 into main"; then
                    echo "SUCCESS" > .git/merge_result_batch_3
                else
                    echo "[!] Conflict detected. Aborting merge."
                    git merge --abort
                    echo "CONFLICT" > .git/merge_result_batch_3
                fi
            ) 200>.git/merge_lock.lock

            MERGE_RESULT=$(cat .git/merge_result_batch_3)
            rm -f .git/merge_result_batch_3

            if [ "$MERGE_RESULT" = "SUCCESS" ]; then
                echo "[$(date)] Successfully merged batch 3 into main."
                break
            fi

            echo "[$(date)] Merge conflict! Initiating Sub-Agent Self-Heal Rebase..."
            cd "/Users/11176728/gemini-cli/dev-turms-im-refactor/turms-worker-batch-3"

            # 开始 Rebase main
            git rebase main || {
                # 触发大模型自动修复冲突
                echo "[$(date)] Handing over to Claude for conflict resolution..."
                claude -p --dangerously-skip-permissions "A git rebase conflict was detected in feature/fix-batch-3.
Please resolve it. Here is the current status:
$(git status)

And the unmerged diff:
$(git diff --diff-filter=U)

You MUST:
1. Search and open the conflicted files to understand the collision.
2. Fix all conflict markers (<<<<<<< ======= >>>>>>>), keeping our bug fixes while incorporating main's latest code.
3. Run 'git add .' to stage the resolved files.
4. Run 'git rebase --continue' to finalize the conflict resolution.
Do NOT attempt to run standard git merge. Finish the rebase process. Keep your thoughts and logs concise."

                # 检查大模型是否成功继续了 rebase
                if git rebase --show-current-patch >/dev/null 2>&1 || [ -d "/Users/11176728/gemini-cli/dev-turms-im-refactor/turms-worker-batch-3/.git/rebase-merge" ]; then
                    echo "[!] Claude failed to finish the rebase. Aborting pipeline."
                    git rebase --abort
                    exit 1
                fi
            }

            echo "[$(date)] Self-Heal Rebase complete. Loop will retry the merge."
        done

        echo "[$(date)] Pipeline for batch 3 complete."
        exit 0
    fi

    echo "[$(date)] Agent failed with exit code $EXIT_CODE. Retrying in $DELAY seconds..."
    count=$((count+1))
    sleep $DELAY
done

echo "[$(date)] Agent failed after $MAX_RETRIES attempts."
exit 1
