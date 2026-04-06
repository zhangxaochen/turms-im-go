const fs = require('fs');

const tempTaskPath = process.argv[2];
const globalDocPath = process.argv[3];

if (!fs.existsSync(tempTaskPath)) process.exit(0);

const localTask = fs.readFileSync(tempTaskPath, 'utf8');
let globalDoc = fs.readFileSync(globalDocPath, 'utf8');

// 把从 localTask 提取出来的 [x] 项，通过准确的字符串查找在 globalDoc 里全局置换。
const completedTasks = [...localTask.matchAll(/^- \[x\] (.*)$/gm)].map(m => m[1].trim());

for (const task of completedTasks) {
    // 处理正则转义，防止任务文本里包含特殊符号导致匹配失败
    const safeTask = task.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'); 
    const regex = new RegExp(`^- \\[ \\] ${safeTask}`, 'gm');
    globalDoc = globalDoc.replace(regex, `- [x] ${task}`);
}

fs.writeFileSync(globalDocPath, globalDoc);
console.log(`[Sync] Synced ${completedTasks.length} completed tasks to global doc.`);
