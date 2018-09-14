package packetcoder

import "encoding/json"

type JSONPacket struct {
	Name   string
	Fields []JSONField
}

type JSONField struct {
	Name string
	Size uint
}

func ReadSchemesFromString(dataString string) []*BitScheme {
	var listOfSchemes []*BitScheme
	var packets []JSONPacket
	json.Unmarshal([]byte(dataString), &packets)
	var currentScheme *BitScheme
	for i := 0; i < len(packets); i++ {
		currentScheme = NewBitScheme()
		for j := 0; j < len(packets[i].Fields); j++ {
			currentScheme.SetBitField(packets[i].Fields[j].Name, packets[i].Fields[j].Size)
		}
		listOfSchemes = append(listOfSchemes, currentScheme)
	}
	return listOfSchemes
}
