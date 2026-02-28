#!/usr/bin/env python3
import subprocess
from pathlib import Path
import sys

# compile_gsc.py
# compile only plugin-related GSC scripts for PlutoPlugin

GSCTOOL = r"C:\Users\yalla\Documents\Pkgs\gsc.exe"

# root folder containing your plugin GSC scripts
PLUGIN_ROOT = r"..\t6\scripts"

def compile(script: Path) -> bool:
    cmd = [GSCTOOL, "-m", "comp", "-g", "t6", "-s", "pc", str(script)]
    try:
        print(f"Compiling {script.relative_to(Path(PLUGIN_ROOT).resolve())} ...")
        subprocess.check_call(cmd)
        print(f"✅ {script.name} compiled successfully")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Failed to compile {script.name}, exit code {e.returncode}")
        return False
    except FileNotFoundError:
        print(f"GSCTOOL not found at {GSCTOOL}")
        sys.exit(1)

def main():
    root = Path(PLUGIN_ROOT).resolve()
    if not root.exists():
        print(f"Error: plugin root folder does not exist: {root}")
        sys.exit(1)

    scripts = [s for s in root.rglob("*.gsc") if s.is_file()]
    if not scripts:
        print(f"No plugin .gsc files found in {root}")
        return

    failed = 0
    for script in scripts:
        success = compile(script)
        if not success:
            failed += 1

    print("\nCompilation complete.")
    if failed:
        print(f"{failed} script(s) failed to compile.")
        sys.exit(1)
    else:
        print("All plugin scripts compiled successfully!")

if __name__ == "__main__":
    main()