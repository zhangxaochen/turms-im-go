#!/bin/bash
# 查阅Turms Parity bugs进度的脚本

cd "$(dirname "$0")"

MD_FILE="docs/pending_bugs.md"

if [ ! -f "$MD_FILE" ]; then
    echo "错误: 未找到 $MD_FILE"
    exit 1
fi

TOTAL_BUGS=$(grep -E '^-[ ]+\[[ x]\]' "$MD_FILE" | wc -l | awk '{print $1}')
COMPLETED_BUGS=$(grep -E '^-[ ]+\[x\]' "$MD_FILE" | wc -l | awk '{print $1}')
PENDING_BUGS=$(grep -E '^-[ ]+\[ \]' "$MD_FILE" | wc -l | awk '{print $1}')

if [ "$TOTAL_BUGS" -eq 0 ]; then
    echo "未在 $MD_FILE 中发现记录。"
    exit 0
fi

PERCENTAGE=$(awk "BEGIN {printf \"%.2f\", ($COMPLETED_BUGS/$TOTAL_BUGS)*100}")

echo "========================================="
echo " Turms Parity Bugs 修复进度反馈"
echo "========================================="
echo "总计 Bug 数: $TOTAL_BUGS"
echo "已修复数量:  $COMPLETED_BUGS"
echo "剩余未修复:  $PENDING_BUGS"
echo "完成百分比:  $PERCENTAGE%"
echo "========================================="
