package main

type SysData struct {
	Network string `json:"network"`
	Device  string `json:"device"`
	Disk    string `json:"disk"`
}

type DataPack struct {
	SysInfo          SysData `json:"system"`
	ClipboardContent string  `json:"clipboard"`
}

type Command struct {
	CMD             string `json:"command"`
	InsertClipboard string `json:"clipboard"`
}
