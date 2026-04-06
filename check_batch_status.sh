#!/usr/bin/env bash

# check_batch_status.sh
# 自动统计所有 turms-worker-batch-* 目录下的 temp_task.md 进度

echo "=========================================================="
echo "          Batch Worktree Progress Report"
echo "=========================================================="

# 遍历上一级目录中所有的 turms-worker-batch-* 文件夹
for dir in ../turms-worker-batch-*; do
  if [ -d "$dir" ] && [ -f "$dir/temp_task.md" ]; then
    dirname=$(basename "$dir")
    done=$(grep -c "\[x\]" "$dir/temp_task.md")
    pending=$(grep -c "\[ \]" "$dir/temp_task.md")
    
    # 状态判定
    status=""
    if [ "$pending" -eq 0 ] && [ "$done" -gt 0 ]; then
      status="✅ DONE (完工)"
    elif [ "$done" -eq 0 ] && [ "$pending" -gt 0 ]; then
      status="🛑 BLOCKED (一点没动)"
    else
      status="⚠️ PARTIAL (干到一半中断)"
    fi

    # 格式化输出
    printf "%-25s | %3d done | %3d pending | %s\n" "$dirname" "$done" "$pending" "$status"
  fi
done
echo "=========================================================="
