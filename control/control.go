package control

// BSD 3-Clause License Copyright (c) 2020
// v0.2

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"text/template"

	"github.com/rootVIII/visipy/utils"
)

// Controller implements methods for AppController.
type Controller interface {
	RunVisipy()
}

// AppController controls the entire GUI application.
type AppController struct {
	STDOUT   []string
	Build    bytes.Buffer
	MapBuild map[string]map[string]interface{}
	I1b      []byte
	I1       string
	I2       string
}

// AppParser inherits AppController for parsing output.
type AppParser struct {
	AppController
	Executable string
	VisiPath   string
	Project    string
	HaveIcon   bool
	Utils      utils.Bootstrap
}

// RunVisipy runs the application and listens to incoming commands.
func (app *AppParser) RunVisipy() {
	app.I1b = []byte{0x20, 0x20, 0x20, 0x20}
	app.I1 = string(app.I1b)
	app.I2 = app.I1 + app.I1

	// Set AppController's initial current builds and buffers.
	app.initUserApp()
	app.RunTemplate(true)

	command := exec.Command(app.Executable, app.VisiPath)
	stdout, _ := command.StdoutPipe()
	command.Start()
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {

		app.STDOUT = strings.Split(scanner.Text(), "|$|")

		switch app.STDOUT[0] {
		case "ADD":
			app.SetWidget(app.STDOUT[1], strings.Split(app.STDOUT[2], "|:|"))
		case "EXIT":
			break
		case "LOADUSERPROJ":
			app.loadExistingProject(app.STDOUT[1])
		case "RESET":
			app.initUserApp()
		case "BUILD":
			app.RunJob()
		case "REMOVE":
			delete(app.MapBuild, app.STDOUT[1])
			if app.STDOUT[1] == "ICON" {
				app.HaveIcon = false
			}
		case "WRITE":
			jsonBytes, _ := json.Marshal(app.MapBuild)
			if app.STDOUT[1][len(app.STDOUT[1])-3:] != ".py" {
				app.STDOUT[1] += ".py"
			}
			app.Utils.WriteFile(app.STDOUT[1]+".project", jsonBytes)
			app.Utils.WriteFile(app.STDOUT[1], app.Build.Bytes())
		case "APPCOLOR":
			app.MapBuild["APPCOLOR"]["appcolor"] = app.STDOUT[1]
		case "DIMENSIONS":
			app.MapBuild["DIMENSIONS"]["dimensions"] = app.STDOUT[1]
		case "ICON":
			app.MapBuild["ICON"] = make(map[string]interface{})
			app.MapBuild["ICON"]["iconpath"] = app.STDOUT[1]
			app.HaveIcon = true
		case "MENU":
			menuItems := strings.Split(app.STDOUT[1], ",")
			_, hasMenu := app.MapBuild[menuItems[0]]
			if !hasMenu {
				app.MapBuild[menuItems[0]] = make(map[string]interface{})
			}
			for index, value := range menuItems[1:] {
				app.MapBuild[menuItems[0]][fmt.Sprintf("submenu%d", index)] = value
			}
		case "THEME":
			app.MapBuild["THEME"]["theme"] = app.STDOUT[1]
		case "TITLE":
			app.MapBuild["TITLE"]["title"] = app.STDOUT[1]
		case "MENUCOLOR":
			colors := strings.Split(app.STDOUT[1], "|:|")
			app.MapBuild["MENUCOLOR"]["foreground"] = colors[0]
			app.MapBuild["MENUCOLOR"]["background"] = colors[1]
		}
		app.RunTemplate(false)
	}
}

func (app *AppParser) initUserApp() {
	app.MapBuild = make(map[string]map[string]interface{})
	app.MapBuild["TITLE"] = make(map[string]interface{})
	app.MapBuild["APPCOLOR"] = make(map[string]interface{})
	app.MapBuild["DIMENSIONS"] = make(map[string]interface{})
	app.MapBuild["MENUCOLOR"] = make(map[string]interface{})
	app.MapBuild["THEME"] = make(map[string]interface{})
	app.MapBuild["TITLE"]["title"] = "MyApp"
	app.MapBuild["MENUCOLOR"]["foreground"] = "#d9d9d9"
	app.MapBuild["MENUCOLOR"]["background"] = "#666666"
	app.MapBuild["DIMENSIONS"]["dimensions"] = "300x400"
	app.MapBuild["APPCOLOR"]["appcolor"] = "#000000"
	app.MapBuild["THEME"]["theme"] = "default"
	app.HaveIcon = false
}

