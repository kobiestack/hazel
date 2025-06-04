package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		panic(err)
	}

	var writer io.Writer
	if os.Getenv("ENV") != "dev" {
		writer, err = os.OpenFile("hazel.log", os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			panic(err)
		}
	} else {
		writer = os.Stdout
	}

	logger := slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	slog.SetDefault(logger)
}
