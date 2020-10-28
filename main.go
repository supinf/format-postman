package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	Prefix        = ".json"
	Info          = "info"
	Item          = "item"
	Event         = "event"
	Script        = "script"
	ScriptID      = "id"
	PostmanID     = "_postman_id"
	CollectionDir = "POSTMAN_COLLECTION_DIR"
)

func main() {
	fileDir := os.Getenv(CollectionDir)
	if fileDir == "" {
		err := fmt.Errorf(`ERROR: no environment variables are set "%s"`, CollectionDir)
		panic(err)
	}

	files, err := ioutil.ReadDir(fileDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if ex := filepath.Ext(file.Name()); ex == Prefix {
			if err := formatPostManFile(filepath.Join(fileDir, file.Name())); err != nil {
				panic(err)
			}
		}
	}
}

func formatPostManFile(filePath string) error {

	// read json file
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// decode to json
	var decodeData map[string]interface{}
	if err := json.Unmarshal(bytes, &decodeData); err != nil {
		return err
	}

	// delete _postman_id
	if info, ok := decodeData[Info]; ok {
		delete(info.(map[string]interface{}), PostmanID)
	}

	// delete script id
	deleteScriptIDFromItem(decodeData)

	// convert json to byte
	b, err := json.MarshalIndent(decodeData, "", "\t")
	if err != nil {
		return err
	}

	// write json file
	fp, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fp.Close()

	writer := bufio.NewWriter(fp)
	if _, err := writer.Write(b); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}

func deleteScriptIDFromItem(data map[string]interface{}) {

	if event, ok := data[Event]; ok {
		for _, e := range event.([]interface{}) {
			if script, ok := e.(map[string]interface{})[Script]; ok {
				delete(script.(map[string]interface{}), ScriptID)
			}
		}
	}

	if item, ok := data[Item]; ok {
		for _, i := range item.([]interface{}) {
			deleteScriptIDFromItem(i.(map[string]interface{}))
		}
	}
}