// RunTemplate templates map values into code snippets.
func (app *AppParser) RunTemplate(initialBuild bool) {
	app.Build.Reset()
	app.Build.Write(app.getSysImport())
	app.Build.Write(app.getStyleImport())
	app.Build.Write(app.getImports())

	out, _ := template.New("classinit").Parse(app.getClassInit())
	out.Execute(&app.Build, app.MapBuild["TITLE"])
	out, _ = template.New("apptitle").Parse(app.getAppTitle())
	out.Execute(&app.Build, app.MapBuild["TITLE"])
	out, _ = template.New("appcolor").Parse(app.getAppColor())
	out.Execute(&app.Build, app.MapBuild["APPCOLOR"])
	out, _ = template.New("dimensions").Parse(app.getAppDimensions())
	out.Execute(&app.Build, app.MapBuild["DIMENSIONS"])
	out, _ = template.New("menucolor").Parse(app.getMenuInit())
	out.Execute(&app.Build, app.MapBuild["MENUCOLOR"])

	if app.HaveIcon {
		out, _ = template.New("icon").Parse(app.getAppIcon())
		out.Execute(&app.Build, app.MapBuild["ICON"])
	}
	for key, value := range app.MapBuild {
		_, hasSub := value["submenu0"]
		if hasSub {
			out, _ = template.New(key).Parse(app.getMenuItem(key, int8(len(value))))
			out.Execute(&app.Build, app.MapBuild[key])
		}
	}

	var methods []string
	for key, value := range app.MapBuild {
		_, isWidget := value["row"]
		if !isWidget {
			continue
		}
		for k, v := range value {
			ival, isint := v.(int)
			if isint && ival == -1 {
				delete(value, k)
			}
			sval, isstring := v.(string)
			if isstring && len(sval) < 1 {
				delete(value, k)
			}
		}

		var tmpbuf bytes.Buffer
		_, isImage := value["image"]
		if isImage {
			tmp, _ := template.New(key).Parse(app.getImgWidget())
			tmp.Execute(&tmpbuf, value)
			app.ReviseWidget(tmpbuf)
			continue
		}

		methodName, hasCommand := value["command"]
		if hasCommand && len(methodName.(string)) > 0 {
			methods = append(methods, methodName.(string))
		}
		tmp, _ := template.New(key).Parse(app.getWidget())
		tmp.Execute(&tmpbuf, value)
		app.ReviseWidget(tmpbuf)
	}

	if len(methods) > 0 {
		app.Build.Write(app.getMethods(methods))
	}
	app.Build.Write(app.getQuit())
	out, _ = template.New("theme").Parse(app.getGui())
	out.Execute(&app.Build, app.MapBuild)
	app.Build.Write(app.getMain())
	rawProject, _ := json.Marshal(app.MapBuild)

	if initialBuild {
		app.Utils.WriteFile(fmt.Sprintf("%s.py", app.Project), app.Build.Bytes())
	} else {
		app.Utils.WriteFile(fmt.Sprintf("%s.json.update", app.Project), rawProject)
		app.Utils.WriteFile(fmt.Sprintf("%s.py.update", app.Project), app.Build.Bytes())
	}
}

func (app *AppParser) loadExistingProject(projectPath string) {
	for key, value := range app.Utils.ReadJSON(projectPath) {
		if key == "ICON" {
			app.HaveIcon = true
		}
		_, exists := app.MapBuild[key]
		if !exists {
			app.MapBuild[key] = make(map[string]interface{})
		}
		for innerKey, innerValue := range value.(map[string]interface{}) {
			app.MapBuild[key][innerKey] = innerValue
		}
	}
}

// RunJob runs the user's current Python app.
func (app *AppParser) RunJob() {
	cmd := exec.Command(app.Executable, fmt.Sprintf("%s.py", app.Project))
	_ = cmd.Start()
}
