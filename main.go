package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "uadata"

	"github.com/mrscorpio/uahelper/pkg/tgbot"
	"google.golang.org/grpc"
)

type uaServiceServer struct {
	pb.UnimplementedUaServiceServer

	Vals     map[string]float32
	Tagnames []string
	Wd       int32
}

func (s *uaServiceServer) GetData(ctx context.Context, req *pb.SrvRequest) (*pb.SrvResponse, error) {
	s.Wd = req.Wd

	if len(s.Tagnames) > 0 && len(req.V) >= len(s.Tagnames) {
		for i, v := range s.Tagnames {
			s.Vals[v] = req.V[i]
		}
		if s.Wd%100 == 0 {
			log.Println(s.Vals)
		}
	}

	if len(req.Tags) > 0 {
		s.Tagnames = make([]string, len(req.Tags))
		copy(s.Tagnames, req.Tags)
	}

	//log.Println(req.V)
	return &pb.SrvResponse{
		Wd: req.Wd + 66,
	}, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tgToken := os.Getenv("TGTOKEN")
	tgBoss := os.Getenv("TGBOSS")
	bot, err := tgbot.NewBot(ctx, tgToken, false)
	if err != nil {
		log.Println(err)
	}

	listener, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalln(err)
	}
	vals := make(map[string]float32)
	tagnames := make([]string, 0)
	grpcServ := grpc.NewServer()
	uaServ := &uaServiceServer{Vals: vals, Tagnames: tagnames}
	pb.RegisterUaServiceServer(grpcServ, uaServ)
	log.Println(net.InterfaceAddrs())

	go func() {
		for {
			wdBuf := uaServ.Wd
			time.Sleep(22 * time.Second)
			if uaServ.Wd == wdBuf {
				if bot != nil {
					err := bot.SendTxt("логер вырубился", tgBoss)
					if err != nil {
						log.Println(err)
					}
				}
				log.Println("uahelper is dead, last wd =", uaServ.Wd)
				break
			}
		}
	}()

	if err := grpcServ.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
