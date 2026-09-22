package cus_cli

import(
	"github.com/alecthomas/kong"
)


type Cli struct{
	Type	string		`arg:"" enum:"send, recv" default:"send" help:"send or recv"`
	IP 	string		`short:"i" optional:"" help:"send to or recieve from a certain IP adress"` 
	Port 	int		`short:"p" optional:"" help:"send to or recieve from a certain Port"` 
	File 	[]string	`arg:"" help:"all the files that are to be transfered"`
}

func NewKong()*Cli{
	var cli Cli
	return &cli
}
