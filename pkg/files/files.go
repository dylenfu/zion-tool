package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func ReadFile(filepath string) ([]byte, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("OpenFile %s error %s", filepath, err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(fmt.Sprintf("File %s close error %s", filepath, err))
		}
	}()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll %s error %s", filepath, err)
	}
	return data, nil
}

func FullPath(workspace, dir string, fileName string) string {
	if !path.IsAbs(workspace) {
		panic("path should be absolute path")
	}
	return path.Join(workspace, dir, fileName)
}

// ReadJsonFile read file and unmarshal to struct instance
func ReadJsonFile(path string, ptr interface{}) error {
	enc, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(enc, ptr)
}

// WriteJsonFile encode struct instance to bytes and persis in file
func WriteJsonFile(path string, ptr interface{}, indent bool) (err error) {
	var enc []byte

	if indent {
		enc, err = json.MarshalIndent(ptr, "", "    ")
	} else {
		enc, err = json.Marshal(ptr)
	}

	if err != nil {
		return
	}

	return ioutil.WriteFile(path, enc, os.ModePerm)
}
