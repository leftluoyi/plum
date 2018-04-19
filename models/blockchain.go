package models

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Content
}

type Content struct {
	Affiliation		string
}

