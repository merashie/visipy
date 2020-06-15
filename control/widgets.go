package control

// BSD 3-Clause License Copyright (c) 2020
// v0.2

import (
	"bytes"
	"fmt"
	"strings"
)

func (cont AppController) getAppColor() string {
	color := fmt.Sprintf("%[1]s# App Color\n%[1]s", cont.I2)
	return color + "self.master.configure(bg='{{.appcolor}}')\n\n"
}

func (cont AppController) getAppTitle() string {
	title := fmt.Sprintf("%[1]s# App Title\n%[1]s", cont.I2)
	return title + "self.master.title('{{.title}}')\n\n"
}

func (cont AppController) getAppDimensions() string {
	dim := fmt.Sprintf("%[1]s# Overall Dimensions\n%[1]sself.master", cont.I2)
	return dim + ".geometry('{{.dimensions}}')\n\n"
}

func (cont AppController) getMenuInit() string {
	menuInit := []byte(
		`		# Window Menu Color
		menu = Menu(self.master)
		menu.config(foreground='{{.foreground}}', background='{{.background}}')
		self.master.config(menu=menu)

`)
	return string(bytes.ReplaceAll(menuInit, []byte{0x09}, cont.I1b))
}

func (cont AppController) getMenuItem(title string, numberSubmenus int8) string {
	var menu bytes.Buffer
	var titleLower = strings.ToLower(title)
	menu.WriteString(fmt.Sprintf("%s# %s\n", cont.I2, title))
	menu.WriteString(fmt.Sprintf("%s%s_menu = Menu(menu)\n", cont.I2, titleLower))

	var index int8
	for index = 0; index < numberSubmenus; index++ {
		menu.WriteString(fmt.Sprintf("%s%s_menu.add_command(\n", cont.I2, titleLower))
		menu.WriteString(fmt.Sprintf("%s%slabel='{{.submenu%d}}',\n", cont.I1, cont.I2, index))
		menu.WriteString(fmt.Sprintf("%s%scommand=quit_\n%s)\n\n", cont.I1, cont.I2, cont.I2))
	}

	menu.WriteString(fmt.Sprintf("%smenu.add_cascade(label='%s', ", cont.I2, title))
	menu.WriteString(fmt.Sprintf("menu=%s_menu)\n\n", titleLower))
	return menu.String()
}

func (cont AppController) getMethods(methodNames []string) []byte {
	todo := "TODO: Add handling code here"
	var collection bytes.Buffer
	for _, name := range methodNames {
		collection.WriteString(fmt.Sprintf("%sdef %s(self):\n", cont.I1, name))
		collection.WriteString(fmt.Sprintf("%s\"\"\" %s \"\"\"\n", cont.I2, todo))
		msg := fmt.Sprintf("%sprint('Handle %s here')\n\n", cont.I2, name)
		collection.WriteString(msg)
	}
	return collection.Bytes()
}

func (cont AppController) getAppIcon() string {
	icon := []byte(
		`		# ICON
		icon_path = '{{.iconpath}}'
		self.icon = PhotoImage(file=icon_path)
		master.iconphoto(False, self.icon)


`)
	return string(bytes.ReplaceAll(icon, []byte{0x09}, cont.I1b))
}

func (cont AppController) getImports() []byte {
	return []byte(
		`from tkinter import *


`)
}

func (cont AppController) getSysImport() []byte {
	return []byte(
		`from sys import exit
`)
}

func (cont AppController) getStyleImport() []byte {
	return []byte(
		`from tkinter.ttk import Style
`)
}

func (cont AppController) getClassInit() string {
	classInit := []byte(
		`class {{.title}}:
	def __init__(self, master):
		self.master = master

`)
	return string(bytes.ReplaceAll(classInit, []byte{0x09}, cont.I1b))
}

func (cont AppController) getQuit() []byte {
	qu := []byte(
		`
def quit_():
	exit()

`)
	return bytes.ReplaceAll(qu, []byte{0x09}, cont.I1b)
}

func (cont AppController) getGui() string {
	tk := []byte(
		`
# App Theme
def run_gui():
	root = Tk()
	root.style = Style()
	root.style.theme_use('{{.THEME.theme}}')
	{{.TITLE.title}}(root)
	root.mainloop()


`)
	return string(bytes.ReplaceAll(tk, []byte{0x09}, cont.I1b))
}

func (cont AppController) getMain() []byte {
	main := []byte(
		`if __name__ == '__main__':
	run_gui()
`)
	return bytes.ReplaceAll(main, []byte{0x09}, cont.I1b)
}

