const fs = require('fs');

const path = '/Users/11176728/gemini-cli/dev-turms-im-refactor/turms-go/docs/pending_bugs.md';
let lines = fs.readFileSync(path, 'utf8').split('\n');

let insideBerBuffer = false;
let checkCount = 0;

for (let i = 0; i < lines.length; i++) {
    // Start tracking when we hit BerBuffer.java or any of the consecutive ones
    if (lines[i].trim() === '# BerBuffer.java') {
        insideBerBuffer = true;
    }
    
    // Stop tracking when we hit SearchRequest since our changes stopped before it
    // Wait, we implemented Filter.Write which is before SearchRequest.java
    if (lines[i].trim() === '# SearchRequest.java') {
        insideBerBuffer = false;
        break;
    }
    
    if (insideBerBuffer && lines[i].includes('- [ ]')) {
        lines[i] = lines[i].replace('- [ ]', '- [x]');
        checkCount++;
    }
}

fs.writeFileSync(path, lines.join('\n'));
console.log(`Successfully checked off ${checkCount} items.`);
