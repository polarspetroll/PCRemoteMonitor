package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io"
	"net"
	"strings"
	"time"
)

var (
	Preferences fyne.Preferences
	Clipboard   fyne.Clipboard
	QueueCMD    string
	a           fyne.App
)

func main() {
	cn := make(chan bool)
	var address string
	var initial_container *fyne.Container
	a = app.NewWithID("main1")
	w := a.NewWindow("PC Monitor")
	Preferences = a.Preferences()

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Address(IP:PORT)")
	label := widget.NewLabel("Please enter the address")
	statusLable := widget.NewLabel("")
	statusLable.Alignment = 1
	label.Alignment = 1
	var button *widget.Button
	button = widget.NewButton("Connect",
		func() {
			button.Disable()
			if entry.Text == "" {
				statusLable.Text = "Please enter a valid IP address"
				statusLable.Refresh()
				button.Enable()
				return
			}
			address = entry.Text
			Preferences.SetString("PrevAddress", address)
			var runloop bool = true

			go func() {
				statusLable.Text = "Connecting   "
				statusLable.Refresh()
				time.Sleep(300 * time.Millisecond)
				for runloop {
					for i := 1; i < 4; i++ {
						if <-cn {
							return
						}
						statusLable.Text = strings.Replace(statusLable.Text, " ", ".", 1)
						statusLable.Refresh()
						time.Sleep(500 * time.Millisecond)
					}
					statusLable.Text = statusLable.Text[:10] + "   "
					statusLable.Refresh()
				}

			}()
			time.Sleep(500 * time.Millisecond)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				runloop = false
				cn <- true
				statusLable.Text = "Unable to connect"
				statusLable.Refresh()
			} else {
				runloop = false
				initial_container.Hide()
				w.SetOnClosed(func() { conn.Close() })
				go Core(w, conn)
			}
			button.Enable()
		},
	)

	address = Preferences.String("PrevAddress")
	if address != "" {
		saved_label := widget.NewLabel("Previous connection")
		saved_button := widget.NewButton(address, func() {
			entry.Text = address
			entry.Refresh()
		})
		initial_container = container.NewVBox(label, entry, button, statusLable, saved_label, saved_button)
	} else {
		initial_container = container.NewVBox(label, entry, button, statusLable)
	}

	w.SetContent(initial_container)
	w.ShowAndRun()
}

func Core(window fyne.Window, c net.Conn) {
	var out DataPack
	l1 := widget.NewLabel(out.SysInfo.Network)
	l2 := widget.NewLabel(out.SysInfo.Device)
	l3 := widget.NewLabel(out.SysInfo.Disk)

	networkCard := widget.NewCard("Network", "", l1)
	deviceCard := widget.NewCard("Device Info", "", l2)
	diskCard := widget.NewCard("Disk Status", "", l3)
	infoContainer := container.NewVBox(networkCard, deviceCard, diskCard)

	suspendBtn := widget.NewButton("Suspend", func() { Send(c, "suspend", "") })
	shutdownBtn := widget.NewButton("Shutdown", func() { Send(c, "shutdown", "") })
	powerContainers := container.NewVBox(suspendBtn, shutdownBtn)

	tab1 := container.NewTabItem("Info", infoContainer)
	tab2 := container.NewTabItem("Power Actions", powerContainers)

	window.SetContent(container.NewAppTabs(tab1, tab2))

	Clipboard = window.Clipboard()
	buf := make([]byte, 512)
	written_clipboard := Clipboard.Content()
	Send(c, QueueCMD, written_clipboard)
	QueueCMD = ""

	var current_content string
	var n int
	var err error
	for {
		c.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
		n, err = c.Read(buf)
		if err == io.EOF {
			fail := DataPack{
				ClipboardContent: "",
				SysInfo:          SysData{Network: "-", Device: "-", Disk: "-\n\n\nDevice Is Offline"},
			}
			UpdateInfo(fail, l1, l2, l3)
			a.SendNotification(fyne.NewNotification("Alert", "Device Was Disconnected"))
			break
		}
		if n != 0 {
			err = json.Unmarshal(buf[:n], &out)
			if err != nil {
				goto CLIPBOARD
			}
			if out.ClipboardContent != "" {
				Clipboard.SetContent(out.ClipboardContent)
			}
			UpdateInfo(out, l1, l2, l3)

		}

	CLIPBOARD:
		current_content = Clipboard.Content()
		if current_content != "" && written_clipboard != current_content {
			Send(c, "", current_content)
			written_clipboard = current_content
			current_content = ""
			time.Sleep(10 * time.Millisecond)
		}

		if QueueCMD != "" {
			Send(c, QueueCMD, "")
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func Send(c net.Conn, cmd, cboard string) {
	data := Command{CMD: cmd, InsertClipboard: cboard}
	j, err := json.Marshal(data)
	if err != nil {
		return
	}
	c.Write(j)
}

func UpdateInfo(info DataPack, l1, l2, l3 *widget.Label) {
	l1.Text = info.SysInfo.Network
	l2.Text = info.SysInfo.Device
	l3.Text = info.SysInfo.Disk

	l1.Refresh()
	l2.Refresh()
	l3.Refresh()

}
