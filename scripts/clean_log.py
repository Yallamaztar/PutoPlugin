import re

def remove_ansi_colors(content: str) -> str:
    return re.sub(r'\x1B\[[0-9;]*[mK]', '', content)

def read_log(file: str) -> str:
    with open(file, "r") as f:
        return f.read()

def write_log(file: str, content: str) -> None:
    with open(file, "w") as f:
        f.write(content)

if __name__ == '__main__':
    file    = input("enter path to your log file: ")
    content = read_log(file)
    clean   = remove_ansi_colors(content)
    write_log(file, clean)
