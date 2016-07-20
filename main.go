package main

import (
	"encoding/json"
	"hash/crc64"
	"io"
	"os"
)

func main() {

	loadstate()
	defer savestate()

	var obj json.RawMessage
	var uuidObj struct {
		ID string `json:"uuid"`
	}
	dec := json.NewDecoder(os.Stdin) //TODO: buffered reader?

	for {
		err := dec.Decode(&obj)
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}
		if err := json.Unmarshal(obj, &uuidObj); err != nil {
			panic(err)
		}
		if lastHashCheck(uuidObj.ID, obj) {
			if _, err := os.Stdout.Write(obj); err != nil {
				panic(err)
			}
			if _, err := os.Stdout.Write([]byte("\n")); err != nil {
				panic(err)
			}
		}
	}
}

var last map[string]uint64

var tab = crc64.MakeTable(crc64.ECMA)

func lastHashCheck(uuid string, data []byte) (changed bool) {
	cs := crc64.Checksum(data, tab)
	if last[uuid] != cs {
		changed = true
		last[uuid] = cs
	}
	return
}

func loadstate() {

	last = make(map[string]uint64)
	f, err := os.Open("changedstate.json")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&last)
}

func savestate() {
	f, err := os.Create("changedstate.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.Encode(last)
}
