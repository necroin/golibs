//go:build windows

package winutils

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/necroin/golibs/utils"
	"github.com/necroin/golibs/utils/winapi"
	"golang.org/x/sys/windows"
)

type Process struct {
	Pid        winapi.ProcessId
	Ppid       winapi.ProcessId
	Executable string
}

func NewWindowsProcess(entry *windows.ProcessEntry32) *Process {
	return &Process{
		Pid:        winapi.ProcessId(entry.ProcessID),
		Ppid:       winapi.ProcessId(entry.ParentProcessID),
		Executable: syscall.UTF16ToString(entry.ExeFile[:]),
	}
}

func (process *Process) String() string {
	return fmt.Sprintf("{ pid: %d, ppid: %d, executable: %s }", process.Pid, process.Ppid, process.Executable)
}

func GetAllProcesses() ([]*Process, error) {
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, fmt.Errorf("[GetAllProcesses] failed get processes snap: %s", err)
	}
	defer windows.CloseHandle(handle)

	entry := windows.ProcessEntry32{}
	entry.Size = uint32(unsafe.Sizeof(entry))

	if err := windows.Process32First(handle, &entry); err != nil {
		return nil, fmt.Errorf("[GetAllProcesses] failed get first process: %s", err)
	}

	result := []*Process{}
	for {
		result = append(result, NewWindowsProcess(&entry))
		if err := windows.Process32Next(handle, &entry); err != nil {
			break
		}
	}

	return result, nil
}

func FindProcessByPid(processes []*Process, pid winapi.ProcessId) *Process {
	processesByParentId := utils.SliceToMultimap[winapi.ProcessId, *Process](processes, func(element *Process) winapi.ProcessId { return element.Pid })
	pidProcesses := processesByParentId[pid]
	if len(pidProcesses) == 0 {
		return nil
	}
	return pidProcesses[0]
}

func FindProcessesByParentPid(processes []*Process, pid winapi.ProcessId) []*Process {
	processesByParentId := utils.SliceToMultimap[winapi.ProcessId, *Process](processes, func(element *Process) winapi.ProcessId { return element.Ppid })
	return processesByParentId[pid]
}
