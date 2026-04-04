const fs = require('fs');
const path = require('path');

const mdPath = path.join(__dirname, 'docs/refactor_progress_report.md');
let content = fs.readFileSync(mdPath, 'utf8');

const regex = /- \[x\] \`(.*?)\` -> \[(.*?):(.*?)\]\((.*?)\)/g;

let updatedContent = content.replace(regex, (match, javaMethod, goPathStr, goMethodSig, goLink) => {
    // goMethodSig example: ChannelRegistered(isAvailable bool)
    const methodNameMatch = goMethodSig.match(/^([A-Za-z0-9_]+)/);
    if (!methodNameMatch) {
         console.warn(`Could not extract method name from: ${goMethodSig}`);
         return `- [x] \`${javaMethod}\` -> [${goMethodSig}](${goLink})`;
    }
    const methodName = methodNameMatch[1];
    
    // goLink: ../internal/domain/gateway/access/client/common/service_availability.go
    // In __dirname, files are at relative to __dirname. 
    // goLink usually starts with ../, e.g., ../internal, from docs directory.
    // However, if we run this in turms-go root (__dirname), the file is actually simply goLink.replace('../', './')
    
    let filePath = goLink.replace('../', './');
    if (filePath.includes('#')) {
       filePath = filePath.split('#')[0];
    }
    
    let lineNum = 1;
    if (fs.existsSync(filePath)) {
        const fileLines = fs.readFileSync(filePath, 'utf8').split('\n');
        for (let i = 0; i < fileLines.length; i++) {
            const line = fileLines[i];
            // Match `func (receiver) MethodName(` or `func MethodName(`
            const funcRegex = new RegExp(`func\\s*(?:\\([^)]+\\)\\s*)?${methodName}\\s*\\(`, 'i');
            if (funcRegex.test(line)) {
                lineNum = i + 1;
                break;
            }
        }
    } else {
        console.warn(`File not found: ${filePath}`);
    }
    
    const newLink = `${goLink.split('#')[0]}#L${lineNum}`;
    return `- [x] \`${javaMethod}\` -> [${goMethodSig}](${newLink})`;
});

// Write it back
fs.writeFileSync(mdPath, updatedContent, 'utf8');
console.log('Update complete.');

