package torrentfile

import (
	"testing"
)

func TestNewMagnetURI_AllFields(t *testing.T) {
	magnet := &MagnetLink{
		InfoHash:          "abcdef1234567890",
		DisplayName:       "TestName",
		Trackers:          []string{"tracker.com/announce", "tracker2"},
		Length:            42,
		Keywords:          []string{"key1,key2"},
		AcceptableSources: []string{"source1", "source2"},
		ExactSources:      []string{"exact1", "exact2"},
		ManifestTopic:     []string{"topic1", "topic2"},
	}
	uri := NewMagnetURI(magnet)
	expected := "magnet:?xt=urn:btih:abcdef1234567890" +
		"&dn=TestName" +
		"&tr=tracker.com%2Fannounce" +
		"&tr=tracker2" +
		"&xl=42" +
		"&kt=key1%2Ckey2" +
		"&as=source1" +
		"&as=source2" +
		"&xs=exact1" +
		"&xs=exact2" +
		"&mt=topic1" +
		"&mt=topic2"
	if uri != expected {
		t.Errorf("unexpected URI:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestNewMagnetURI_Minimal(t *testing.T) {
	magnet := &MagnetLink{
		InfoHash: "abcdef1234567890",
	}
	uri := NewMagnetURI(magnet)
	expected := "magnet:?xt=urn:btih:abcdef1234567890"
	if uri != expected {
		t.Errorf("unexpected URI:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestNewMagnetURI_EmptyFields(t *testing.T) {
	magnet := &MagnetLink{}
	uri := NewMagnetURI(magnet)
	expected := "magnet:?xt=urn:btih:"
	if uri != expected {
		t.Errorf("unexpected URI:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestNewMagnetURI_Escaping(t *testing.T) {
	magnet := &MagnetLink{
		InfoHash:    "abcdef1234567890",
		DisplayName: "Name With Spaces",
		Trackers:    []string{"tracker.com/announce?foo=bar&baz=qux"},
	}
	uri := NewMagnetURI(magnet)
	expected := "magnet:?xt=urn:btih:abcdef1234567890" +
		"&dn=Name+With+Spaces" +
		"&tr=tracker.com%2Fannounce%3Ffoo%3Dbar%26baz%3Dqux"
	if uri != expected {
		t.Errorf("unexpected URI:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestParseMagnetURI_AllFields(t *testing.T) {
	uri := "magnet:?xt=urn:btih:abcdef1234567890&dn=TestName&tr=tracker.com%2Fannounce&tr=tracker2&xl=42&kt=key1,key2&as=source1&as=source2&xs=exact1&xs=exact2&mt=topic1&mt=topic2"
	magnet, err := ParseMagnetURI(uri)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if magnet.InfoHash != "abcdef1234567890" {
		t.Errorf("unexpected InfoHash: %s", magnet.InfoHash)
	}
	if magnet.DisplayName != "TestName" {
		t.Errorf("unexpected DisplayName: %s", magnet.DisplayName)
	}
	if len(magnet.Trackers) != 2 || magnet.Trackers[0] != "tracker.com/announce" || magnet.Trackers[1] != "tracker2" {
		t.Errorf("unexpected Trackers: %v", magnet.Trackers)
	}
	if magnet.Length != 42 {
		t.Errorf("unexpected Length: %v", magnet.Length)
	}
	if len(magnet.Keywords) != 1 || magnet.Keywords[0] != "key1,key2" {
		t.Errorf("unexpected Keywords: %v", magnet.Keywords)
	}
	if len(magnet.AcceptableSources) != 2 || magnet.AcceptableSources[0] != "source1" || magnet.AcceptableSources[1] != "source2" {
		t.Errorf("unexpected AcceptableSources: %v", magnet.AcceptableSources)
	}
	if len(magnet.ExactSources) != 2 || magnet.ExactSources[0] != "exact1" || magnet.ExactSources[1] != "exact2" {
		t.Errorf("unexpected ExactSources: %v", magnet.ExactSources)
	}
	if len(magnet.ManifestTopic) != 2 || magnet.ManifestTopic[0] != "topic1" || magnet.ManifestTopic[1] != "topic2" {
		t.Errorf("unexpected ManifestTopic: %v", magnet.ManifestTopic)
	}
}

func TestParseMagnetURI_InvalidPrefix(t *testing.T) {
	uri := "invalid:?xt=urn:btih:abcdef1234567890"
	_, err := ParseMagnetURI(uri)
	if err == nil {
		t.Error("expected error for invalid prefix")
	}
}

func TestParseMagnetURI_InvalidXL(t *testing.T) {
	uri := "magnet:?xt=urn:btih:abcdef1234567890&xl=notanumber"
	_, err := ParseMagnetURI(uri)
	if err == nil {
		t.Error("expected error for invalid xl parameter")
	}
}

func TestParseMagnetURI_InvalidXT(t *testing.T) {
	uri := "magnet:?xt=%ZZ"
	_, err := ParseMagnetURI(uri)
	if err == nil {
		t.Error("expected error for invalid xt parameter")
	}
}

func TestParseMagnetURI_EmptyParams(t *testing.T) {
	uri := "magnet:?"
	magnet, err := ParseMagnetURI(uri)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if magnet.InfoHash != "" || magnet.DisplayName != "" || len(magnet.Trackers) != 0 || magnet.Length != 0 {
		t.Errorf("unexpected non-empty MagnetLink: %+v", magnet)
	}
}
