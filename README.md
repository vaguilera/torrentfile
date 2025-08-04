# Torrentfile
Library for managing Torrent files and magnet links

### Parse and create torrent files

Include a couple of CLI tools as an example of use

# Usage

## Reading Torrent Files

You can read a file from disk of from memory

```go
tFile, err := Open("path/example.torrent") // Open torrent file from path
tFile, err := FromMem(buffer) // Open torrent file from memory 
```

TFile is a struct of type `TorrentFile` containing the parsed fields

## Creating Torrent Files

```go
 buffer, err := NewTorrentFile(&config, "path/file.torrent")
```

Create a new torrent file based on the provided config. Second argument can be a single file or a folder for a multi-file torrents.

If a folder is specified as a second argument, library will traverse the whole folder structure adding every file to the torrent file

You can specify the current fields in config:

```go
type TConfig struct {
    Announce     string
    AnnounceList [][]string
    Encoding     string
    Comment      string
    CreatedBy    string
    UrlList      []string
    PieceSize    int
    Md5          bool
    Private      bool
}
```

Piece Size is 32Kb by default. Some trackers will need some specific setting here depending on the size of the torrent.file

## Magnet Links

Library allows to create and read Magnet links

```go
    func ParseMagnetURI(uri string) (*MagnetLink, error) // Returns a struct with the parsed Magnet Link
    func NewMagnetURI(magnet *MagnetLink) string // Returns a string with the generted Magnet URI
```

Fields supported by the library:

- xt (Exact Topic)
- dn (Display Name)
- tr (Address Tracker)
- xl (Exact Length)
- kt (Keyword Topic)
- as (Acceptable Sources)
- xs (Exact Sources)
- mt (Manifest Topic)

Note: Most of the BitTorrent clients accept only a single "mt" field. But as this is not strictly defined I decided to allow multiple manifest topics. If you just want to support 1 read the first element of the array