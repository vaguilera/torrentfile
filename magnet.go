package torrentfile

import (
	"fmt"
	"net/url"
	"strconv"
)

type MagnetLink struct {
	InfoHash          string
	Trackers          []string
	DisplayName       string
	Length            int64
	Keywords          []string
	AcceptableSources []string
	ExactSources      []string
	ManifestTopic     []string
}

func ParseMagnetURI(uri string) (*MagnetLink, error) {
	if uri[0:8] != "magnet:?" {
		return nil, fmt.Errorf("invalid magnet URI")
	}

	params, err := url.ParseQuery(uri[8:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse magnet URI: %w", err)
	}

	magnet := &MagnetLink{}

	for key, values := range params {
		switch key {
		case "xt":
			ifParam := values[0]
			urn := ifParam[0:9]
			if urn != "urn:btih:" {
				return nil, fmt.Errorf("not a BitTorrent info hash: %s", urn)
			}
			magnet.InfoHash = ifParam[9:]
		case "dn":
			magnet.DisplayName = values[0]
		case "tr":
			magnet.Trackers = values
		case "xl":
			magnet.Length, err = strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid xl parameter: %w", err)
			}

		case "kt":
			magnet.Keywords = values
		case "as":
			magnet.AcceptableSources = values
		case "xs":
			magnet.ExactSources = values
		case "mt":
			magnet.ManifestTopic = values
		default:
		}
	}

	return magnet, nil
}

func NewMagnetURI(magnet *MagnetLink) string {
	uri := "magnet:?xt=urn:btih:" + magnet.InfoHash

	if magnet.DisplayName != "" {
		uri += "&dn=" + url.QueryEscape(magnet.DisplayName)
	}
	for _, tracker := range magnet.Trackers {
		uri += "&tr=" + url.QueryEscape(tracker)
	}
	if magnet.Length > 0 {
		uri += "&xl=" + strconv.FormatInt(magnet.Length, 10)
	}
	if len(magnet.Keywords) > 0 {
		uri += "&kt=" + url.QueryEscape(magnet.Keywords[0]) // Assuming single keyword for simplicity
	}
	for _, source := range magnet.AcceptableSources {
		uri += "&as=" + url.QueryEscape(source)
	}
	for _, source := range magnet.ExactSources {
		uri += "&xs=" + url.QueryEscape(source)
	}
	for _, topic := range magnet.ManifestTopic {
		uri += "&mt=" + url.QueryEscape(topic)
	}

	return uri
}
