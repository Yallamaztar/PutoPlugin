#!/usr/bin/env python3
import re
import sys
from pathlib import Path

# clean_log.py
# cleans the log file from ANSI color codes: \033

def remove_ansi_colors(content: str) -> str:
    return re.sub(r'\x1B\[[0-9;]*[mK]', '', content)

def read_log(file: Path) -> str:
    with file.open("r", encoding="utf-8") as f:
        return f.read()

def write_log(file: Path, content: str) -> None:
    with file.open("w", encoding="utf-8") as f:
        f.write(content)

def main() -> None:
    if len(sys.argv) < 2:
        print("Usage: clean_log.py <path_to_log_file>")
        sys.exit(1)

    log_path = Path(sys.argv[1])
    if not log_path.exists() or not log_path.is_file():
        print(f"Error: file not found: {log_path}")
        sys.exit(1)

    content = read_log(log_path)
    write_log(log_path, remove_ansi_colors(content))
    print(f"Cleaned ANSI codes from {log_path}")

if __name__ == '__main__':
    main()