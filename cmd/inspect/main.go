package main

import (
	"fmt"
	"github.com/vaguilera/torrentfile"
	"time"
)

func main() {
	torrent, err := torrentfile.LoadTorrentFile("./fixtures/test.torrent")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Torrent: \t%s\nCreated by: \t%s on %s\nEncoding: \t%s\nComment: \t%s\nMain Tracker: \t%s\n\n",
		torrent.Info.Name,
		torrent.CreatedBy,
		torrent.CreationDate.Format(time.RFC3339),
		torrent.Encoding,
		torrent.Comment,
		torrent.Announce,
	)
	fmt.Printf("Piece length: \t%d\nTotal length: \t%d\n", torrent.Info.PieceLength, torrent.Info.Length)

	if len(torrent.Info.Files) > 0 {
		fmt.Println("Files:")
		for _, f := range torrent.Info.Files {
			var path string
			for _, p := range f.Path {
				path = path + "/" + p
			}
			fmt.Printf("%s size: %d MD5: %s\n", path, f.Length, f.MD5sum)

		}
	}
	//fmt.Println("Secondary Trackers:")
	//for _, tiers := range torrent.AnnounceList {
	//	for _, tracker := range tiers {
	//		fmt.Printf("%s\n", tracker)
	//	}
	//}
	//fmt.Println()

}
