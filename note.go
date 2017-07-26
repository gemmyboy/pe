package pe

import (
	"bytes"
	"encoding/gob"
	"log"
)

/*
	note.go
		by Marcus Shannon

	Used as a simple slimmable structure to communicate between internal systems
*/

//Note -: data structure for passing data around quickly
type Note struct {
	Flag uint32
	From uint32
	Data []interface{}
} //End Note

//Slim -: Trim note into an array
func Slim(n *Note) []byte {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(n)
	if err != nil {
		log.Println(err)

	}
	return buf.Bytes()
} //End Slim

//Fatten -: Convert byte array into a Note
func Fatten(b []byte) *Note {
	buf := bytes.NewBuffer(b)
	var n *Note

	dec := gob.NewDecoder(buf)
	err := dec.Decode(n)
	if err != nil {
		log.Println(err)

	}
	return n
} //End Fatten
