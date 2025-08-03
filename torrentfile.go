package torrentfile

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	bencode "github.com/IncSW/go-bencode"
)

const createdBy = "TorrentFile lib"

type SHA1 [20]byte

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

type MultiFileInfo struct {
	Length int64
	Path   []string
	MD5sum string
}

type Info struct {
	Pieces      []SHA1
	PieceLength int64
	Length      int64
	Name        string
	Private     bool
	Source      string
	Files       []MultiFileInfo
}

type TorrentFile struct {
	Announce     string
	AnnounceList [][]string
	CreationDate time.Time
	Encoding     string
	Info         Info
	Comment      string
	CreatedBy    string
	UrlList      []string

	rawData map[string]interface{}
}

func Open(filename string) (*TorrentFile, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return FromMem(f)

}

func FromMem(buffer []byte) (*TorrentFile, error) {
	rawTFile, err := bencode.Unmarshal(buffer)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling torrent file: %w", err)
	}

	t := &TorrentFile{}
	var ok bool
	t.rawData, ok = rawTFile.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid torrent file")
	}
	if err = t.unmarshalTFile(); err != nil {
		return nil, fmt.Errorf("error unmarshalling torrent file: %w", err)
	}

	return t, nil
}

func NewTorrentFile(t *TConfig, path string) ([]byte, error) {
	tInterface := map[string]interface{}{}
	info := map[string]interface{}{}

	if t.Announce != "" {
		tInterface["announce"] = t.Announce
	}
	if t.AnnounceList != nil {
		tInterface["announce-list"] = t.AnnounceList
	}
	if t.Encoding != "" {
		tInterface["encoding"] = t.Encoding
	}
	if t.Comment != "" {
		tInterface["comment"] = t.Comment
	}
	if t.UrlList != nil {
		tInterface["url-list"] = t.UrlList
	}

	tInterface["creation date"] = time.Now().UnixMilli()
	tInterface["created by"] = createdBy

	if t.Private == true {
		info["private"] = true
	}

	files, err := getFiles(path)
	if err != nil {
		return nil, err
	}
	switch len(files) {
	case 0:
		return nil, errors.New("no files found")
	case 1:
		info["name"] = files[0].Path[len(files[0].Path)-1]
	}

	pSize := 32
	if t.PieceSize > 0 {
		pSize = t.PieceSize
	}

	pieces, err := calculateParts(files, pSize)
	if err != nil {
		return nil, err
	}

	info["pieces"] = pieces
	tInterface["info"] = info

	encodedT, err := bencode.Marshal(tInterface)
	if err != nil {
		return nil, err
	}
	return encodedT, nil
}

func (tf *TorrentFile) GetRawData() map[string]interface{} {
	return tf.rawData
}

func (tf *TorrentFile) unmarshalTFile() error {
	if creationDate, ok := tf.rawData["creation date"].(int64); ok {
		tf.CreationDate = time.Unix(creationDate, 0)
	}

	tf.Encoding = castToString(tf.rawData["encoding"])
	tf.Comment = castToString(tf.rawData["comment"])
	tf.CreatedBy = castToString(tf.rawData["created by"])
	tf.Announce = castToString(tf.rawData["announce"])
	tf.AnnounceList = unmarshalAnnounceList(tf.rawData["announce-list"])
	switch ul := tf.rawData["url-list"].(type) {
	case []byte:
		tf.UrlList = []string{string(ul)}
	case []interface{}:
		for _, url := range ul {
			tf.UrlList = append(tf.UrlList, castToString(url))
		}
	}

	var err error
	tf.Info, err = unmarshalInfo(tf.rawData["info"])
	return err
}

func (tf *TorrentFile) GetInfoHash() (ih [20]byte, err error) {
	info, ok := tf.rawData["info"].(map[string]interface{})
	if !ok {
		return ih, errors.New("invalid torrent info field")
	}
	encInfo, err := bencode.Marshal(info)
	if err != nil {
		return ih, err
	}

	return sha1.Sum(encInfo), nil
}

