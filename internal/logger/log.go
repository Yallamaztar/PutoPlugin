package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	slug string
	base *log.Logger
	file *os.File
}

func New(slug, logPath string) *Logger {
	writers := []io.Writer{os.Stdout}

	var file *os.File
	if logPath != "" {
		var err error
		file, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			writers = append(writers, file)
		} else {
			fmt.Printf("Failed to open log file %s: %v\n", logPath, err)
		}
	}

	multi := io.MultiWriter(writers...)
	return &Logger{
		slug: slug,
		base: log.New(multi, "", log.LstdFlags|log.Lshortfile),
		file: file,
	}
}

func (l *Logger) prefix() string {
	return fmt.Sprintf("[%s]", l.slug)
}

func (l *Logger) Println(v ...any) {
	l.base.Println(append([]any{l.prefix()}, v...)...)
}

func (l *Logger) Printf(format string, v ...any) {
	l.base.Printf(l.prefix()+format, v...)
}

func (l *Logger) Fatal(v ...any) {
	l.base.Fatal(append([]any{l.prefix()}, v...)...)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.base.Fatalf(l.prefix()+format, v...)
}

func (l *Logger) Warnf(format string, v ...any) {
	l.base.Printf(l.prefix()+" [WARN] "+format+"\n", v...)
}

func (l *Logger) Warnln(v ...any) {
	l.base.Println(append([]any{l.prefix(), "[WARN]"}, v...)...)
}
