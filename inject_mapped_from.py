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
        # Format 1 (backticks): - [x] `verify(UserLoginInfo)` -> `internal/file.go:Verify()`
        # Format 2 (markdown link): - [x] `verify(UserLoginInfo)` -> [internal/file.go:Verify()](../internal/file.go)
        # We can extract the Java side from inside backticks.
        java_match = re.search(r'- \[x\] `(.*?)` \-> ', line)
        if not java_match:
            continue
        java_sig = java_match.group(1)
        
        # Now try to extract the Go file and method
        # It could be `internal/file.go:MethodName(...)` OR [internal/file.go:MethodName(...)](...)
        go_side_match = re.search(r'->\s+\[?(.*?):([A-Za-z0-9_]+)(?:\(.*?\))?\]?\(?.*?\)?', line)
        
        # We need a more precise match
        go_side_str = line[java_match.end():].strip()
        
        # Try markdown link format: [file.go:MethodName(args...)](...)
        mlink_match = re.match(r'\[(.*?):([A-Za-z0-9_]+)(?:\(.*?\))?\]\(.*?\)', go_side_str)
        
        # Try backtick format: `file.go:MethodName(args...)`
        btick_match = re.match(r'`(.*?):([A-Za-z0-9_]+)(?:\(.*?\))?`', go_side_str)
        
        if mlink_match:
            go_file = mlink_match.group(1)
            go_method = mlink_match.group(2)
        elif btick_match:
            go_file = btick_match.group(1)
            go_method = btick_match.group(2)
        else:
            print(f"Failed to parse Go side in line: {line}")
            continue
            
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
