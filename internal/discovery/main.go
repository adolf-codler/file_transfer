package discovery

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	PHRASE = "DISOLF_UDPLER" 
)
func Broadcast(ctx context.Context, broadPort string){// {{{
	broadIP, err:=getSubnet()
	if err!=nil{
		fmt.Println("getSubnet error:", err)
	} else {
		fmt.Println("Broadcasting at", broadIP)
	}
	broadAddr:= net.JoinHostPort(broadIP, broadPort)
	addr, err:= net.ResolveUDPAddr("udp4", broadAddr)
	if err != nil{
		log.Println("Resolve Error:",err)
	}
	conn, err:=net.DialUDP("udp4", nil, addr)
	if err != nil{
		log.Println("Dial Error:",err)
	}
	defer conn.Close()

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	// Broadcast immediately the first time
	msg:=[]byte(PHRASE)
	conn.Write(msg)
	fmt.Println("Waiting for Receiver ...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nReceiver connected! Stopping UDP broadcast...")
			return
		case <-ticker.C:
			_, err:=conn.Write(msg)
			if err != nil{
				log.Println("Write Error:",err)
			}
		}
	}
}// }}}

func ListenBroadcast(broadPort string)(net.UDPAddr){// {{{
	broadAddr := net.JoinHostPort("", broadPort)
	addr, err:= net.ResolveUDPAddr("udp4", broadAddr)
	if err != nil{
		log.Println("Resolve Error", err)
	}
	conn, err:=net.ListenUDP("udp4", addr)
	if err != nil{
		log.Println("Listen Error:", err)
	}
	defer conn.Close()
	fmt.Println("Waiting for Sender ...")
	buf:=make([]byte,1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error reading:", err)
			continue
		}
		msg := string(buf[:n])
		fmt.Printf("Received '%s' from %s\n", msg, remoteAddr)
		if msg == PHRASE{
			return *remoteAddr
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

