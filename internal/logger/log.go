package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	slug string
	base *log.Logger
}

func New(slug string) *Logger {
	return &Logger{
		slug: slug,
		base: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
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
