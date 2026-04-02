import re
import os

report_path = "docs/refactor_progress_report.md"

if not os.path.exists(report_path):
    print("Report not found")
    exit(1)

with open(report_path, "r", encoding="utf-8") as f:
    lines = f.readlines()

updates = {}  # file_path -> { method_name -> set(java_signatures) }

for line in lines:
    line = line.strip()
    if line.startswith("- [x] "):
        # Example line: - [x] `verifyAndGrant(UserLoginInfo userLoginInfo)` -> `internal/domain/gateway/access/client/common/ip_request_throttler.go:VerifyAndGrant(ctx context.Context, u int)`
        # Note: Sometimes it has full Go signature: `internal/.../file.go:MethodName(abc)`
        m = re.match(r'- \[x\] `(.*?)` -> `(.*?):([A-Za-z0-9_]+)(?:\(.*?\))?`', line)
        if m:
            java_sig = m.group(1)
            go_file = m.group(2)
            go_method = m.group(3)
            
            if go_file not in updates:
                updates[go_file] = {}
            if go_method not in updates[go_file]:
                updates[go_file][go_method] = set()
            updates[go_file][go_method].add(java_sig)

for go_file, methods in updates.items():
    if not os.path.exists(go_file):
        continue
    
    with open(go_file, "r", encoding="utf-8") as f:
        content = f.read()

    new_lines = []
    file_lines = content.split('\n')
    i = 0
    while i < len(file_lines):
        line = file_lines[i]
        
        # Check if line is a function declaration matching one of our methods
        func_match = re.search(r'^func\s+(?:\([^\)]+\)\s+)?([A-Za-z0-9_]+)\(', line)
        if func_match:
            method_name = func_match.group(1)
            if method_name in methods:
                java_sigs = list(methods[method_name])
                # Check if tags already exist right above
                existing_tags = []
                j = len(new_lines) - 1
                while j >= 0 and new_lines[j].strip().startswith("// @MappedFrom"):
                    existing_tags.append(new_lines[j].strip()[15:].strip())
                    j -= 1
                
                for sig in java_sigs:
                    # Ignore exact duplicates or already annotated
                    if sig not in existing_tags:
                        new_lines.append(f"// @MappedFrom {sig}")
        
        new_lines.append(line)
        i += 1
    
    # Save back
    with open(go_file, "w", encoding="utf-8") as f:
        f.write('\n'.join(new_lines))

print("Injected tags")
