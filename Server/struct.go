package main

type Config struct {
	Lhost string `json:"Lhost"`
	Lport int    `json:"Lport"`
}

type SysData struct {
	Network string `json:"network"`
	Device  string `json:"device"`
	Disk    string `json:"disk"`
}

type DiskStatus struct {
	All  uint64 `json:"All"`
	Used uint64 `json:"Used"`
	Free uint64 `json:"Free"`
}

type Command struct {
	CMD             string `json:"command"`
	InsertClipboard string `json:"clipboard"`
}

type CompleteResp struct {
	SysInfo          SysData `json:"system"`
	ClipboardContent string  `json:"clipboard"`
}
