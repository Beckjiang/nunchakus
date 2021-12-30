package nunchakus

import (
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

type FileProcess struct {
	filePath string

	loaded bool
	v int
}

func NewFileProcess(filePath string) *FileProcess {
	fp := &FileProcess{filePath: filePath}
	fp.init()

	return fp
}

func (fp *FileProcess) init() {
	go func() {
		t := time.NewTicker(time.Second)
		if fp.filePath == "" {
			fp.filePath = "./progress.storage"
		}
		lastIndex := fp.v
		for {
			select {
			case <-t.C:
				if lastIndex == fp.v {
					continue
				}

				if err := ioutil.WriteFile(fp.filePath, []byte(strconv.Itoa(fp.v)), 777); err != nil{
					log.Println(err.Error())
				}
			}
		}
	}()
}

func (fp *FileProcess) Save(i int) {
	fp.v = i
}

func (fp *FileProcess) Load() int {
	if !fp.loaded {
		if fp.filePath == "" {
			fp.filePath = "./progress.storage"
		}
		log.Println("process store file:", fp.filePath)

		content ,err :=ioutil.ReadFile(fp.filePath)
		if err != nil || len(content) == 0 {
			fp.v = 0
		} else {
			v, err := strconv.Atoi(string(content))
			if err == nil {
				fp.v = v
			} else {
				fp.v = 0
			}
		}

		fp.loaded = true
	}

	return fp.v
}
