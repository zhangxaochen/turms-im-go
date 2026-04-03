const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const rootDir = process.cwd();
const mdPath = path.join(rootDir, 'docs', 'refactor_progress_report.md');
const content = fs.readFileSync(mdPath, 'utf-8');

const lines = content.split('\n');

let currentClass = null;
let currentJavaFile = null;
let currentGoFile = null;
let checkedMethods = [];
let pendingBugs = '';

function checkClass() {
    if (!currentClass || checkedMethods.length === 0) return true;
    
    console.log(`Checking ${currentClass} (${checkedMethods.length} methods)...`);
    const javaPath = path.resolve(rootDir, 'docs', currentJavaFile);
    const goPath = path.resolve(rootDir, 'docs', currentGoFile);
    
    if (!fs.existsSync(javaPath)) {
        console.log(`  [Skipped] Java file missing: ${javaPath}`);
        return true;
    }
    if (!fs.existsSync(goPath)) {
        console.log(`  [Skipped] Go file missing: ${goPath}`);
        return true;
    }
    
    let javaCode = fs.readFileSync(javaPath, 'utf-8');
    let goCode = fs.readFileSync(goPath, 'utf-8');
    
    if (javaCode.length > 8000) javaCode = javaCode.substring(0, 8000) + '\n...[TRUNCATED]';
    if (goCode.length > 8000) goCode = goCode.substring(0, 8000) + '\n...[TRUNCATED]';

    const promptText = `I am refactoring a Java project to Go. Act as a strict code reviewer. 
Compare the Go implementation against the original Java code.
Focus ONLY on these methods: ${checkedMethods.join(', ')}.

Java File:
\`\`\`java
${javaCode}
\`\`\`

Go File:
\`\`\`go
${goCode}
\`\`\`

Are there any missing logic, bugs, or differences in behavior compared to the Java version?
If the Go code implements the Java logic flawlessly for these methods, reply with EXACTLY 'NO_BUGS'.
Otherwise, list the bugs clearly in markdown bullet points. Do not include introductory text, just the issues.`;

    fs.writeFileSync('temp_prompt.txt', promptText);
    
    let success = false;
    let attempts = 0;
    while (!success && attempts < 3) {
        try {
            console.log(`  Invoking AI reviewer (attempt ${attempts + 1})...`);
            const cmd = `gemini -y -e none -p "$(cat temp_prompt.txt)"`;
            const output = execSync(cmd, { encoding: 'utf-8', stdio: ['pipe', 'pipe', 'pipe'] });
            
            if (!output.includes('NO_BUGS')) {
                console.log(`  [Issue Found] ${currentClass}`);
                pendingBugs += `\n### ${currentClass}\nChecked methods: ${checkedMethods.join(', ')}\n${output.trim()}\n`;
            } else {
                console.log(`  [OK] ${currentClass}`);
            }
            success = true;
        } catch (e) {
            const stderr = e.stderr ? e.stderr.toString() : e.message;
            if (stderr.includes('429') || stderr.includes('Resource has been exhausted')) {
                console.log(`  [Rate Limit 429] Waiting 60 seconds...`);
                // Wait to avoid completely spamming out
                execSync('sleep 60');
                attempts++;
            } else {
                console.error(`  [Error] ${stderr}`);
                pendingBugs += `\n### ${currentClass} (Error checking)\n${stderr}\n`;
                success = true; // Skip to next
            }
        }
    }

    if (!success) {
        console.log(`  [Failed] Could not check ${currentClass} after retries.`);
        return false; // Stop script on permanent rate limit
    }
    
    return true;
}

let proceed = true;
for (let line of lines) {
    if (!proceed) break;

    // Match class header
    const classMatch = line.match(/- \*\*(.*?)\*\* \(\[(.*?)\]\((.*?)\)\) ➡️ \[\`(.*?)\`\]\((.*?)\)/);
    if (classMatch) {
        proceed = checkClass();
        
        currentClass = classMatch[1];
        currentJavaFile = classMatch[3];
        currentGoFile = classMatch[5];
        checkedMethods = [];
        continue;
    }
    
    // Match checked method: e.g. - [x] `method()`
    const methodMatch = line.match(/^\s*- \[x\] \`(.*?)\`/);
    if (methodMatch) {
         checkedMethods.push(methodMatch[1]);
    }
}
if (proceed) checkClass();

if (pendingBugs) {
    const bugsFile = path.join(rootDir, 'pending_bugs.md');
    const existing = fs.existsSync(bugsFile) ? fs.readFileSync(bugsFile, 'utf-8') : '';
    fs.writeFileSync(bugsFile, existing + '\n\n## Auto-Review Findings\n' + pendingBugs);
    console.log('Saved bugs to pending_bugs.md');
} else {
    console.log('No bugs found across all checked methods.');
}

if (fs.existsSync('temp_prompt.txt')) fs.unlinkSync('temp_prompt.txt');
