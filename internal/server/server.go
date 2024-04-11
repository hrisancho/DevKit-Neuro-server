package server

import (
	"DevKit-Neuro-server/internal/config"
	main_logger "DevKit-Neuro-server/internal/logger"
	"DevKit-Neuro-server/internal/validator"
	controller "DevKit-Neuro-server/proto"
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
)

type Server struct {
	App       *net.UDPConn
	Config    config.Config
	Logger    *main_logger.Logger
	Validator *validator.AppValidator
}

func NewServer(config config.Config, logger *main_logger.Logger, validator *validator.AppValidator) (server *Server, err error) {

	udpAddr, err := net.ResolveUDPAddr("udp", config.Addr)

	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalln(err)
	}

	server = &Server{
		App:       conn,
		Config:    config,
		Logger:    logger,
		Validator: validator,
	}

	return
}

func (server Server) Run(ctx context.Context) {

	var buf [1024]byte
	// Read from UDP listener in endless loop
	for {
		n, addr, err := server.App.ReadFromUDP(buf[:])
		if err != nil {
			log.Fatalln(err)
		}
		data := &controller.RawDataPack{}
		err = proto.Unmarshal(buf[0:n], data)
		if err != nil {
			log.Fatalln(err)
		}

		log.Print("> ", data)

		// Write back the message over UPD
		server.App.WriteToUDP([]byte("Hello UDP Client\n"), addr)
	}
}
