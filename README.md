# Torrentfile
Library for managing Torrent files

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
