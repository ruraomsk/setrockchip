package main

import (
	"embed"
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var err error
var port int
var buildDate string
var logLines []string
var maxLogLines = 100
var isRunning bool

//go:embed torockchip

var resource embed.FS
var app *tview.Application
var form *tview.Form
var logView *tview.TextView
var statusView *tview.TextView
var check = false

func main() {
	app = tview.NewApplication()

	// Включаем поддержку мыши
	app.EnableMouse(true)

	// Переменные для данных

	// Создаем форму для ввода параметров
	form = tview.NewForm().
		AddInputField("IP Устройства", "192.168.88.100", 30, nil, nil).
		AddInputField("Номер порта SSH ", "22", 10, nil, nil).
		AddInputField("Пользователь", "root", 10, nil, nil).
		AddInputField("Пароль", "1234", 10, nil, nil).AddCheckbox("Сбросить настройки", check, func(checked bool) {
		check = checked
	})

	form.SetBorder(true).SetTitle(fmt.Sprintf(" Установка программы potop %s на RockChip. Параметры подключения ", buildDate))

	// Создаем текстовое поле для логов
	logView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	logView.SetBorder(true).SetTitle(" Лог выполнения ")

	// Включаем прокрутку для лога
	logView.SetScrollable(true)

	// Создаем текстовое поле для статуса
	statusView = tview.NewTextView().
		SetDynamicColors(true)
	statusView.SetBorder(true).SetTitle(" Статус ")

	// Функция для добавления лога

	// Канал для остановки выполнения

	// Функция выполнения задачи
	executeTask := func() {

		if isRunning {
			addLog("[yellow]Установка уже выполняется")
			return
		}

		// Получаем значения из формы
		host := form.GetFormItem(0).(*tview.InputField).GetText()
		sport := form.GetFormItem(1).(*tview.InputField).GetText()
		user := form.GetFormItem(2).(*tview.InputField).GetText()
		password := form.GetFormItem(3).(*tview.InputField).GetText()

		// Парсим параметры

		port, err = strconv.Atoi(sport)
		if err != nil {
			port = 22
		}
		if host == "" {
			host = "192.168.88.100"
		}
		if user == "" {
			user = "root"
		}
		if password == "" {
			password = "1234"
		}

		// Обновляем UI в основном потоке
		app.QueueUpdateDraw(func() {
			addLog("[blue]IP Устройства: [white]%s [blue]Порт SSH: [white]%d", host, port)
			addLog("[blue]Пользователь: [white]%s [blue]Пароль: [white]%s ", user, password)
			if check {
				addLog("[blue]Настройки: [white]%s ", "Сбросить")
			} else {
				addLog("[blue]Настройки: [white]%s ", "Сохранить")
			}
			updateStatus("[yellow]Подготовка к выполнению...")
		})

		loopMain(host, port, user, password, check)

	}

	// Добавляем кнопки в форму
	form.AddButton("Запуск", func() {
		// Запускаем в горутине
		go executeTask()
	})

	form.AddButton("Очистить лог", func() {
		clearLog()
	})

	form.AddButton("Выход", func() {
		app.Stop()
	})

	// Упрощенная обработка мыши для прокрутки лога
	logView.SetDoneFunc(func(key tcell.Key) {
		// Обработка клавиш для прокрутки
		switch key {
		case tcell.KeyUp:
			row, _ := logView.GetScrollOffset()
			logView.ScrollTo(row-1, 0)
		case tcell.KeyDown:
			row, _ := logView.GetScrollOffset()
			logView.ScrollTo(row+1, 0)
		case tcell.KeyHome:
			logView.ScrollToBeginning()
		case tcell.KeyEnd:
			logView.ScrollToEnd()
		}
	})

	// Создаем основную компоновку
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(form, 15, 0, true)       // Форма
	flex.AddItem(logView, 0, 1, false)    // Лог (расширяемый)
	flex.AddItem(statusView, 3, 0, false) // Статус

	// Настраиваем обработчики клавиш
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlL:
			clearLog()
			return nil
		case tcell.KeyCtrlQ:
			app.Stop()
			return nil
		case tcell.KeyF1:
			// Фокус на форму
			app.SetFocus(form)
			return nil
		case tcell.KeyF2:
			// Фокус на лог
			app.SetFocus(logView)
			return nil
		}
		return event
	})

	// Устанавливаем начальный статус
	updateStatus("[blue]Готов к работе - используйте мышь или клавиатуру")
	addLog("[green]Добро пожаловать!")
	addLog("[blue]Управление:")
	addLog("[white]• Нажимайте кнопки мышью")
	addLog("[white]• Прокрутка лога: колесо мыши или клавиши Up/Down")
	addLog("[white]• Ctrl+L: очистка лога")
	addLog("[white]• Esc: остановка задачи")
	addLog("[white]• Ctrl+Q: выход")
	addLog("[white]• F1: фокус на форму, F2: фокус на лог")
	addLog("")

	// Запускаем приложение
	if err := app.SetRoot(flex, true).SetFocus(form).Run(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
}
