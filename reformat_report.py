import re

report_path = "docs/refactor_progress_report.md"

with open(report_path, "r", encoding="utf-8") as f:
    lines = f.readlines()

new_lines = []
for line in lines:
    line_stripped = line.strip()
    # Looking for lines like: - [x] `create(...)` -> `internal/domain/.../file.go:Function(...)`
    # or empty `-> ` cmd/turms-gateway/main.go:main()`
    
    # We want to replace the `...go:Method(...)` with [...go:Method(...)](../...go)
    # But only if it's currently using backticks.
    
    match = re.search(r'-\s+\[(x| )\]\s+`(.*?)`\s+->\s+`(.*?):(.*?)`', line)
    if match:
        check = match.group(1)
        java_sig = match.group(2)
        go_path = match.group(3)
        go_method = match.group(4)
        
        # fix the specific error with ThrowableInfo
        if "ThrowableInfo" in java_sig and go_method.startswith("Create("):
            go_method = go_method.replace("Create(", "CreateFromError(")
            
        new_line = f"  - [{check}] `{java_sig}` -> [{go_path}:{go_method}](../{go_path})\n"
        new_lines.append(new_line)
        continue
        
    new_lines.append(line)

with open(report_path, "w", encoding="utf-8") as f:
    f.write("".join(new_lines))

print("Done reformatting.")
