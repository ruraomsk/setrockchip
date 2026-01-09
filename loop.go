package main

import (
	"fmt"

	"github.com/ruraomsk/setrockchip/command"
	"github.com/ruraomsk/setrockchip/scp"
	"golang.org/x/crypto/ssh"
)

var conn *ssh.Client
var allgood = false

func loopMain(host string, port int, user string, password string, setupclear bool) {
	isRunning = true
	defer func() {
		isRunning = false
		if !allgood {
			StatusMessage("[red] Во время выполнения обнаружена ошибка")
		} else {
			StatusMessage("[green] Готово")
		}
	}()
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	LogMessage("Подключаемся к устройству")
	conn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		ErrorMessage(fmt.Sprintf("Failed to dial: %s", err.Error()))
		return
	}
	defer conn.Close()
	command.Connection(conn)

	LogMessage("Определяем операционную систему")
	name, err := command.GetSystem()
	if err != nil {
		ErrorMessage(err.Error())
		return
	}
	if name != "RK356X" {
		ErrorMessage("Устройство не rockchip!")
		return
	}

	LogMessage(fmt.Sprintf("Устройство с операционной системой:%s", name))
	LogMessage("Останавливаем potop на целевой системе")
	progress(1, 9)
	err = command.KillProc("gobanana")
	if err != nil {
		ErrorMessage(err.Error())
		return
	}
	err = command.KillProc("potop")
	if err != nil {
		ErrorMessage(err.Error())
		return
	}
	err = command.KillProc("vpot")
	if err != nil {
		ErrorMessage(err.Error())
		return
	}
	LogMessage("Подключаем систему sfp")
	progress(2, 9)
	err = scp.Connection(conn)
	if err != nil {
		ErrorMessage(err.Error())
		return
	}
	if setupclear {
		LogMessage("Удаляем прошлые настройки программ")
		LogMessage("Очищаем каталог с программами")
		progress(3, 9)
		err = command.DeleteDir("/home/rura")
		if err != nil {
			ErrorMessage(err.Error())
			return
		}
	} else {
		LogMessage("Сохраняем прошлые настройки программ")
		LogMessage("Очищаем логи программ")
		progress(3, 9)
		err = command.DeleteDir("/home/rura/log")
		if err != nil {
			ErrorMessage(err.Error())
			return
		}

	}
	LogMessage("Создаем каталоги для программ")
	progress(4, 9)
	err = command.CreateDir("/home/rura")
	if err != nil {
		ErrorMessage(err.Error())
		return
	}
	err = command.CreateDir("/etc/qt5")
	if err != nil {
		ErrorMessage(err.Error())
		return
	}
	LogMessage("Записываем скрипты запуска программ")
	progress(5, 9)
	datas, _ := resource.ReadFile("torockchip/gopotop.sh")
	err = scp.WriteFile("/home/rura/gopotop.sh", datas, true)
	if err != nil {
		ErrorMessage(err.Error())
		return
	}
	datas, _ = resource.ReadFile("torockchip/govpot.sh")
	err = scp.WriteFile("/home/rura/govpot.sh", datas, true)
	if err != nil {
		ErrorMessage(err.Error())
		return
	}
	datas, _ = resource.ReadFile("torockchip/eglfs_kms.json")
	err = scp.WriteFile("/etc/qt5/eglfs_kms.json", datas, false)
	if err != nil {
		ErrorMessage(err.Error())
		return
	}

	LogMessage("Записываем скрипт автоматичекого запуска программ")
	progress(6, 9)
	datas, _ = resource.ReadFile("torockchip/S99potop")
	err = scp.WriteFile("/etc/init.d/potop", datas, true)
	if err != nil {
		ErrorMessage(err.Error())
		return
	}
	LogMessage("Записываем сервисы автоматического запуска программ")
	progress(7, 9)
	// datas, _ = resource.ReadFile("torockchip/potop.service")
	// err = scp.WriteFile("/etc/systemd/system/potop.service", datas, false)
	// if err != nil {
	// 	ErrorMessage(err.Error())
	// 	return
	// }
	// datas, _ = resource.ReadFile("torockchip/vpot.service")
	// err = scp.WriteFile("/etc/systemd/system/vpot.service", datas, false)
	// if err != nil {
	// 	ErrorMessage(err.Error())
	// 	return
	// }
	LogMessage("Записываем программу potop")
	progress(8, 9)
	err = scp.CopyFile("potop", "/home/rura/potop", true)
	if err != nil {
		ErrorMessage(err.Error())
		return
	}
	LogMessage("Записываем программу vpot")
	progress(8, 9)
	err = scp.CopyFile("vpot", "/home/rura/vpot", true)
	if err != nil {
		ErrorMessage(err.Error())
		return
	}
	LogMessage("Включаем сервисы автоматического запуска программ")
	progress(9, 9)

	LogMessage("Устройство нужно перезагрузить")
	allgood = true
}
