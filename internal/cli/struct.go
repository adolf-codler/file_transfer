package cus_cli// {{{

import (
)// }}}


type Cli struct{// {{{
	Type	string		`enum:"send, recv" default:"send" short:"t" help:"send or recv"`
	File 	[]string	`arg:"" help:"all the files that are to be transfered"`
	IP 	string		`short:"i" optional:"" help:"send to or recieve from a certain IP adress"` 
	Port 	int		`short:"p" optional:"" help:"send to or recieve from a certain Port"` 
}// }}}

