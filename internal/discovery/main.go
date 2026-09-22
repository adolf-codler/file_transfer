package discovery// {{{

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)// }}}


const (
	PHRASE = "DISOLF_UDPLER" 
)

func Broadcast(ctx context.Context, broadPort string)(error){// {{{
	// initializing broadcast{{{
	broadIP, err:=getSubnet()
	if err!=nil{
		return fmt.Errorf("Subnet Error: %w", err)
	} else {
		fmt.Println("Broadcasting at", broadIP)
	}
	broadAddr:= net.JoinHostPort(broadIP, broadPort)
	addr, err:= net.ResolveUDPAddr("udp4", broadAddr)
	if err != nil{
		return fmt.Errorf("Resolve Error: %w",err)
	}
	conn, err:=net.DialUDP("udp4", nil, addr)
	if err != nil{
		return fmt.Errorf("Dial Error: %w",err)
	}
	defer conn.Close()// }}}

	// periodic broadcasting{{{
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	msg:=[]byte(PHRASE)
	conn.Write(msg)
	fmt.Println("Waiting for Receiver ...")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nReceiver connected! Stopping UDP broadcast...")
			return nil
		case <-ticker.C:
			_, err:=conn.Write(msg)
			if err != nil{
				log.Println("Write Error:",err)
			}
		}
	}// }}}
}// }}}

func ListenBroadcast(broadPort string)(net.UDPAddr, error){// {{{
	// initializing listener{{{
	broadAddr := net.JoinHostPort("", broadPort)
	addr, err:= net.ResolveUDPAddr("udp4", broadAddr)
	if err != nil{
		return net.UDPAddr{}, fmt.Errorf("Resolve Error: %w", err)
	}
	conn, err:=net.ListenUDP("udp4", addr)
	if err != nil{
		return net.UDPAddr{}, fmt.Errorf("Listen Error: %w", err)
	}
	defer conn.Close()// }}}

	fmt.Println("Waiting for Sender ...")
	buf:=make([]byte,1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return net.UDPAddr{}, fmt.Errorf("Reading from udp Error: %w", err)
		}
		msg := string(buf[:n])
		fmt.Printf("Received '%s' from %s\n", msg, remoteAddr)
		if msg == PHRASE{
			return *remoteAddr, nil
		} else{
			continue
		}
	}
}// }}}

func getSubnet()(string, error){// {{{
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.To4()
				mask := ipnet.Mask
				broadcast := net.IP(make([]byte, 4))
				for i := range ip {
					broadcast[i] = ip[i] | ^mask[i]
				}
				return broadcast.String(), nil
			}
		}
	}
	return "255.255.255.255", nil
}// }}}

