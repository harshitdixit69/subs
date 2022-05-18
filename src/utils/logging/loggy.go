package loggy

import (
	"fmt"
	"github.com/fatih/color"
	"runtime"
	"strings"
	"time"
)

func printLabel(label string, clr color.Attribute, contents ...interface{}) {
	mainColor := color.New(clr)
	cyan := color.New(color.FgCyan)

	_, _ = mainColor.Printf("[%v] -- ", label)
	_, _ = cyan.Print(time.Now().Format("15:04.05"))
	_, _ = mainColor.Print(" -- ")
	fmt.Println(getCallReference())
	out := fmt.Sprint(contents...)
	for _, line := range strings.Split(out, "\n") {
		_, _ = mainColor.Print("│ ")
		fmt.Println(line)
	}

	_, _ = mainColor.Println("└--------------------")
}

func Warning(v ...interface{}) {
	printLabel("WARNING", color.FgYellow, v...)
}

func Error(v ...interface{}) {
	printLabel("ERROR", color.FgRed, v...)
}

func Info(v ...interface{}) {
	printLabel("INFO", color.FgGreen, v...)
}

func WTF(v ...interface{}) {
	printLabel("WTF", color.FgBlue, v...)
}

func Debug(v ...interface{}) {
	printLabel("Debug", color.FgHiYellow, v...)
}

func getCallReference() string {
	_, file, line, ok := runtime.Caller(3)
	if ok {
		return fmt.Sprintf("%v:%v", file, line)
	}
	return "<unknown-file>"
}
