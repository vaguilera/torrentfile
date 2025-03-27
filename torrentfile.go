package torrentfile

import (
	"errors"
	"fmt"
	"os"
	"time"

	bencode "github.com/IncSW/go-bencode"
)

type SHA1 [20]byte

type TorrentMultiFileInfo struct {
	Length int64
	Path   []string
	MD5sum string
}

type TorrentFileInfo struct {
	Pieces      []SHA1
	PieceLength int64
	Length      int64
	Name        string
	Files       []TorrentMultiFileInfo
}

type TorrentFile struct {
	Announce     string
	AnnounceList [][]string
	CreationDate time.Time
	Encoding     string
	Info         TorrentFileInfo
	Comment      string
	CreatedBy    string

	rawData map[string]interface{}
}

func LoadTorrentFile(filename string) (*TorrentFile, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	rawTFile, err := bencode.Unmarshal(f)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling torrent file: %v", err)
	}

	t := &TorrentFile{}
	var ok bool
	t.rawData, ok = rawTFile.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid torrent file")
	}

	if err = t.unmarshallTFile(); err != nil {
		fmt.Println("error unmarshalling torrent file:", err)
	}

	return t, nil
}

func (tf *TorrentFile) GetRawData() map[string]interface{} {
	return tf.rawData
}

func (tf *TorrentFile) unmarshallTFile() error {
	if creationDate, ok := tf.rawData["creation date"].(int64); ok {
		tf.CreationDate = time.Unix(creationDate, 0)
	}

	tf.Encoding = castToString(tf.rawData["encoding"])
	tf.Comment = castToString(tf.rawData["comment"])
	tf.CreatedBy = castToString(tf.rawData["created by"])
	tf.Announce = castToString(tf.rawData["announce"])
	tf.AnnounceList = unMarshallAnnounceList(tf.rawData["announce-list"])

	var err error
	tf.Info, err = unMarshallInfo(tf.rawData["info"])
	return err
}

func unMarshallInfo(rawInfo interface{}) (TorrentFileInfo, error) {
	res := TorrentFileInfo{}
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

	pieces := []SHA1{}
	if pbuffer, ok := info["pieces"].([]uint8); ok {
		if len(pbuffer)%20 != 0 {
			return res, errors.New("consistency check error: pieces length isn't multiple of 20")
		}
		for i := 0; i < len(pbuffer); i += 20 {
			var chunk [20]byte
			copy(chunk[:], pbuffer[i:i+20])
			pieces = append(pieces, chunk)
		}
	}

	var err error
	res.Files, err = unMarshallInfoFiles(info["files"])
	return res, err

}

func unMarshallInfoFiles(rawInfoFiles interface{}) ([]TorrentMultiFileInfo, error) {
	res := []TorrentMultiFileInfo{}
	files, ok := rawInfoFiles.([]interface{})
	if !ok {
		return res, errors.New("corrupted torrent file list")
	}
	for _, file := range files {
		var f TorrentMultiFileInfo
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

func unMarshallAnnounceList(announceList interface{}) [][]string {
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
