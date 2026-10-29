package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "uadata"

	"github.com/mrscorpio/uachecker/pkg/influxdbold"
	"github.com/mrscorpio/uahelper/pkg/tgbot"
	"google.golang.org/grpc"
)

type uaServiceServer struct {
	pb.UnimplementedUaServiceServer

	Vals     map[string]any
	Tagnames []string
	Wd       int32
	Alive    bool

	DbCl  *influxdbold.InfluxCl
	TgBot *tgbot.TgBot
}

func (s *uaServiceServer) GetData(ctx context.Context, req *pb.SrvRequest) (*pb.SrvResponse, error) {
	s.Wd = req.Wd

	if len(s.Tagnames) > 0 && len(req.V) >= len(s.Tagnames) {
		for i, v := range s.Tagnames {
			s.Vals[v] = req.V[i]
		}
		if s.DbCl != nil {
			tm, err := time.Parse("2006-01-02 15:04:05.000 -07", req.Tm)
			if err != nil {
				log.Println(err)
			}

			err = s.DbCl.AddPoint(s.Vals, tm)
			if err != nil {
				log.Println(err)
			}
		}
		if s.Wd%666 == 0 {
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

func (s *uaServiceServer) SendTg(ctx context.Context, req *pb.ToTelegram) (*pb.TelegramResponse, error) {
	var err error
	if s.TgBot != nil {
		err = s.TgBot.SendTxt(req.Mes, req.ChatId)
		if err != nil {
			return &pb.TelegramResponse{Er: err.Error()}, err
		}
		return &pb.TelegramResponse{Er: "message sent"}, nil
	}
	return &pb.TelegramResponse{Er: "telegram bot is nil"}, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tgBoss := os.Getenv("TGBOSS")

	bot, err := tgbot.NewBot(ctx, os.Getenv("TGTOKEN"), false)
	if err != nil {
		log.Println(err)
	}

	dbCl, err := influxdbold.NewDbConn("http://localhost:8086", os.Getenv("INFLUX_USER"), os.Getenv("INFLUX_PASS"))
	if err != nil {
		log.Println(err)
	}

	listener, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalln(err)
	}
	vals := make(map[string]any)
	tagnames := make([]string, 0)
	grpcServ := grpc.NewServer()
	uaServ := &uaServiceServer{Vals: vals, Tagnames: tagnames, DbCl: dbCl, TgBot: bot}
	pb.RegisterUaServiceServer(grpcServ, uaServ)
	log.Println(net.InterfaceAddrs())

	go func() {
		for {
			wdBuf := uaServ.Wd
			time.Sleep(22 * time.Second)
			if uaServ.Wd == wdBuf && uaServ.Alive {
				if bot != nil {
					err := bot.SendTxt("логер вырубился", tgBoss)
					if err != nil {
						log.Println(err)
					}
				}
				log.Println("uahelper is dead, last wd =", uaServ.Wd)
				uaServ.Alive = false
			}
			if uaServ.Wd != wdBuf && !uaServ.Alive {
				uaServ.Alive = true
				if bot != nil {
					err := bot.SendTxt("логер в строю", tgBoss)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}()

	if err := grpcServ.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
