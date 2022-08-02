package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//User takes the environment supplied during the run of the function and produces key-val pairs
type ManipulatorFunc func(env Env) (string, []KeyValPair)
type FileType string
type Key string      //This is a dot seperated path to the key in the json, eg: student.profile.name
type Val interface{} //This is the value to be put on the key
type KeyValPair struct {
	Key
	Val
}

const (
	JSON FileType = "JSON"
	YAML FileType = "YAML"
)

type Env struct {
	ParentDirectoryName string
}
type Copier struct {
	FilePath string   `json:"filepath"`
	RootDir  string   `json:"rootdir"`
	FileType FileType `json:"filetype"`
	errors   []error
	mx       sync.Mutex
}

func NewCopier(fp string, ft FileType, rd string) *Copier {
	return &Copier{
		FilePath: fp,
		FileType: ft,
		RootDir:  rd,
		errors:   make([]error, 0),
	}
}

func (c *Copier) Copy(man ManipulatorFunc) error {
	c.mx.Lock() //One instance of Copier can have only once instance of Copy function running
	defer c.mx.Unlock()
	if man == nil {
		return fmt.Errorf("Nil manipulator passed. Pass a valid function")
	}
	file, err := os.Open(c.FilePath)
	if err != nil {
		return fmt.Errorf("Could not open file: %s", err.Error())
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Could not read file: %s", err.Error())
	}
	fi, err := ioutil.ReadDir(c.RootDir)
	if err != nil {
		return fmt.Errorf("Could not read root directory: %s", err.Error())
	}
	for _, file := range fi {
		if !file.IsDir() || file.Name()[0] == '.' {
			continue
		}
		env := Env{ParentDirectoryName: file.Name()}
		newFilename, KeyValPairs := man(env)
		x := make(map[string]interface{}, 1)
		switch c.FileType {
		case JSON:
			err := json.Unmarshal(content, &x)
			if err != nil {
				c.errors = append(c.errors, err)
				continue
			}
		case YAML:
			return fmt.Errorf("Currently yaml not supported")
		}

		for _, v := range KeyValPairs {
			x, err = change(v.Key, v.Val, x)
			if err != nil {
				c.errors = append(c.errors, fmt.Errorf("Could not apply key-val: %s %s due to %s", v.Key, v.Val, err.Error()))
				continue
			}
		}
		switch c.FileType {
		case JSON:
			content, err = json.MarshalIndent(x, "", "\t")
			if err != nil {
				c.errors = append(c.errors, err)
				continue
			}
		case YAML:
			return fmt.Errorf("Currently yaml not supported")
		}
		err := writeToDirectory(content, filepath.Join(c.RootDir, file.Name()), newFilename)
		if err != nil {
			c.errors = append(c.errors, err)
			continue
		}
	}

	return mergeErrors(c.errors)
}

func writeToDirectory(content []byte, dirPath string, newfilename string) error {
	_, err := os.Create(filepath.Join(dirPath, newfilename))
	if err != nil {
		fmt.Println("err here", err.Error())
		return err
	}
	err = os.WriteFile(filepath.Join(dirPath, newfilename), content, 0777)
	if err != nil {
		return err
	}
	return nil
}
func mergeErrors(err []error) error {
	var newerr string
	for _, er := range err {
		newerr += "\n" + er.Error()
	}
	if newerr == "" {
		return nil
	}
	return fmt.Errorf(newerr)
}

//currenly this only supports while traversing through nested objects
func change(key Key, val Val, x map[string]interface{}) (changedContent map[string]interface{}, err error) {
	sep := strings.Split(string(key), ".")
	if len(sep) == 0 { //base condition
		return nil, fmt.Errorf("Could not find the key")
	}
	if len(sep) == 1 { //base condition
		x[sep[0]] = val
		return x, nil
	}
	newKey := sep[1:]
	key = Key(sep[0])
	if x[string(key)] == nil {
		return nil, fmt.Errorf("Cannot find key: %s", key)
	}
	newx, ok := x[string(key)].(map[string]interface{}) //currenly this only supports while traversing through nested objects
	if !ok {
		return nil, fmt.Errorf("Could not convert value of ")
	}
	x[string(key)], err = change(Key(strings.Join(newKey, ".")), val, newx)
	if err != nil {
		return nil, err
	}
	return x, nil
}
