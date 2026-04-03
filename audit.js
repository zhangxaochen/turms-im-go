const fs = require('fs');
const path = require('path');
const { spawnSync } = require('child_process');

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
    
    console.log(`\n======================================================`);
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

    // Instead of cat-ing the files, we just ask Gemini to read them itself!
    // This utilizes the Agent's file-reading tools and avoids shell argument/pipe limits.
    const promptText = `Act as a strict code reviewer. Read the original Java code in ${javaPath} and the Go refactor code in ${goPath}. Compare them focusing ONLY on these ported methods: ${checkedMethods.join(', ')}. 

Identify any missing core logic, missing field assignments, or differences in behavior compared to the Java version.
If the Go code implements the Java logic flawlessly for these specific methods, reply with EXACTLY 'NO_BUGS'.
Otherwise, list the specific bugs clearly in markdown format following these rules strictly:
1. For each method that has bugs, create a level 2 heading with the method name (e.g., '## MethodName').
2. List the bugs under the corresponding method heading as a checklist item starting with '- [ ] '.
Do not include introductory or sign-off text.`;

    let success = false;
    let attempts = 0;

    // Retry loop for rate-limiting
    while (!success && attempts < 3) {
        try {
            console.log(`  [Attempt ${attempts + 1}] Invoking Claude API...`);
            
            // We use spawnSync instead of execSync to avoid buffer limits 
            // and we pass parameters as array to avoid shell quoting hell.
            const result = spawnSync('claude', ['-p', promptText], { 
                encoding: 'utf-8', 
                maxBuffer: 10 * 1024 * 1024 // 10MB buffer
            });
            
            const stderr = result.stderr || '';
            const output = result.stdout || '';

            // Handle errors
            if (result.status !== 0 || stderr.includes('429') || stderr.includes('Resource has been exhausted') || stderr.includes('Error')) {
                if (stderr.includes('429') || stderr.includes('exhausted') || output.includes('429')) {
                    console.log(`  [Rate Limit 429] Waiting 65 seconds before retry...`);
                    spawnSync('sleep', ['65']);
                    attempts++;
                } else {
                    console.error(`  [Unexpected Error] Exit code: ${result.status}\nStderr: ${stderr}\nStdout: ${output}`);
                    let errText = `\n# ${currentClass} (Error checking)\n${stderr || output || 'Unknown execution error'}\n`;
                    pendingBugs += errText;
                    fs.appendFileSync(path.join(rootDir, 'docs', 'pending_bugs.md'), errText);
                    success = true; // Error wasn't 429, skip to next
                }
            } else {
                // Success
                if (!output.includes('NO_BUGS')) {
                    console.log(`  [Issue Found] ${currentClass}`);
                    let bugText = `\n# ${currentClass}\n*Checked methods: ${checkedMethods.join(', ')}*\n\n${output.trim()}\n`;
                    pendingBugs += bugText;
                    fs.appendFileSync(path.join(rootDir, 'docs', 'pending_bugs.md'), bugText);
                } else {
                    console.log(`  [OK] ${currentClass}`);
                }
                success = true;
            }
        } catch (e) {
            console.error(`  [Fatal Error] ${e.message}`);
            let fatalText = `\n# ${currentClass} (Fatal error)\n${e.message}\n`;
            pendingBugs += fatalText;
            fs.appendFileSync(path.join(rootDir, 'docs', 'pending_bugs.md'), fatalText);
            success = true; // Skip to next class
        }
    }

    if (!success) {
        console.log(`  [Failed] Could not check ${currentClass} after retries.`);
        return false; // Stop script on permanent rate limit
    }
    
    return true;
}

let proceeded = true;
for (let line of lines) {
    if (!proceeded) break;

    // Match class header
    const classMatch = line.match(/- \*\*(.*?)\*\* \(\[(.*?)\]\((.*?)\)\) ➡️ \[\`(.*?)\`\]\((.*?)\)/);
    if (classMatch) {
        proceeded = checkClass();
        
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
// Trigger the last matched class check
if (proceeded) checkClass();

// At the end, we just check if any bugs were printed
if (!pendingBugs) {
    console.log('No bugs found across all checked methods.');
} else {
    console.log('Audit completed.');
}
