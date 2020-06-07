package utils

// BSD 3-Clause License Copyright (c) 2020

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// Bootstrap decides whether or not to start the GUI.
type Bootstrap struct {
	ExePy      string
	TempPath   string
	IsPython3  bool
	HaveImgs   bool
	HaveGUI    bool
	HaveConfig bool
}

// WriteFile writes a plain text file.
func (btsrp Bootstrap) WriteFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0700)
}

// ReadJSON reads a serialized JSON file.
func (btsrp Bootstrap) ReadJSON(path string) map[string]interface{} {
	fileContent, _ := ioutil.ReadFile(path)
	var fileMap map[string]interface{}
	json.Unmarshal(fileContent, &fileMap)
	return fileMap
}

// ReadFile reads a file and returns bytes.
func (btsrp Bootstrap) ReadFile(path string) []byte {
	fileContent, _ := ioutil.ReadFile(path)
	return fileContent
}

// CreateProjectFile creates files needed for application runtime.
func (btsrp *Bootstrap) CreateProjectFile(basename string, data []byte, out chan<- struct{}) {
	err := btsrp.WriteFile(btsrp.TempPath+basename, data)
	if basename == "gui.py" && err != nil {
		btsrp.HaveGUI = false
	}
	if basename == "project.json" && err != nil {
		btsrp.HaveConfig = false
	}
	out <- struct{}{}
}

// GetImgs pulls the window icons from the web for life of process only.
func (btsrp *Bootstrap) GetImgs(url string, out chan<- struct{}) {
	response, err := http.Get(url)
	if err != nil {
		btsrp.HaveImgs = false
	} else {
		defer response.Body.Close()
		file, err := os.Create(btsrp.TempPath + url[len(url)-8:])
		if err != nil {
			btsrp.HaveImgs = false
		}
		defer file.Close()
		if btsrp.HaveImgs {
			_, err = io.Copy(file, response.Body)
			if err != nil {
				btsrp.HaveImgs = false
			}
		}
	}
	out <- struct{}{}
}

// GetInitialJSON returns the projects default JSON settings.
func (btsrp Bootstrap) GetInitialJSON() []byte {
	initJSON := []byte(`{"APPCOLOR": {"appcolor": "black"}, "TITLE": {"title": `)
	initJSON = append(initJSON, []byte(`"MyApp"}, "DIMENSIONS": {"dimensions": `)...)
	initJSON = append(initJSON, []byte(`"300x400"}, "MENUCOLOR": {"background"`)...)
	initJSON = append(initJSON, []byte(`: "#666666", "foreground": "#d9d9d9"}, `)...)
	return append(initJSON, `"THEME": {"theme":"default"}}`...)
}

// CheckPython ensures Python3 is installed and in path.
func (btsrp *Bootstrap) CheckPython(out chan<- struct{}) {
	for index, exe := range [4]string{"python", "python", "python3", "python3"} {
		var output bytes.Buffer
		pyCheck := exec.Command(exe, "--version")
		if index%2 != 0 {
			pyCheck.Stdout = &output
		} else {
			pyCheck.Stderr = &output
		}
		err := pyCheck.Run()
		if err == nil && bytes.Contains(output.Bytes(), []byte("Python 3")) {
			btsrp.IsPython3 = true
			btsrp.ExePy = exe
			break
		}
	}
	out <- struct{}{}
}

// ErrorExit exits and leaves log in cwd if initial startup errors occur.
func (btsrp Bootstrap) ErrorExit(warn string) {
	logfile, _ := os.OpenFile("VISIPY-ERROR-LOG", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer logfile.Close()
	log.SetOutput(logfile)
	log.Println(warn)

	if !btsrp.IsPython3 {
		log.Println("Failed to find python3 installation.")
	}
	if !btsrp.HaveGUI {
		log.Println("Failed to write Python GUI.")
	}
	if !btsrp.HaveImgs {
		log.Println("Failed to pull application images from the web.")
	}
	if !btsrp.HaveConfig {
		log.Println("Failed to write project JSON config.")
	}
	log.Fatalf("Exiting...")
}

// MasterLightOffChecklist examines bootstrap status before opening app.
func (btsrp *Bootstrap) MasterLightOffChecklist() {
	if !btsrp.IsPython3 || !btsrp.HaveImgs || !btsrp.HaveGUI || !btsrp.HaveConfig {
		btsrp.ErrorExit("Unable to start application.")
	}
}
