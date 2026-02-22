package main

import (
	"fmt"
	"strings"
	"time"
)

var addLog = func(format string, args ...any) {
	message := fmt.Sprintf("[white][%s] %s", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
	logLines = append(logLines, message)

	// Ограничиваем количество строк
	if len(logLines) > maxLogLines {
		logLines = logLines[len(logLines)-maxLogLines:]
	}

	logView.SetText(strings.Join(logLines, "\n"))
	logView.ScrollToEnd()
}
var LogMessage = func(message string) {
	app.QueueUpdateDraw(func() {
		addLog(fmt.Sprintf("[green]%s", message))
	})
}
var ErrorMessage = func(message string) {
	app.QueueUpdateDraw(func() {
		addLog(fmt.Sprintf("[red]%s", message))
	})
}
var StatusMessage = func(message string) {
	app.QueueUpdateDraw(func() {
		updateStatus(fmt.Sprintf("[yellow]%s", message))
	})
}

// Функция для обновления статуса
var updateStatus = func(message string) {
	statusView.SetText(message)
}

// Функция для очистки лога
var clearLog = func() {
	if isRunning {
		addLog("[yellow]Нельзя очистить лог во время выполнения задачи")
		return
	}
	logLines = []string{}
	logView.SetText("")
	updateStatus("Лог очищен")
}

func progress(stage int, stages int) {
	updateStatus(fmt.Sprintf("Выпоняется %02d / %02d", stage, stages))
}
