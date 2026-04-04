import sys

def main():
    lines_to_fix = set()
    with open('/tmp/checkboxes2.txt', 'r') as f:
        for line in f:
            if line.strip():
                parts = line.split(":", 1)
                if len(parts) == 2 and parts[0].isdigit():
                    lines_to_fix.add(int(parts[0]))

    with open('docs/pending_bugs.md', 'r') as f:
        content = f.readlines()

    for idx, num in enumerate(lines_to_fix):
        idx_0 = num - 1
        if "- [ ]" in content[idx_0]:
            content[idx_0] = content[idx_0].replace("- [ ]", "- [x]", 1)

    with open('docs/pending_bugs.md', 'w') as f:
        f.writelines(content)
    
    print(f"Fixed {len(lines_to_fix)} checkboxes.")

if __name__ == "__main__":
    main()
