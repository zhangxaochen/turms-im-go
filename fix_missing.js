const fs = require('fs');
const cp = require('child_process');

let md = fs.readFileSync('docs/refactor_progress_report.md', 'utf8');
const lines = md.split('\n');

let currentClass = '';
const result = [];

for (let i = 0; i < lines.length; i++) {
  let line = lines[i];

  const classMatch = line.match(/^- \*\*(.*?)\.java\*\*/);
  if (classMatch) {
    currentClass = classMatch[1];
    
    const searchCmd = `grep -Hn -r "// @MappedFrom ${currentClass}\\b" internal/ || true`;
    const files = cp.execSync(searchCmd, { encoding: 'utf8' }).trim().split('\n').filter(x => x);
    let matchedFile = null;
    if (files.length > 0) {
        matchedFile = files[0].split(':')[0];
    }
    if (matchedFile && !line.includes('➡️')) {
        const originalLinkMatch = line.match(/^(- \*\*.*?\*\* \(\[.*?\]\(.*?\)\))/);
        if (originalLinkMatch) {
            line = `${originalLinkMatch[1]} ➡️ [\`${matchedFile}\`](../${matchedFile})`;
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
        line = `${indent}- [x] \`${methodSig}\` ➡️ [\`${matchedFile}\`](../${matchedFile})`;
      }
    }
  }

  result.push(line);
}

fs.writeFileSync('docs/refactor_progress_report.md', result.join('\n'));
