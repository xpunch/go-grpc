package main

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/micro/go-grpc"
	hello "github.com/micro/go-grpc/examples/greeter/server/proto/hello"
	"github.com/micro/go-micro"
	mtls "github.com/micro/util/go/lib/tls"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	cert, _ := mtls.Certificate("127.0.0.1:9090")
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			cert,
		},
	}

	service := grpc.NewService(
		micro.Name("go.micro.srv.greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		grpc.WithTLS(tlsConfig),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	hello.RegisterSayHandler(service.Server(), new(Say))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
