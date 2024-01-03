package main

import (
	"log"
	"net"

	// "net/http"
	// "net/rpc"
	// "os"
	// "os/signal"
	// "syscall"

	pb "spacemesher/proto"

	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)




func main(){
	_, err := flags.Parse(&tags)
	if err != nil {
		return
	}

	err = NewPlot(tags.Datadir)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	listener, err := net.Listen("tcp", tags.Listen)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	server = grpc.NewServer()
	reflection.Register(server)
	pb.RegisterSpacemesherServer(server, &Server{})
	log.Printf("Server: %s is up", tags.Listen)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("error: %v", err)
	}
}