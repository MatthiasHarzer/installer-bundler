package windowsutil

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modShell32               = windows.NewLazySystemDLL("shell32.dll")
	procSHGetKnownFolderPath = modShell32.NewProc("SHGetKnownFolderPath")
)

func shGetKnownFolderPath(rfid *windows.GUID, dwFlags uint32, hToken windows.Handle) (string, error) {
	var pszPath *uint16
	r, _, err := procSHGetKnownFolderPath.Call(
		uintptr(unsafe.Pointer(rfid)),
		uintptr(dwFlags),
		uintptr(hToken),
		uintptr(unsafe.Pointer(&pszPath)),
	)
	if r != 0 {
		return "", err
	}
	defer windows.CoTaskMemFree(unsafe.Pointer(pszPath))
	return windows.UTF16PtrToString(pszPath), nil
}

func GetKnownFolderPath(id KnownFolderID) (string, error) {
	guid, err := windows.GUIDFromString(string(id))
	if err != nil {
		return "", err
	}
	return shGetKnownFolderPath(&guid, 0, 0)
}
