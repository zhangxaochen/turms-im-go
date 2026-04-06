const fs = require('fs');
const path = require('path');
const { execSync, spawn } = require('child_process');

(async function main() {
    const args = process.argv.slice(2);
    if (args.length < 2) {
        console.error("Usage: node dispatch_claude_workers.js <startBatchIndex> <endBatchIndex>");
        console.error("Example: node dispatch_claude_workers.js 0 2 (to run batch 0, 1, 2)");
        console.error("Note: Max allowed batches is 15 in total.");
        process.exit(1);
    }

    const startBatchIndex = parseInt(args[0], 10);
    const endBatchIndex = parseInt(args[1], 10);
    const BATCH_SIZE = 50;

    if (endBatchIndex - startBatchIndex + 1 > 15) {
        console.error("Error: Cannot start more than 15 sub-agents at once to prevent system overload.");
        process.exit(1);
    }

    const projectRoot = __dirname;
    const markdownUrl = path.join(projectRoot, 'docs', 'pending_bugs.md');
    if (!fs.existsSync(markdownUrl)) {
        console.error(`Error: Could not find ${markdownUrl}`);
        process.exit(1);
    }

    const content = fs.readFileSync(markdownUrl, 'utf-8');

    // 匹配起始锚点
    const anchor = "- [ ]";
    const anchorIndex = content.indexOf(anchor);

    if (anchorIndex === -1) {
        console.error("Error: Could not find the starting anchor in pending_bugs.md");
        process.exit(1);
    }

    const relevantContent = content.substring(anchorIndex);
    const lines = relevantContent.split('\n');

    let bugList = [];
    let currentBug = [];
    let currentHeader = '';

    for (const line of lines) {
        if (line.trim().startsWith('#')) {
            currentHeader = line.trim(); // 记录最近的标题上下文
        } else if (line.trim().startsWith('- [ ]')) {
            if (currentBug.length > 0) {
                bugList.push(currentBug.join('\n'));
            }
            currentBug = [currentHeader ? `[Context: ${currentHeader}]\n${line}` : line];
        } else if (currentBug.length > 0 && line.trim() !== '') {
            currentBug.push(line);
        }
    }
    if (currentBug.length > 0) {
        bugList.push(currentBug.join('\n'));
    }

    console.log(`Found a total of ${bugList.length} unresolved bugs after the anchor.`);

    for (let i = startBatchIndex; i <= endBatchIndex; i++) {
        const startIdx = i * BATCH_SIZE;
        const endIdx = startIdx + BATCH_SIZE;
        const batchBugs = bugList.slice(startIdx, endIdx);

        if (batchBugs.length === 0) {
            console.log(`Batch ${i} has no bugs. Stopping.`);
            break;
        }

        console.log(`\n=== Starting Worker for Batch ${i} (Bugs ${startIdx} to ${startIdx + batchBugs.length - 1}) ===`);

        const worktreeName = `turms-worker-batch-${i}`;
        const worktreePath = path.resolve(projectRoot, '..', worktreeName);
        const branchName = `feature/fix-batch-${i}`;

        // 1. Create worktree
        if (!fs.existsSync(worktreePath)) {
            console.log(`Creating git worktree at ${worktreePath}...`);
            try {
                execSync(`git worktree add ${worktreePath} -b ${branchName}`, { stdio: 'inherit', cwd: projectRoot });
            } catch (e) {
                // 如果分支已存在等异常则尝试直接检出
                execSync(`git worktree add ${worktreePath} ${branchName}`, { stdio: 'inherit', cwd: projectRoot });
            }
        } else {
            console.log(`Worktree ${worktreePath} already exists, reusing...`);
        }

        // 2. 执行 equivalent of 'wwwww'
        const envPath = path.join(projectRoot, '.env');
        if (fs.existsSync(envPath)) {
            console.log('Copying .env over...');
            fs.copyFileSync(envPath, path.join(worktreePath, '.env'));
        }

        const mainModelsPath = path.join(projectRoot, 'public', 'models');
        const targetModelsPath = path.join(worktreePath, 'public', 'models');
        if (fs.existsSync(mainModelsPath) && !fs.existsSync(targetModelsPath)) {
            console.log('Symlinking public/models...');
            fs.mkdirSync(path.join(worktreePath, 'public'), { recursive: true });
            execSync(`ln -sf ${mainModelsPath} ${targetModelsPath}`);
        }

        // 3. Setup prompt instructions for Claude
        const promptPath = path.join(worktreePath, `claude_prompt.txt`);
        const promptText = `You are a specialized coding sub-agent. Your task is to resolve exactly ${batchBugs.length} pending bugs in the turms-go codebase.\n`
            + `CRITICAL RULES:\n`
            + `1. Be extremely concise in your output and thoughts. Minimize conversational filler.\n`
            + `2. Check 'git status', 'git diff', and 'git log main..HEAD' first. You might be resuming an interrupted job where some bugs are already partially fixed, staged, or committed in this worktree.\n`
            + `3. As you finish resolving bugs, YOU MUST check them off by changing '- [ ]' to '- [x]' in docs/pending_bugs.md for your specific assigned bugs.\n`
            + `4. Test the implementations, and YOU MUST automatically commit your final changes with a descriptive commit message.\n\n`
            + `BUGS TO FIX:\n====================\n\n`
            + batchBugs.join('\n\n---\n\n')
            + `\n\n====================\nRemember: Always mark completed tasks in docs/pending_bugs.md, and ALWAYS commit the result when finished. KEEP YOUR LOGS AND RESPONSES AS CONCISE AS POSSIBLE.`;

        fs.writeFileSync(promptPath, promptText);

        // 4. 生成包含自动重试策略的 Shell 脚本
        const runnerPath = path.join(worktreePath, 'run_agent.sh');
        const runnerContent = `#!/bin/bash
MAX_RETRIES=15
DELAY=60
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
        git commit -m "fix(automation): resolve parity bugs for batch ${i}" || true
        
        # 2. 回到主分支执行安全的合并。遇到冲突则让 AI 进行 Rebase + Self-Heal 自动修复
        while true; do
            cd "${projectRoot}"
            
            # 使用 mkdir 实现全平台(特别是 macOS)兼容的原子锁排队
            while ! mkdir .git/merge_lock_dir 2>/dev/null; do
                sleep 2
            done
            
            echo "[$(date)] Attempting merge for feature/fix-batch-${i} into main..."
            if git merge "feature/fix-batch-${i}" --no-edit -m "Merge auto-fix batch ${i} into main"; then
                MERGE_RESULT="SUCCESS"
            else
                echo "[!] Conflict detected. Aborting merge."
                git merge --abort
                MERGE_RESULT="CONFLICT"
            fi
            
            # 取出结果后立刻放行排队
            rmdir .git/merge_lock_dir
            
            if [ "$MERGE_RESULT" = "SUCCESS" ]; then
                echo "[$(date)] Successfully merged batch ${i} into main."
                break
            fi
            
            echo "[$(date)] Merge conflict! Initiating Sub-Agent Self-Heal Rebase..."
            cd "${worktreePath}"
            
            # 开始 Rebase main
            git rebase main || {
                # 触发大模型自动修复冲突
                echo "[$(date)] Handing over to Claude for conflict resolution..."
                claude -p --dangerously-skip-permissions "A git rebase conflict was detected in feature/fix-batch-${i}.
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
                if git rebase --show-current-patch >/dev/null 2>&1 || [ -d "${worktreePath}/.git/rebase-merge" ]; then
                    echo "[!] Claude failed to finish the rebase. Aborting pipeline."
                    git rebase --abort
                    exit 1
                fi
            }
            
            echo "[$(date)] Self-Heal Rebase complete. Loop will retry the merge."
        done
        
        echo "[$(date)] Pipeline for batch ${i} complete."
        exit 0
    fi
    
    echo "[$(date)] Agent failed with exit code $EXIT_CODE. Retrying in $DELAY seconds..."
    count=$((count+1))
    sleep $DELAY
done

echo "[$(date)] Agent failed after $MAX_RETRIES attempts."
exit 1
`;
        fs.writeFileSync(runnerPath, runnerContent, { mode: 0o755 });

        // 5. Launch retry script in background
        console.log(`Starting run_agent.sh for batch ${i} in background...`);

        const outLog = fs.openSync(path.join(projectRoot, `claude_worker_batch_${i}.log`), 'a');
        const errLog = fs.openSync(path.join(projectRoot, `claude_worker_batch_${i}.err`), 'a');

        const child = spawn('bash', ['run_agent.sh'], {
            cwd: worktreePath,
            detached: true,
            stdio: ['ignore', outLog, errLog]
        });

        child.unref(); // 允许主进程退出，子代理仍在后台自行运行
        console.log(`[Success] Claude Sub-agent (Batch ${i}) launched with PID ${child.pid}. Logs: claude_worker_batch_${i}.log`);

        if (i < endBatchIndex) {
            // 错峰启动：每个批次拉起后随机延时 5~10 秒，缓解瞬间并发导致的限流
            const delayMs = Math.floor(Math.random() * 5000) + 5000;
            console.log(`Waiting for ${delayMs / 1000}s to avoid rate limits before launching next batch...`);
            await new Promise(resolve => setTimeout(resolve, delayMs));
        }
    }

    console.log("\nAll requested sub-agents have been dispatched!");
})();
