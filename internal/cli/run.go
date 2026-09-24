package cus_cli

import (
	"adolf-codler/file_transfer/internal/discovery"
	"adolf-codler/file_transfer/internal/sockets"
	"context"
	"log"
)

const(
	DEFAULT_UDP_PORT = "6967"
	DEFAULT_PORT = "6769"
)

func (c *Cli) Run() error{
	if c.Type=="send"{
		ctx, cancel := context.WithCancel(context.Background())
		if err:=discovery.Broadcast(ctx, DEFAULT_PORT); err!=nil{
			return err
		}
		if len(c.File)==1{
			go func(){
				if err := discovery.Broadcast(ctx, DEFAULT_UDP_PORT); err != nil{
					log.Fatalf("Broadcasting error: %s", err)
				}
			}()
			err:=cus_soc.StartServer([]string{"."}, cancel)
			if err!=nil{
				log.Fatalf("Server Error: %s", err)
			}
		} else {
			return nil
		}	
	} else {
		ip, err:=discovery.ListenBroadcast(DEFAULT_UDP_PORT)
		if err!=nil{
			log.Fatalf("Listening Error: %s", err)
		}
		err=cus_soc.StartClient([]string{"."}, ip.IP.String())
		if err!=nil{
			log.Fatalf("Client Error: %s", err)
		}
	}
	return nil
}
