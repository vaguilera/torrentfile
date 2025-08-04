package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/vaguilera/torrentfile"
)

func main() {
	cmd := &cli.Command{
		Name:  "TorrentInfo",
		Usage: "Shows torrent file data",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose mode",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				cli.ShowAppHelpAndExit(cmd, 0)
			}
			verbose := cmd.Bool("verbose")
			return showTorrentInfo(cmd.Args().Get(0), verbose)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func showTorrentInfo(file string, verbose bool) error {
	torrent, err := torrentfile.Open(file)
	if err != nil {
		return err
	}

	infoHash, _ := torrent.GetInfoHash()

	fmt.Printf("Torrent: \t%s\nCreated by: \t%s on %s\nEncoding: \t%s\nComment: \t%s\nPrivate: \t%t\nSource: \t%s\nMain Tracker: \t%s\nInfoHash: \t%x\n\n",
		torrent.Info.Name,
		torrent.CreatedBy,
		torrent.CreationDate.Format(time.RFC3339),
		torrent.Encoding,
		torrent.Comment,
		torrent.Info.Private,
		torrent.Info.Source,
		torrent.Announce,
		infoHash,
	)
	fmt.Printf("Piece length: \t%d\nTotal length: \t%d\n", torrent.Info.PieceLength, torrent.Info.Length)

	if len(torrent.Info.Files) > 0 {
		fmt.Println("files:")
		for _, f := range torrent.Info.Files {
			var path string
			for _, p := range f.Path {
				path = path + "/" + p
			}
			fmt.Printf("%s size: %d MD5: %s\n", path, f.Length, f.MD5sum)

		}
	}
	if verbose {
		fmt.Println("Url List:")
		for _, url := range torrent.UrlList {
			fmt.Printf("\t%s\n", url)
		}
		fmt.Println("Secondary Trackers:")
		for _, tiers := range torrent.AnnounceList {
			for _, tracker := range tiers {
				fmt.Printf("\t%s\n", tracker)
			}
		}
	}
	return nil

}
