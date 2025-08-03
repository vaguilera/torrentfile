package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v3"
	"github.com/vaguilera/torrentfile"
)

func main() {
	cmd := &cli.Command{
		Name:  "TorrentInfo",
		Usage: "Create torrent from template",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				cli.ShowAppHelpAndExit(cmd, 0)
			}
			return createTorrentFile(cmd.Args().Get(0))
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func createTorrentFile(path string) error {
	b, err := os.ReadFile("cmd/create/template.toml")
	if err != nil {
		log.Fatal(err)
	}

	var cfg torrentfile.TConfig
	err = toml.Unmarshal(b, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	data, err := torrentfile.NewTorrentFile(&cfg, path)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("testo.torrent", data, 0644)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
