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

    const today = new Date().toISOString().split('T')[0].replace(/-/g, '');
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
            currentBug = [line]; // 严格复制整行，不额外拼接 Context 等前缀，防止 sync 脚本报错或混淆
        } else if (line.trim().match(/^- \[[xX]\]/)) {
            // 遇见已完成的任务，将其作为边界，把之前的未完成任务存入，并清空当前块（丢弃已完成的任务及其附带说明内容）
            if (currentBug.length > 0) {
                bugList.push(currentBug.join('\n'));
                currentBug = [];
            }
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

        const worktreeName = `turms-worker-${today}-batch-${i}`;
        const worktreePath = path.resolve(projectRoot, '..', worktreeName);
        const branchName = `feature/${today}/fix-batch-${i}`;

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

        // 2.5 Generate Local Scratchpad (temp_task.md)
        const tempTaskPath = path.join(worktreePath, 'temp_task.md');
        if (!fs.existsSync(tempTaskPath)) {
            const taskContent = "# Local Progress Tracker for Batch " + i + "\n\n" + batchBugs.join('\n\n---\n\n');
            fs.writeFileSync(tempTaskPath, taskContent);
        }

        const excludePath = path.join(projectRoot, '.git', 'info', 'exclude');
        if (fs.existsSync(excludePath)) {
            const excludes = fs.readFileSync(excludePath, 'utf8');
            if (!excludes.includes('temp_task.md')) {
                fs.appendFileSync(excludePath, '\ntemp_task.md\n');
            }
        }

        // 3. Setup prompt instructions for Claude
        const promptPath = path.join(worktreePath, `claude_prompt.txt`);
        const promptText = `You are a specialized coding sub-agent resolving bugs in the turms-go codebase.\n`
            + `CRITICAL RULES:\n`
            + `1. Your specific assigned bugs are listed securely in the file 'temp_task.md' located in the root of your workspace. Read 'temp_task.md' immediately to figure out what you need to do.\n`
            + `2. THIS IS A FRESH BATCH. Treat it as an independent task. Even if you have worked on similar batches before, DO NOT assume any bug is already fixed unless it is checked off in YOUR local 'temp_task.md'.\n`
            + `3. Check 'git status', 'git diff', and 'git log main..HEAD' first! You might be resuming an interrupted execution where some bugs are already fixed or partially staged.\n`
            + `4. As you fix each bug, YOU MUST open 'temp_task.md' and change its '- [ ]' to '- [x]'. This file acts as your single source of truth for resumption.\n`
            + `5. Very Important: You ONLY need to check off ('- [x]') the bugs in your local 'temp_task.md'. Under NO CIRCUMSTANCES should you modify 'docs/pending_bugs.md'. The main scheduler will sync it automatically.\n`
            + `6. ABSOLUTELY CRITICAL: When updating 'temp_task.md', you MUST ONLY change '- [ ]' to '- [x]'. DO NOT rewrite, reformat, or change a single word of the task description text. Changing the text will break the global regex sync script.\n`
            + `7. The pipeline is only considered complete when ALL tasks in 'temp_task.md' are checked off. BEFORE finishing, you MUST run 'go build ./...' to verify no compilation errors exist globally!\n`
            + `8. Finally, use 'git add -A' (NOT git commit -a) so that ALL newly created files are staged, and then 'git commit' with a neat descriptive message.\n`
            + `KEEP LOGS CONCISE. Stop and commit when all tasks in the scratchpad are fully resolved.`;

        fs.writeFileSync(promptPath, promptText);

        // 4. 生成包含自动重试策略的 Shell 脚本
        const runnerPath = path.join(worktreePath, 'run_agent.sh');
        const runnerContent = `#!/bin/bash
MAX_RETRIES=100
count=0

echo "Starting Claude Sub-agent..."
while [ $count -lt $MAX_RETRIES ]; do
    echo "[$(date)] Attempt $((count+1)) / $MAX_RETRIES"
    
    # 强制检查 temp_task.md 是否已经全部打钩了
    if ! grep -q "\\- \\[ \\]" temp_task.md; then
        echo "[$(date)] All tasks in temp_task.md are checked off! Skipping claude execution."
        EXIT_CODE=0
    else
        # 执行 claude 命令
        claude -p --dangerously-skip-permissions "$(cat claude_prompt.txt)"
        
        if grep -q "\\- \\[ \\]" temp_task.md; then
            echo "[!] temp_task.md still has unfinished tasks (- [ ]). Claude agent did not complete everything!"
            EXIT_CODE=1
        else
            echo "[Checking Validation] Running 'go build ./...' to verify codebase integrity..."
            if ! go build ./...; then
                echo "[!] Compilation failed. Rejecting claim of completion and sending back to Claude."
                # 回滚已打完的对勾，强制 Claude 重来并修复编译问题
                # 兼容 macOS 的 sed 用法
                sed -i '' 's/- \\[x\\]/- \\[ \\]/g' temp_task.md || true
                EXIT_CODE=1
            else
                echo "[Success] Code compiles cleanly. All tasks in temp_task.md are checked off."
                EXIT_CODE=0
            fi
        fi
    fi
    
    if [ $EXIT_CODE -eq 0 ]; then
        echo "[$(date)] Agent finished successfully. Attempting to commit and merge..."
        # 1. 兜底提交（防止 Claude 忘了主动跑 Commit 命令，强制使用 -A 避免漏掉新文件）
        git add -A
        git commit -m "fix(automation): resolve parity bugs for batch ${i}" || true
        
        # 2. 回到主分支执行安全的合并。遇到冲突则让 AI 进行 Rebase + Self-Heal 自动修复
        while true; do
            cd "${projectRoot}"
            
            # 使用 mkdir 实现全平台(特别是 macOS)兼容的原子锁排队
            while ! mkdir .git/merge_lock_dir 2>/dev/null; do
                sleep 2
            done
            
            echo "[$(date)] Attempting merge for ${branchName} into main..."
            if git merge "${branchName}" --no-edit -m "Merge auto-fix ${today} batch ${i} into main"; then
                MERGE_RESULT="SUCCESS"
            else
                echo "[!] Conflict detected. Aborting merge."
                git merge --abort
                MERGE_RESULT="CONFLICT"
            fi
            
            if [ "$MERGE_RESULT" = "SUCCESS" ]; then
                echo "[$(date)] Successfully merged batch ${i} into main."
                
                # --- 主控在持有合法的原根目录读写锁时同步最新的全量状态 ---
                node "${projectRoot}/sync_markdown_status.js" "${worktreePath}/temp_task.md" "${projectRoot}/docs/pending_bugs.md"
                
                git add docs/pending_bugs.md
                git commit -m "docs: sync global progress from batch ${i}" || true
                
                rmdir .git/merge_lock_dir
                break
            fi
            
            # 若失败则只释放锁，进入重试流
            rmdir .git/merge_lock_dir
            
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
    
    # Delay strategy: 1min - 2min - 1min - 2min ... (count % 2)
    if [ $((count % 2)) -eq 0 ]; then
        DELAY=60
    else
        DELAY=120
    fi
    echo "[$(date)] Pipeline did not complete. Retrying in $((DELAY / 60)) minute(s)..."
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
