package main

import (
	"context"
	"log"
	"net"

	flag "github.com/spf13/pflag"
	"github.com/zahar517/brute-force-attack-prevention/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	host  string
	port  string
	add   = "add"
	black = "black"
	rm    = "rm"
	white = "white"
)

func init() {
	flag.StringVarP(&host, "host", "h", "127.0.0.1", "host")
	flag.StringVarP(&port, "port", "p", "5531", "port")
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 3 {
		log.Fatal("bad args. usage example: ./bin/cli add/rm black/white 192.168.0.1/32")
	}

	command := args[0]
	if command != add && command != rm {
		log.Fatal("use add or rm command")
	}

	list := args[1]
	if list != black && list != white {
		log.Fatal("use black or white list")
	}

	conn, err := grpc.Dial(net.JoinHostPort(host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := grpcserver.NewBFAPToolClient(conn)
	ctx := context.Background()
	sr := &grpcserver.SubnetRequest{Subnet: args[2]}
	var e error

	if command == add && list == black {
		_, e = client.AddToBlacklist(ctx, sr)
	}

	if command == rm && list == black {
		_, e = client.RemoveFromBlacklist(ctx, sr)
	}

	if command == add && list == white {
		_, e = client.AddToWhitelist(ctx, sr)
	}

	if command == rm && list == white {
		_, e = client.RemoveFromWhitelist(ctx, sr)
	}

	if e != nil {
		log.Println(e)
	}
}
