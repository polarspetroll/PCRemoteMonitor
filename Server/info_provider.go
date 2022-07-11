package main

import (
	"fmt"
	"golang.design/x/clipboard"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"syscall"
)

func initi() {
	err := clipboard.Init()
	if err != nil {
		log.Println(err)
	}
}

func NetStat() string {
	ifcs, _ := net.Interfaces()
	var format string
	for _, v := range ifcs {

		adr, _ := v.Addrs()
		format += v.Name + "\n"
		for _, vv := range adr {
			format += fmt.Sprintf("   %s\n", vv.String())
		}
	}
	return format
}

func DeviceInfo() string {
	hostname, _ := os.Hostname()
	usr, _ := user.Current()

	return fmt.Sprintf("  Username: %s\n  Name:%s\n  Host Name: %s\n",
		usr.Username,
		usr.Name,
		hostname,
	)
}

func ReadClipboard() string {
	return string(clipboard.Read(clipboard.FmtText))
}

func Copy(s string) {
	clipboard.Write(clipboard.FmtText, []byte(s))
}

func DiskInfo() string {
	var disk DiskStatus
	fs := syscall.Statfs_t{}
	syscall.Statfs("/", &fs)
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return fmt.Sprintf("  Total: %.2f GB\n  Used: %.2f GB\n  Free: %.2f GB\n",
		float64(disk.All)/float64(1073741824),
		float64(disk.Used)/float64(1073741824),
		float64(disk.Free)/float64(1073741824),
	)
}

func PowerAction(action string) {
	if action == "shutdown" {
		if runtime.GOOS == "windows" {
			exec.Command("shutdown", "/s").Output()
		} else {
			exec.Command("shutdown", "-h", "now").Output()
		}
	}

	if action == "suspend" {
		if runtime.GOOS == "windows" {
			exec.Command("rundll32.exe", "powrprof.dll,", "SetSuspendState", "Sleep").Output()
		} else if runtime.GOOS == "darwin" {
			exec.Command("pmset", "sleepnow").Output()
		} else {
			exec.Command("systemctl", "suspend").Output()
		}
	}

}
