## PacketCoder
The main idea of this library is provide mechanism to create packets descriptions and generate packets based on these descriptions.

## Installation
For getting PacketCoder run on your machine, you just:
```go
go get github.com/SealNTibbers/packetcoder
```

## How to use

When you use this library you should choose one of two ways to set scheme:

1) Write scheme description manually
``` go
  packet := NewPacket()
	scheme := NewBitScheme("testPacket")
	scheme.AddBitField("head", 4)
	scheme.AddBitField("type", 8)
	scheme.AddStuffBits("fill", 4)
	scheme.AddBitField("crc", 8)
	packet.SetScheme(scheme)
```
2) Fill scheme from JSON string with description of packet
``` go
  var sampleJSON = `{
		"name": "testPacket",
		"fields":
                [
					{"name": "head", "size": 4},
					{"name": "type", "size": 8},
					{"name": "fill", "size": 4},
					{"name": "crc", "size": 8}
				]
	}`
  var scheme *BitScheme
	scheme = ReadSchemeFromString(sampleJSON)
  packet.SetScheme(scheme)
 ```
To fill packet with data you can use:
``` go
    WriteValue64(fieldName string, value uint64) // for writing single data field
    WriteBytes(fieldName string, value []byte) // for writing array data field
    WriteStuff(fieldName string) // for writing stub
```

To read packet from data you can use:
``` go
    ReadValue64(fieldName string) // for reading single data field
    ReadBytesValue(fieldName string) // for reading array data field
```
 
Now, you can work with packet. If you want to encode packet into bytes buffer you can use:
``` go
  buf := bytes.NewBuffer(nil)
	packet.EncodeTo(buf)
```
for decoding:

``` go
  buf := bytes.NewBuffer([]byte{86, 144, 99, 86, 144, 140, 86, 144})
	_, err := packet.DecodeFrom(buf)
```