func unmarshalInfo(rawInfo interface{}) (Info, error) {
	res := Info{}
	info, ok := rawInfo.(map[string]interface{})
	if !ok {
		return res, errors.New("error unmarshalling torrent info")
	}

	if plength, ok := info["piece length"].(int64); ok {
		res.PieceLength = plength
	}
	if length, ok := info["length"].(int64); ok {
		res.Length = length
	}
	if name, ok := info["name"].([]uint8); ok {
		res.Name = string(name)
	}
	if private, ok := info["private"].(int64); ok {
		res.Private = private != 0
	}
	if source, ok := info["source"].([]uint8); ok {
		res.Source = string(source)
	}

	if pbuffer, ok := info["pieces"].([]uint8); ok {
		if len(pbuffer)%20 != 0 {
			return res, errors.New("consistency check error: pieces length isn't multiple of 20")
		}
		piecesCount := len(pbuffer) / 20
		res.Pieces = make([]SHA1, piecesCount)
		for i := 0; i < len(res.Pieces); i++ {
			copy(res.Pieces[i][:], pbuffer[i*20:(i+1)*20])
		}
	}

	var err error
	if _, ok := info["files"]; ok {
		res.Files, err = unmarshalInfoFiles(info["files"])
	}

	return res, err
}

func unmarshalInfoFiles(rawInfoFiles interface{}) ([]MultiFileInfo, error) {
	var res []MultiFileInfo
	files, ok := rawInfoFiles.([]interface{})
	if !ok {
		return res, errors.New("corrupted torrent file list")
	}
	for _, file := range files {
		var f MultiFileInfo
		if fmap, ok := file.(map[string]interface{}); ok {
			if flength, ok := fmap["length"].(int64); ok {
				f.Length = flength
			}
			if md5sum, ok := fmap["md5sum"].([]uint8); ok {
				f.MD5sum = string(md5sum)
			}
			if ps, ok := fmap["path"].([]interface{}); ok {
				for _, p := range ps {
					if cp, ok := p.([]uint8); ok {
						f.Path = append(f.Path, string(cp))
					}
				}
			}
			res = append(res, f)
		}
	}
	return res, nil
}

func getFiles(path string) ([]MultiFileInfo, error) {
	var mf []MultiFileInfo
	if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			pathParts := strings.Split(path, string(filepath.Separator))
			info, err := d.Info()
			if err != nil {
				return err
			}
			fi := MultiFileInfo{
				Path:   pathParts,
				Length: info.Size(),
			}
			mf = append(mf, fi)

		}
		return nil
	}); err != nil {
		log.Fatalf("impossible to walk directories: %s", err)
	}
	fmt.Println(mf)
	return mf, nil
}

func calculateParts(files []MultiFileInfo, pSize int) ([]byte, error) {
	pieceSize := pSize * 1024
	buffer := make([]byte, pieceSize)
	var readBytes int
	var pieces []byte
	for _, file := range files {
		f, err := os.Open(strings.Join(file.Path, "/"))
		if err != nil {
			return nil, err
		}
		fs, err := f.Stat()
		if err != nil {
			return nil, err
		}
		fmt.Println(fs.Size(), pieceSize)

		for {
			n, err := f.Read(buffer[readBytes:])
			if err == io.EOF {
				break
			}
			readBytes += n
			if readBytes == pieceSize {
				h := sha1.Sum(buffer)
				pieces = append(pieces, h[:]...)
				readBytes = 0
			}
		}
	}
	h := sha1.Sum(buffer[:readBytes])
	pieces = append(pieces, h[:]...)

	return pieces, nil
}

func unmarshalAnnounceList(announceList interface{}) [][]string {
	var res [][]string
	if anlist, ok := announceList.([]interface{}); ok {
		for _, tier := range anlist {
			if trackerlist, ok := tier.([]interface{}); ok {
				strTrackerList := []string{}
				for _, tracker := range trackerlist {
					if tr, ok := tracker.([]uint8); ok {
						strTrackerList = append(strTrackerList, string(tr))
					}
				}
				res = append(res, strTrackerList)
			}
		}
		return res
	}
	return nil
}

func castToString(buffer interface{}) string {
	if data, ok := buffer.([]uint8); ok {
		return string(data)
	}
	return ""
}
