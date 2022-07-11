package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"
)

var Info SysData
var Clipboard string

func main() {
	conf := LoadConfig()
	server, err := net.Listen(
		"tcp",
		fmt.Sprintf("%s:%d", conf.Lhost, conf.Lport))
	Fatal(err)

	//System Info
	go func() {
		for {
			UpdateSysInfo()
			time.Sleep(10 * time.Second)
		}
	}()

	//Clipboard
	go func() {
		for {
			Clipboard = ReadClipboard()
			time.Sleep(300 * time.Millisecond)
		}
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}

		go Handler(conn)

	}
}

func Handler(c net.Conn) {
	r := CompleteResp{SysInfo: Info, ClipboardContent: Clipboard}
	j, _ := json.Marshal(r)
	c.Write(j)
	i := Info
	cb := Clipboard
	buf := make([]byte, 512)

	for {
		time.Sleep(100 * time.Millisecond)
		if i != Info || cb != Clipboard {
			r = CompleteResp{SysInfo: Info, ClipboardContent: Clipboard}
			j, _ = json.Marshal(r)
			c.Write(j)
			i = Info
			cb = Clipboard
		}

		n, err := c.Read(buf)
		if n == 0 || err != nil {
			continue
		} else {
			var output Command
			err := json.Unmarshal(buf[:n], &output)
			if err != nil {
				println(err.Error())
				continue
			}
			if output.CMD != "" {
				PowerAction(output.CMD)
			}
			if output.InsertClipboard != "" {
				Copy(output.InsertClipboard)
			}

		}

	}
}

func UpdateSysInfo() {
	Info.Network = NetStat()
	Info.Device = DeviceInfo()
	Info.Disk = DiskInfo()
}

func LoadConfig() (o Config) {
	b, err := ioutil.ReadFile("config.json")
	Fatal(err)
	err = json.Unmarshal(b, &o)
	Fatal(err)
	return o
}

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
