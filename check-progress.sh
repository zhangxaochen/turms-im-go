#!/usr/bin/env bash
# 查阅Turms Parity bugs进度以及 Batch Worktree 进度的脚本

cd "$(dirname "$0")"

MD_FILE="docs/pending_bugs.md"

if [ -f "$MD_FILE" ]; then
    TOTAL_BUGS=$(grep -E '^-[ ]+\[[ x]\]' "$MD_FILE" | wc -l | awk '{print $1}')
    COMPLETED_BUGS=$(grep -E '^-[ ]+\[x\]' "$MD_FILE" | wc -l | awk '{print $1}')
    PENDING_BUGS=$(grep -E '^-[ ]+\[ \]' "$MD_FILE" | wc -l | awk '{print $1}')

    PERCENTAGE=0
    if [ "$TOTAL_BUGS" -gt 0 ]; then
        PERCENTAGE=$(awk "BEGIN {printf \"%.2f\", ($COMPLETED_BUGS/$TOTAL_BUGS)*100}")
    fi

    echo "=========================================================="
    echo "            Turms Parity Bugs 修复全局进度"
    echo "=========================================================="
    echo "总计 Bug 数: $TOTAL_BUGS"
    echo "已修复数量:  $COMPLETED_BUGS"
    echo "剩余未修复:  $PENDING_BUGS"
    echo "完成百分比:  $PERCENTAGE%"
else
    echo "警告: 未找到 $MD_FILE"
fi

echo ""
echo "=========================================================="
echo "               Batch Worktree Progress Report"
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
