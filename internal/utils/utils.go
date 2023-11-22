package utils

import (
	"encoding/json"
	"os"
	"os/exec"
	"syscall"
	"unicode/utf8"
	"unsafe"
)

var (
	config Config
)

func Clear() {
	c := exec.Command("cmd", "/c", "cls")
	c.Stdout = os.Stdout
	c.Run()
}

func ContainsEmoji(s string) bool {
	for _, r := range s {
		if utf8.RuneLen(r) > 3 {
			return true
		}
	}
	return false
}

func SetTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}

	defer syscall.FreeLibrary(handle)
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}

	newTitle, _ := syscall.UTF16PtrFromString(title)
	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(newTitle)), 0, 0)
	return int(r), err
}

func ReadConfig(data []byte) (*Config, error) {
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
