package main

import (
	"github.com/RenzoFudo/GoCityMarket/cmd/internal/config"
	"github.com/RenzoFudo/GoCityMarket/cmd/internal/server"
	"github.com/RenzoFudo/GoCityMarket/cmd/internal/storage"
	"log"
)

func main() {
	cfg := config.ReadConfig()
	log.Println(cfg)
	storage := storage.New()
	server := server.New(cfg.Host, storage)

	if err := server.Run(); err != nil {
		panic(err)
	}
}

отправить ссылку на github разобрать ошибки что скинет Михаил