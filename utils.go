package packetcoder

import "encoding/json"

type JSONPacket struct {
	Name   string
	Fields []JSONField
}

type JSONField struct {
	Name         string
	Size         uint
	LittleEndian bool
}

func ReadSchemeFromString(dataString string) *BitScheme {
	var packet JSONPacket
	json.Unmarshal([]byte(dataString), &packet)
	var currentScheme *BitScheme
	currentScheme = NewBitScheme(packet.Name)
	for j := 0; j < len(packet.Fields); j++ {
		field := packet.Fields[j]
		if field.LittleEndian {
			currentScheme.AddBitFieldLittleEndian(field.Name, field.Size)
		} else {
			currentScheme.AddBitField(field.Name, field.Size)
		}
	}
	return currentScheme
}

func ReadSchemesFromString(dataString string) map[string]*BitScheme {
	var mapOfSchemes map[string]*BitScheme
	mapOfSchemes = make(map[string]*BitScheme)
	var packets []JSONPacket
	json.Unmarshal([]byte(dataString), &packets)
	var currentScheme *BitScheme
	for i := 0; i < len(packets); i++ {
		currentScheme = NewBitScheme(packets[i].Name)
		for j := 0; j < len(packets[i].Fields); j++ {
			field := packets[i].Fields[j]
			if field.LittleEndian {
				currentScheme.AddBitFieldLittleEndian(field.Name, field.Size)
			} else {
				currentScheme.AddBitField(field.Name, field.Size)
			}
		}
		mapOfSchemes[packets[i].Name] = currentScheme
	}
	return mapOfSchemes
}