func (cont AppController) getWidget() string {
	anonWidget := []byte(
		`		# {{.name}}
		self.{{.name}} = {{.widget}}(
			master,
			foreground='{{.foreground}}',
			background='{{.background}}',
			font='{{.font}}',
			text='{{.text}}',
			activeforeground='{{.activeforeground}}',
			activebackground='{{.activebackground}}',
			anchor={{.anchor}},
			command=self.{{.command}},
			highlightcolor='{{.highlightcolor}}',
			indicatoron={{.indicatoron}},
			selectcolor='{{.selectcolor}}',
			show='{{.show}}',
			insertbackground='{{.insertbackground}}',
			selectforeground='{{.selectforeground}}',
			selectbackground='{{.selectbackground}}',
			relief={{.relief}},
			justify={{.justify}},
			troughcolor='{{.troughcolor}}',
			orient={{.orient}},
			sliderrelief={{.sliderrelief}},
			values=({{.values}}),
			wrap={{.wrap}},
			selectmode={{.selectmode}},
			activestyle={{.activestyle}},
			length={{.length}},
			from_={{.from}},
			to={{.to}},
			tickinterval={{.tickinterval}},
			sliderlength={{.sliderlength}},
			width={{.width}},
			height={{.height}},
			borderwidth={{.borderwidth}},
			highlightthickness={{.highlightthickness}},
			selectborderwidth={{.selectborderwidth}},
			insertontime={{.insertontime}},
			insertofftime={{.insertofftime}},
		)

		self.{{.name}}.grid(
			row={{.row}},
			rowspan={{.rowspan}},
			column={{.column}},
			columnspan={{.columnspan}},
			padx={{.padx}},
			pady={{.pady}},
			sticky={{.sticky}},
		)`)

	return string(bytes.ReplaceAll(anonWidget, []byte{0x09}, cont.I1b))
}

func (cont AppController) getImgWidget() string {
	anonWidget := []byte(
		`		# {{.name}}
		self.{{.name}} = Label(master)
        {{.name}}_gif = '{{.image}}'
        self.{{.name}}.img = PhotoImage(file={{.name}}_gif)
        self.{{.name}}.config(
			background='{{.background}}',
			borderwidth={{.borderwidth}},
			image=self.{{.name}}.img
		)

		self.{{.name}}.grid(
			row={{.row}},
			rowspan={{.rowspan}},
			column={{.column}},
			columnspan={{.columnspan}},
			padx={{.padx}},
			pady={{.pady}},
			sticky={{.sticky}},
		)`)

	return string(bytes.ReplaceAll(anonWidget, []byte{0x09}, cont.I1b))
}

// SetWidget sets any TK widget into the current build.
func (app *AppParser) SetWidget(widgetType string, update []string) {

	tmp := make(map[string]interface{})
	for _, attr := range update {
		kv := strings.Split(attr, "|@|")
		tmp[kv[0]] = kv[1]
	}

	app.MapBuild[tmp["name"].(string)] = make(map[string]interface{})
	app.MapBuild[tmp["name"].(string)] = map[string]interface{}{
		"foreground":       "",
		"background":       "",
		"font":             "",
		"text":             "",
		"activeforeground": "",
		"activebackground": "",
		"anchor":           "",
		"command":          "",
		"highlightcolor":   "",
		"indicatoron":      "",
		"selectcolor":      "",
		"show":             "",
		"insertbackground": "",
		"selectforeground": "",
		"selectbackground": "",
		"relief":           "",
		"justify":          "",
		"image":            "",
		"troughcolor":      "",
		"orient":           "",
		"sliderrelief":     "",
		"values":           "",
		"wrap":             "",
		"selectmode":       "",
		"activestyle":      "",
		"sticky":           "",
		"widget":           widgetType,
		"name":             tmp["name"].(string),

		"length":             -1,
		"from":               -1,
		"to":                 -1,
		"tickinterval":       -1,
		"sliderlength":       -1,
		"width":              -1,
		"height":             -1,
		"borderwidth":        -1,
		"highlightthickness": -1,
		"selectborderwidth":  -1,
		"insertontime":       -1,
		"insertofftime":      -1,
		"row":                -1,
		"rowspan":            -1,
		"column":             -1,
		"columnspan":         -1,
		"padx":               -1,
		"pady":               -1,
	}

	for attr, val := range tmp {
		app.MapBuild[tmp["name"].(string)][attr] = val
	}
}

// ReviseWidget removes un-templated lines from the build buffer.
func (app *AppParser) ReviseWidget(tmpbuff bytes.Buffer) {

	end := []byte(",\n" + app.I2 + ")")

	widgetUpdate := bytes.Split(tmpbuff.Bytes(), []byte("\n"))
	tmpbuff.Reset()

	for _, line := range widgetUpdate {
		if !bytes.Contains(line, []byte("<no value>")) {
			tmpbuff.Write(line)
			tmpbuff.Write([]byte("\n"))
		}
	}

	tmpbuff.Write([]byte("\n"))
	app.Build.Write(bytes.ReplaceAll(tmpbuff.Bytes(), end, end[1:]))
}
