const fs = require('fs');
const cp = require('child_process');

let md = fs.readFileSync('docs/refactor_progress_report.md', 'utf8');

// replace all "➡️ [`.../group_repositories_stubs.go`](...)" with empty and `[x]` with `[ ]`
md = md.replace(/- \[x\] (.*?) ➡️ \[\`internal\/domain\/group\/repository\/group_repositories_stubs\.go\`\].*/g, '- [ ] $1');

const lines = md.split('\n');

let currentClass = '';
const result = [];

for (let i = 0; i < lines.length; i++) {
  let line = lines[i];

  const classMatch = line.match(/^- \*\*(.*?)\.java\*\*/);
  if (classMatch) {
    currentClass = classMatch[1];
    
    // Check if the current class has a MappedFrom in internal
    // Be careful with GroupBlocklistRepository since we mapped it in group_blocked_user_repository.go
    const searchCmd = `grep -Hn -r "// @MappedFrom ${currentClass}\\b" internal/ || true`;
    const files = cp.execSync(searchCmd, { encoding: 'utf8' }).trim().split('\n').filter(x => x);
    let matchedFile = null;
    if (files.length > 0) {
        matchedFile = files[0].split(':')[0];
    }
    
    // Also try to fix mutated lines that might have lost their original links.
    // In our previous bad script run, we had `- **ClassName.java** ([java/im/turms/...](../turms-orig/turms-service/src/...)) ➡️ [\`xyz\`](../xyz)`
    if (matchedFile) {
        // Find existing link component
        const originalLinkMatch = line.match(/^(- \*\*.*?\*\* \(\[.*?\]\(.*?\)\))/);
        if (originalLinkMatch) {
            line = `${originalLinkMatch[1]} ➡️ [\`${matchedFile}\`](../${matchedFile})`;
        } else {
            // fallback if it was messed up
            line = `- **${currentClass}.java** (unlinked) ➡️ [\`${matchedFile}\`](../${matchedFile})`;
        }
    }
  }

  const methodMatch = line.match(/^(\s*)- \[ \] `(.*?)`/);
  if (methodMatch) {
    const indent = methodMatch[1];
    let methodSig = methodMatch[2];
    
    let methodNameMatch = methodSig.match(/^([a-zA-Z0-9_]+)\b/);
    if (methodNameMatch) {
      let methodName = methodNameMatch[1];
      let goMethodName = methodName.charAt(0).toUpperCase() + methodName.slice(1);
      
      const searchCmd = `grep -Hn -r "${goMethodName}(" internal/domain/group/ internal/infra/validator/ || true`;
      const files = cp.execSync(searchCmd, { encoding: 'utf8' }).trim().split('\n').filter(x => x);

      let matchedFile = null;
      if (files.length > 0) {
        matchedFile = files[0].split(':')[0];
      }

      if (matchedFile) {
        let lineNum = 1;
        if (fs.existsSync(matchedFile)) {
            const fileLines = fs.readFileSync(matchedFile, 'utf8').split('\n');
            for (let j = 0; j < fileLines.length; j++) {
                const funcRegex = new RegExp(`func\\s*(?:\\([^)]+\\)\\s*)?${goMethodName}\\s*\\(`, 'i');
                if (funcRegex.test(fileLines[j])) {
                    lineNum = j + 1;
                    break;
                }
            }
        }
        line = `${indent}- [x] \`${methodSig}\` -> [${goMethodName}()](../${matchedFile}#L${lineNum})`;
      }
    }
  }

  result.push(line);
}

fs.writeFileSync('docs/refactor_progress_report.md', result.join('\n'));
