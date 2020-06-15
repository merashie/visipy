package utils

import "bytes"

// BSD 3-Clause License Copyright (c) 2020
// v0.2

// GetGui gets the app's Python/Tkinter GUI front end.
// Below code is Pycodestyle valid once written to disk
// by the application. It's removed after execution.
func GetGui() []byte {

	vpy := []byte(`
import builtins
from json import load
from keyword import iskeyword
from os import rename, environ
from os.path import realpath, basename, isfile
from sys import stdout, modules
from threading import Thread
from time import sleep
from tkinter import Tk, Menu, Label, Spinbox, Entry, LEFT, CENTER
from tkinter import Button, Listbox, Text, FLAT, SUNKEN
from tkinter import PhotoImage, Scrollbar, Scale, Toplevel
from tkinter import E, W, END, HORIZONTAL, NORMAL, DISABLED
from tkinter.messagebox import askyesno
from tkinter.filedialog import askopenfilename
from tkinter.filedialog import asksaveasfilename
from tkinter.font import names
from tkinter.ttk import Style


class VisiPy:
	def __init__(self, master):
		self.dark = '#404040'
		self.light = '#CCCCCC'
		self.bold = 'TkHeadingFont 10'
		self.normal = 'TkTextFont 8'
		self.small = 'TkSmallCaptionFont 8'
		rpath = realpath(__file__)[:-len(basename(__file__))]
		self.data_path = '%sproject.json' % rpath
		self.code_path = '%sproject.py' % rpath

		self.available_widgets = [
			'Button', 'Checkbutton', 'Entry', 'Image', 'Label',
			'Listbox', 'Radiobutton', 'Scale', 'Spinbox', 'Text'
		]
		self.valid_colors = [
			'white', 'black', 'red', 'green',
			'blue', 'cyan', 'yellow', 'magenta'
		]
		self.reserved = [
			'REMOVE', 'THEME', 'WRITE', 'TITLE', 'QUIT',
			'APPCOLOR', 'GUI', 'DIMENSIONS', 'BUILD'
			'LOADUSERPROJ', 'MENU', 'MENUCOLOR', 'exit'
		]
		self.reserved += [module for module in dir(modules[__name__])]
		self.reserved += [name for name in dir(builtins) if name.islower()]

		self.master = master
		self.master.title('VisiPy')
		self.master.configure(bg='black')

		self.icon, self.sel, self.popup, self.theme, self.code, self.title = (
			None for _ in range(6))
		self.project, self.theme_layout, self.menu_layout, self.menu_color, \
			self.color, self.xydim, self.font = ({} for _ in range(7))
		self.is_existing = False
		self.updated = True

		self.chars = tuple([str(hex(index)[-1:]) for index in range(16)])

		desktop = environ.get('DESKTOP_SESSION')
		envs = {
			'plasma': '1085x905',
			'ubuntu': '1130x905'
		}
		if desktop is None or desktop == 'ubuntu' or desktop not in envs:
			self.master.geometry(envs['ubuntu'])
			self.listbox_width = 35
			self.master.maxsize('1160', '905')
			self.master.minsize('1130', '750')
		elif desktop == 'plasma':
			self.master.geometry(envs['plasma'])
			self.listbox_width = 30
			self.master.maxsize('1115', '905')
			self.master.minsize('1085', '750')

		self.icon = PhotoImage(file=rpath + 'icon.png')
		master.iconphoto(False, self.icon)

		menu = Menu(self.master)
		menu.config(background=self.light, borderwidth=0, fg=self.dark)
		self.master.config(menu=menu)
		file_menu = Menu(menu)
		file_menu.add_command(
			label='Available Fonts',
			command=self.available_fonts
		)
		file_menu.add_command(
			label='Load Existing',
			command=self.load_user_project
		)
		file_menu.add_command(
			label='Write To File',
			command=self.write_file
		)
		file_menu.add_command(
			label='Reset All',
			command=self.reset
		)
		file_menu.add_command(
			label='Exit', command=self.quit)
		menu.add_cascade(label='File', menu=file_menu)

		edit_menu = Menu(menu)
		edit_menu.add_command(
			label='App Title',
			command=self.app_title
		)
		edit_menu.add_command(
			label='App Color',
			command=self.app_color
		)
		edit_menu.add_command(
			label='Window Menu Color',
			command=self.window_menu_color
		)
		edit_menu.add_command(
			label='Overall Dimensions',
			command=self.display_overall_dimensions
		)
		edit_menu.add_command(
			label='App Theme',
			command=self.app_theme
		)
		menu.add_cascade(label='Edit', menu=edit_menu)

		extras_menu = Menu(menu)
		extras_menu.add_command(
			label='Add 16x16 icon.png',
			command=self.add_icon)
		extras_menu.add_command(
			label='Add Window Menu',
			command=self.add_menu
		)
		menu.add_cascade(label='Extras', menu=extras_menu)

		self.image_path = PhotoImage(file=rpath + 'icon.gif')
		self.image = Label(
			master,
			image=self.image_path,
			borderwidth=0,
			bg='black'
		)
		self.image.grid(row=0, column=0, sticky=W)

		self.overall_dim_label = Label(
			master,
			fg=self.light,
			bg='black',
			anchor=E,
			font=self.bold,
			text='Overall Dimensions:',
			width=20
		)
		self.overall_dim_label.grid(row=0, column=0, sticky=E, padx=5)

		self.xdim = Label(
			master,
			fg=self.light,
			bg='black',
			font=self.small,
			text='L :',
			width=20
		)
		self.xdim.grid(row=0, column=1, sticky=W, padx=5, pady=5)

		self.ydim = Label(
			master,
			fg=self.light,
			bg='black',
			font=self.small,
			text='W :',
			width=20
		)
		self.ydim.grid(row=0, column=1, sticky=E, padx=5, pady=5)

		self.current_build_label = Label(
			master,
			fg=self.light,
			bg='black',
			anchor=CENTER,
			justify=CENTER,
			font=self.bold,
			text='Current Build',
			width=20
		)
		self.current_build_label.grid(
			row=0,
			column=2,
			columnspan=2,
			sticky=E+W,
			padx=5
		)

		scrollbar_build = Scrollbar(master, activebackground=self.light)
		self.build_box = Listbox(
			master,
			fg='#79ff4d',
			activestyle='none',
			bg='black',
			borderwidth=0,
			width=self.listbox_width * 2,
			yscrollcommand=scrollbar_build.set,
			height=47,
			highlightthickness=0,
			selectmode='extended'
		)
		scrollbar_build.config(command=self.build_box.yview)
		scrollbar_build.grid(
			row=0,
			column=2,
			columnspan=2,
			sticky=E,
			padx=5,
			pady=5
		)
		self.build_box.grid(row=1, rowspan=20, column=2, padx=5)

		self.widget_selection = ''
		self.select_new_widget_label = Label(
			master,
			fg=self.light,
			bg='black',
			anchor=W,
			font=self.bold,
			text='Add New Widgets:',
			width=24
		)
		self.select_new_widget_label.grid(
			row=1,
			column=0,
			sticky=W,
			padx=5
		)

		self.select_existing_widget_label = Label(
			master,
			fg=self.light,
			bg='black',
			anchor=W,
			font=self.bold,
			text='Edit Existing Widgets:',
			width=24
		)
		self.select_existing_widget_label.grid(
			row=1,
			column=1,
			sticky=W,
			padx=5
		)

		scrollbar_right = Scrollbar(master, activebackground=self.light)
		self.new_box = Listbox(
			master,
			fg='cyan',
			activestyle='none',
			bg='black',
			borderwidth=0,
			width=self.listbox_width,
			height=10,
			highlightthickness=0,
			selectmode='single'
		)
		self.existing_box = Listbox(
			master,
			fg='cyan',
			activestyle='none',
			bg='black',
			borderwidth=0,
			width=self.listbox_width,
			yscrollcommand=scrollbar_right.set,
			height=10,
			highlightthickness=0,
			selectmode='single'
		)

		self.new_box.grid(row=2, column=0, padx=5, sticky=W)
		self.existing_box.grid(row=2, column=1, padx=5, sticky=E)
		scrollbar_right.config(command=self.existing_box.yview)
		scrollbar_right.grid(row=1, column=1, sticky=E, padx=5, pady=5)

		self.widget_attributes_button = Button(
			master,
			fg=self.light,
			bg=self.dark,
			text='Widget Attributes',
			font=self.small,
			width=16,
			highlightbackground='black',
			highlightthickness=2,
			command=self.set_current_widget
		)
		self.widget_attributes_button.grid(
			row=3,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)

		self.existing_attributes_button = Button(
			master,
			fg=self.light,
			bg=self.dark,
			text='Widget Attributes',
			font=self.small,
			width=16,
			highlightbackground='black',
			highlightthickness=2,
			command=self.set_existing_widget
		)
		self.existing_attributes_button.grid(
			row=3,
			column=1,
			sticky=W,
			padx=5,
			pady=5
		)

		self.remove_button = Button(
			master,
			fg=self.light,
			bg=self.dark,
			text='Remove Widget',
			font=self.small,
			activeforeground='#FF0000',
			highlightbackground='black',
			width=16,
			highlightthickness=2,
			command=self.remove)
		self.remove_button.grid(row=3, column=1, sticky=E, padx=5, pady=5)

		self.widget_type_label = Label(
			master,
			fg=self.light,
			text='Type:',
			bg='black',
			anchor=W,
			width=16,
			height=1,
			font=self.bold
		)
		self.widget_type_label.grid(
			row=4,
			column=0,
			sticky=W,
			padx=5,
			pady=5)

		self.widget_type = Text(
			master,
			fg='cyan',
			bg='black',
			width=23,
			height=1,
			insertontime=600,
			insertbackground='#33CC00',
			insertofftime=100,
			font=self.bold
		)
		self.widget_type.grid(
			row=4,
			column=0,
			sticky=E,
			padx=5,
			pady=15
		)

		self.widget_name_label = Label(
			master,
			fg=self.light,
			text='Name:',
			bg='black',
			anchor=W,
			width=16,
			height=1,
			font=self.bold
		)
		self.widget_name_label.grid(
			row=4,
			column=1,
			sticky=W,
			padx=5,
			pady=5
		)

		self.widget_name = Entry(
			master,
			fg='cyan',
			bg='black',
			width=23,
			insertontime=600,
			insertbackground='#33CC00',
			insertofftime=100,
			font=self.bold,
			disabledforeground='cyan',
			disabledbackground='black'
		)
		self.widget_name.grid(
			row=4,
			column=1,
			sticky=E,
			padx=5,
			pady=15
		)

		self.layout = {}
		self.set_layout()
		self.blackout()

		for index in range(1, int(len(self.layout) / 4) + 1):
			name = 'v%d_label' % index
			self.layout[name].grid(
				row=index + 5, column=0, sticky=W, padx=5, pady=5)
			self.layout[name[:-6]].grid(
				row=index + 5, column=0, sticky=E, padx=30, pady=5)
			name = 'v%d_label' % (index + 11)
			self.layout[name].grid(
				row=index + 5, column=1, sticky=W, padx=5, pady=5)
			self.layout[name[:-6]].grid(
				row=index + 5, column=1, sticky=E, padx=30, pady=5)

		self.blank1 = Label(
			master,
			fg='black',
			bg='black',
			text='',
			width=18,
			height=1
		)
		self.blank1.grid(row=17, column=0, columnspan=2, padx=5)

		self.clear_button = Button(
			master,
			fg=self.light,
			bg=self.dark,
			text='Clear',
			font=self.small,
			width=22,
			highlightbackground='black',
			activeforeground='yellow',
			highlightthickness=2,
			command=self.clear
		)
		self.clear_button.grid(
			row=18,
			column=0,
			sticky=W,
			padx=20,
			pady=5
		)

		self.build_button = Button(
			master,
			fg=self.light,
			bg=self.dark,
			text='Run',
			font=self.small,
			width=22,
			highlightbackground='black',
			activeforeground=self.light,
			highlightthickness=2,
			command=self.build
		)
		self.build_button.grid(
			row=18,
			column=0,
			columnspan=2,
			padx=20,
			pady=5
		)

		self.update_button = Button(
			master,
			fg=self.light,
			bg=self.dark,
			text='Update',
			font=self.small,
			width=22,
			highlightbackground='black',
			activeforeground='#33CC00',
			highlightthickness=2,
			command=self.add_new_widget
		)
		self.update_button.grid(
			row=18,
			column=1,
			sticky=E,
			padx=20,
			pady=5
		)

		self.blank2 = Label(
			master,
			fg='black',
			bg='black',
			text='',
			width=18,
			height=1
		)
		self.blank2.grid(row=19, column=0, padx=5)

		self.status_label = Label(
			master,
			fg=self.light,
			bg='black',
			width=36,
			pady=15,
			font=self.normal
		)
		self.status_label.grid(
			row=20,
			column=0,
			columnspan=2,
			sticky=E+W,
			pady=5
		)

		self.master.grid_columnconfigure(1, weight=1)
		self.master.grid_columnconfigure(2, weight=1)
		self.master.grid_rowconfigure(1, weight=1)
		self.master.grid_rowconfigure(2, weight=1)
		self.master.grid_rowconfigure(16, weight=1)
		self.master.grid_rowconfigure(17, weight=1)
		self.master.grid_rowconfigure(19, weight=1)

		self.populate_add_widgets()
		self.populate_existing_widgets()
		self.populate_code()
		msg = 'Existing widgets found: %d'
		self.set_status(msg % self.existing_box.size())

	def set_layout(self):
		self.layout = {
			'v1': Scale(
				self.master,
				from_=-1, to=100, width=6, length=113,
				borderwidth=0, highlightthickness=0,
				troughcolor='black', showvalue=100,
				bg='black', fg='cyan', orient=HORIZONTAL),
			'v1_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v2': Scale(
				self.master, from_=-1, to=100, width=6, length=113,
				borderwidth=0, highlightthickness=0, troughcolor='black',
				bg='black', fg='cyan', orient=HORIZONTAL),
			'v2_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v3': Scale(
				self.master, from_=-1, to=100, width=6,
				length=113, borderwidth=0, highlightthickness=0,
				troughcolor='black', bg='black', fg='cyan',
				orient=HORIZONTAL, command=self.check_rowspan),
			'v3_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v4': Scale(
				self.master, from_=-1, to=100, width=6,
				length=113, borderwidth=0, highlightthickness=0,
				troughcolor='black', bg='black', fg='cyan',
				orient=HORIZONTAL, command=self.check_columnspan),
			'v4_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v5': Scale(
				self.master, from_=-1, to=100, width=6, length=113,
				borderwidth=0, highlightthickness=0, troughcolor='black',
				bg='black', fg='cyan', orient=HORIZONTAL),
			'v5_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v6': Scale(
				self.master, from_=-1, to=100, width=6, length=113,
				borderwidth=0, highlightthickness=0, troughcolor='black',
				bg='black', fg='cyan', orient=HORIZONTAL),
			'v6_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v7': Scale(
				self.master, from_=-1, to=100, width=6, length=113,
				borderwidth=0, highlightthickness=0, troughcolor='black',
				bg='black', fg='cyan', orient=HORIZONTAL),
			'v7_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v8': Scale(
				self.master, from_=-1, to=100, width=6, length=113,
				borderwidth=0, highlightthickness=0, troughcolor='black',
				bg='black', fg='cyan', orient=HORIZONTAL),
			'v8_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v9': Scale(
				self.master, from_=-1, to=100, width=6, length=113,
				borderwidth=0, highlightthickness=0, troughcolor='black',
				bg='black', fg='cyan', orient=HORIZONTAL),
			'v9_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v10': Scale(
				self.master, from_=-1, to=100, width=6, length=113,
				borderwidth=0, highlightthickness=0, troughcolor='black',
				bg='black', fg='cyan', orient=HORIZONTAL),
			'v10_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v11': Scale(
				self.master, from_=-1, to=100, width=6, length=113,
				borderwidth=0, highlightthickness=0, troughcolor='black',
				bg='black', fg='cyan', orient=HORIZONTAL),
			'v11_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v12': Text(
				self.master, fg='cyan', bg='black', borderwidth=0,
				highlightbackground='#012230', width=16, height=1,
				font=self.normal, insertontime=600,
				insertbackground='#33CC00', insertofftime=100),
			'v12_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v13': Text(
				self.master, fg='cyan', bg='black', borderwidth=0,
				highlightbackground='#012230', width=16, height=1,
				font=self.normal, insertontime=600,
				insertbackground='#33CC00', insertofftime=100),
			'v13_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v14': Text(
				self.master, fg='cyan', bg='black', borderwidth=0,
				highlightbackground='#012230', width=16, height=1,
				font=self.normal, insertontime=600,
				insertbackground='#33CC00', insertofftime=100),
			'v14_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v15': Text(
				self.master, fg='cyan', bg='black', borderwidth=0,
				highlightbackground='#012230', width=16, height=1,
				font=self.normal, insertontime=600,
				insertbackground='#33CC00', insertofftime=100),
			'v15_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v16': Text(
				self.master, fg='cyan', bg='black', borderwidth=0,
				highlightbackground='#012230', width=16, height=1,
				font=self.normal, insertontime=600,
				insertbackground='#33CC00', insertofftime=100),
			'v16_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v17': Text(
				self.master, fg='cyan', bg='black', borderwidth=0,
				highlightbackground='#012230', width=16, height=1,
				font=self.normal, insertontime=600,
				insertbackground='#33CC00', insertofftime=100),
			'v17_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v18': Text(
				self.master, fg='cyan', bg='black', borderwidth=0,
				highlightbackground='#012230', width=16, height=1,
				font=self.normal, insertontime=600,
				insertbackground='#33CC00', insertofftime=100),
			'v18_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v19': Text(
				self.master, fg='cyan', bg='black', borderwidth=0,
				highlightbackground='#012230', width=16, height=1,
				font=self.normal, insertontime=600,
				insertbackground='#33CC00', insertofftime=100),
			'v19_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v20': Text(
				self.master, fg='cyan', bg='black', borderwidth=0,
				highlightbackground='#012230', width=16, height=1,
				font=self.normal, insertontime=600,
				insertbackground='#33CC00', insertofftime=100),
			'v20_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v21': Text(
				self.master, fg='cyan', bg='black', borderwidth=0,
				highlightbackground='#012230', width=16, height=1,
				font=self.normal, insertontime=600,
				insertbackground='#33CC00', insertofftime=100),
			'v21_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
			'v22': Text(
				self.master, fg='cyan', bg='black', borderwidth=0,
				highlightbackground='#012230', width=16, height=1,
				font=self.normal, insertontime=600,
				insertbackground='#33CC00', insertofftime=100),
			'v22_label': Label(
				self.master, fg=self.light, text='', anchor=W,
				disabledforeground='black', font=self.small,
				padx=21, bg='black', width=12, height=1),
		}
		[self.layout[k].set(-1) for k, v in self.layout.items()
			if 'scale' in str(v)]

	def display_button(self):
		self.refresh()
		self.fill_layout({
			'left': {
				'row': -1,
				'column': -1,
				'rowspan': -1,
				'columnspan': -1,
				'padx': -1,
				'pady': -1,
				'width': -1,
				'height': -1,
				'borderwidth': -1,
				'highlightthickness': -1,
				'1': -1
			},
			'right': {
				'foreground': '',
				'background': '',
				'font': '',
				'text': '',
				'activeforeground': '',
				'activebackground': '',
				'highlightcolor': '',
				'anchor': '',
				'sticky': '',
				'command': '',
				'1': ''
			}
		})
		self.disable_remaining()

	def display_checkbutton(self):
		self.refresh()
		self.fill_layout({
			'left': {
				'row': -1,
				'column': -1,
				'rowspan': -1,
				'columnspan': -1,
				'padx': -1,
				'pady': -1,
				'width': -1,
				'height': -1,
				'borderwidth': -1,
				'highlightthickness': -1,
				'1': -1
			},
			'right': {
				'foreground': '',
				'background': '',
				'font': '',
				'text': '',
				'activeforeground': '',
				'activebackground': '',
				'highlightcolor': '',
				'indicatoron': '',
				'selectcolor': '',
				'sticky': '',
				'command': ''
			}
		})
		self.disable_remaining()

	def display_entry(self):
		self.refresh()
		self.fill_layout({
			'left': {
				'row': -1,
				'column': -1,
				'rowspan': -1,
				'columnspan': -1,
				'padx': -1,
				'pady': -1,
				'width': -1,
				'borderwidth': -1,
				'insertontime': -1,
				'insertofftime': -1,
				'1': -1
			},
			'right': {
				'foreground': '',
				'background': '',
				'font': '',
				'show': '',
				'highlightcolor': '',
				'insertbackground': '',
				'selectforeground': '',
				'selectbackground': '',
				'sticky': '',
				'justify': '',
				'1': ''
			}
		})
		self.disable_remaining()

	def display_image(self):
		self.blackout()
		self.refresh()
		img_details = {
			'left': {
				'row': -1,
				'column': -1,
				'rowspan': -1,
				'columnspan': -1,
				'padx': -1,
				'pady': -1,
				'borderwidth': -1
			},
			'right': {
				'image': '',
				'background': '',
				'sticky': ''
			}
		}

		for index in range(1, 5):
			img_details['left'][str(index)] = -1
		for index in range(1, 9):
			img_details['right'][str(index)] = ''

		self.fill_layout(img_details)
		self.disable_remaining()
		if not self.is_existing:
			photo = askopenfilename(
				title='Choose a GIF Image',
				initialdir='/',
				filetypes=(('.gif files', '*.gif'),))

			if not photo or not isfile(photo):
				self.message_thread('Enter a valid image')
			else:
				for k, v in self.layout.items():
					if 'label' in k and self.layout[k].cget('text') == 'image':
						self.layout[k.split('_')[0]].delete('1.0', END)
						self.layout[k.split('_')[0]].insert(END, photo)
						break

	def display_label(self):
		self.refresh()
		self.fill_layout({
			'left': {
				'row': -1,
				'column': -1,
				'rowspan': -1,
				'columnspan': -1,
				'padx': -1,
				'pady': -1,
				'width': -1,
				'height': -1,
				'borderwidth': -1,
				'highlightthickness': -1,
				'1': -1
			},
			'right': {
				'foreground': '',
				'background': '',
				'font': '',
				'text': '',
				'activeforeground': '',
				'activebackground': '',
				'highlightcolor': '',
				'anchor': '',
				'sticky': '',
				'1': '',
				'2': ''
			}
		})
		self.disable_remaining()

	def display_listbox(self):
		self.refresh()
		self.fill_layout({
			'left': {
				'row': -1,
				'column': -1,
				'rowspan': -1,
				'columnspan': -1,
				'padx': -1,
				'pady': -1,
				'width': -1,
				'height': -1,
				'selectborderwidth': -1,
				'borderwidth': -1,
				'highlightthickness': -1
			},
			'right': {
				'foreground': '',
				'background': '',
				'selectmode': '',
				'font': '',
				'activestyle': '',
				'highlightcolor': '',
				'selectforeground': '',
				'selectbackground': '',
				'justify': '',
				'sticky': '',
				'1': ''
			}
		})
		self.disable_remaining()

	def display_radiobutton(self):
		self.refresh()
		self.fill_layout({
			'left': {
				'row': -1,
				'column': -1,
				'rowspan': -1,
				'columnspan': -1,
				'padx': -1,
				'pady': -1,
				'width': -1,
				'height': -1,
				'borderwidth': -1,
				'highlightthickness': -1,
				'1': -1
			},
			'right': {
				'foreground': '',
				'background': '',
				'font': '',
				'text': '',
				'activeforeground': '',
				'activebackground': '',
				'highlightcolor': '',
				'indicatoron': '',
				'selectcolor': '',
				'sticky': '',
				'command': ''
			}
		})
		self.disable_remaining()

	def display_scale(self):
		self.refresh()
		self.fill_layout({
			'left': {
				'row': -1,
				'column': -1,
				'rowspan': -1,
				'columnspan': -1,
				'width': -1,
				'length': -1,
				'from': -1,
				'to': -1,
				'tickinterval': -1,
				'sliderlength': -1,
				'borderwidth': -1
			},
			'right': {
				'foreground': '',
				'background': '',
				'troughcolor': '',
				'activebackground': '',
				'highlightcolor': '',
				'orient': '',
				'sliderrelief': '',
				'relief': '',
				'sticky': '',
				'command': '',
				'1': ''
			}
		})
		self.disable_remaining()

	def display_spinbox(self):
		self.refresh()
		spinbox_layout = {
			'left': {
				'row': -1,
				'column': -1,
				'rowspan': -1,
				'columnspan': -1,
				'padx': -1,
				'pady': -1,
				'width': -1,
				'borderwidth': -1,
				'insertontime': -1,
				'insertofftime': -1,
				'highlightthickness': -1
			},
			'right': {
				'foreground': '',
				'background': '',
				'font': '',
				'activebackground': '',
				'highlightcolor': '',
				'selectforeground': '',
				'selectbackground': '',
				'values': '',
				'sticky': '',
				'command': '',
				'1': ''
			}
		}
		self.fill_layout(spinbox_layout)
		self.disable_remaining()

	def display_text(self):
		self.refresh()
		self.fill_layout({
			'left': {
				'row': -1,
				'column': -1,
				'rowspan': -1,
				'columnspan': -1,
				'padx': -1,
				'pady': -1,
				'width': -1,
				'height': -1,
				'borderwidth': -1,
				'insertontime': -1,
				'insertofftime': -1
			},
			'right': {
				'foreground': '',
				'background': '',
				'font': '',
				'highlightcolor': '',
				'insertbackground': '',
				'selectforeground': '',
				'selectbackground': '',
				'wrap': '',
				'relief': '',
				'sticky': '',
				'1': '',
			}
		})
		self.disable_remaining()

	def display_overall_dimensions(self):
		self.refresh()
		self.blackout()

		self.popup = Toplevel()
		self.popup.title('Overall Dimensions')
		self.popup.geometry(
			"+%d+%d" % (
				self.master.winfo_x() + 100,
				self.master.winfo_y() + 100
			)
		)
		self.popup.configure(bg='black')

		self.xydim = {
			'xdimlabel': Label(
				self.popup,
				fg='white',
				text='x-dimensions:',
				anchor=W,
				bg='black',
				width=18,
				height=1,
				font=self.normal,
				pady=10
			),
			'xdim': Text(
				self.popup,
				fg='cyan',
				bg='black',
				width=18,
				height=1,
				insertontime=600,
				insertbackground='#33CC00',
				insertofftime=100,
				font=self.normal
			),
			'ydimlabel': Label(
				self.popup,
				fg='white',
				text='y-dimensions:',
				anchor=W,
				bg='black',
				width=18,
				height=1,
				font=self.normal,
				pady=10
			),
			'ydim': Text(
				self.popup,
				fg='cyan',
				bg='black',
				width=18,
				height=1,
				insertontime=600,
				insertbackground='#33CC00',
				insertofftime=100,
				font=self.normal
			),
			'ok': Button(
				self.popup,
				fg='green',
				bg='black',
				text='Ok',
				font=self.small,
				width=5,
				command=self.add_dimensions
			),
			'cancel': Button(
				self.popup,
				fg='red',
				bg='black',
				text='Cancel',
				font=self.small,
				width=5,
				command=self.popup.destroy
			),
			'warnlabel': Label(
				self.popup,
				fg='red',
				text='',
				bg='black',
				width=18,
				height=1,
				font=self.normal,
				padx=5,
				pady=5
			)
		}

		self.xydim['xdimlabel'].grid(
			row=0,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.xydim['xdim'].grid(
			row=0,
			column=1,
			sticky=E,
			padx=5,
			pady=5
		)
		self.xydim['ydimlabel'].grid(
			row=1,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.xydim['ydim'].grid(
			row=1,
			column=1,
			sticky=E,
			padx=5,
			pady=5
		)
		self.xydim['ok'].grid(
			row=2,
			column=0,
			sticky=W,
			padx=5,
			pady=5)
		self.xydim['cancel'].grid(
			row=2,
			column=0,
			sticky=E,
			padx=5,
			pady=5
		)
		self.xydim['warnlabel'].grid(
			row=2,
			column=1,
			padx=5,
			pady=5
		)

	def app_theme(self):
		self.refresh()
		self.blackout()
		self.theme = ''
		self.popup = Toplevel()
		self.popup.title('Application Style')
		self.popup.geometry(
			"+%d+%d" % (
				self.master.winfo_x() + 100,
				self.master.winfo_y() + 100
			)
		)
		self.popup.configure(bg='black')

		self.theme_layout = {
			'theme_label': Label(
				self.popup,
				fg='white',
				text='Select TTK Theme:',
				anchor=W,
				bg='black',
				width=18,
				height=1,
				font=self.normal,
				pady=10
			),
			'theme': Spinbox(
				self.popup,
				fg='cyan',
				bg='black',
				width=18,
				values=('', 'clam', 'alt', 'classic', 'default'),
				highlightbackground=self.dark,
				command=self.set_theme
			),
			'ok': Button(
				self.popup,
				fg='green',
				bg='black',
				text='Ok',
				font=self.small,
				width=5,
				command=self.add_style
			),
			'cancel': Button(
				self.popup,
				fg='red',
				bg='black',
				text='Cancel',
				font=self.small,
				width=5,
				command=self.destroy_theme_box
			),
			'warnlabel': Label(
				self.popup,
				fg='red',
				text='',
				bg='black',
				width=18,
				height=1,
				font=self.small,
				pady=10
			)
		}

		self.theme_layout['theme_label'].grid(
			row=0,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.theme_layout['theme'].grid(
			row=0,
			column=1,
			sticky=E,
			padx=5,
			pady=5
		)
		self.theme_layout['ok'].grid(
			row=1,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.theme_layout['cancel'].grid(
			row=1,
			column=0,
			sticky=E,
			padx=5,
			pady=5
		)
		self.theme_layout['warnlabel'].grid(
			row=1,
			column=1,
			padx=5,
			pady=5
		)

	def app_title(self):
		self.refresh()
		self.blackout()
		self.popup = Toplevel()
		self.popup.title('Application Title')
		self.popup.geometry(
			"+%d+%d" % (
				self.master.winfo_x() + 100,
				self.master.winfo_y() + 100
			)
		)
		self.popup.configure(bg='black')

		self.title = {
			'title_label': Label(
				self.popup,
				fg='white',
				text='Enter Title:',
				anchor=W,
				bg='black',
				width=18,
				height=1,
				font=self.normal,
				pady=10
			),
			'title': Text(
				self.popup,
				fg='cyan',
				bg='black',
				width=18,
				height=1,
				insertontime=600,
				insertbackground='#33CC00',
				insertofftime=100,
				font=self.normal
			),
			'ok': Button(
				self.popup,
				fg='green',
				bg='black',
				text='Ok',
				font=self.small,
				width=5,
				command=self.add_title
			),
			'cancel': Button(
				self.popup,
				fg='red',
				bg='black',
				text='Cancel',
				font=self.small,
				width=5,
				command=self.popup.destroy
			),
			'warnlabel': Label(
				self.popup,
				fg='red',
				text='',
				bg='black',
				width=18,
				height=1,
				font=self.small,
				pady=10
			)
		}

		self.title['title_label'].grid(
			row=0,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.title['title'].grid(
			row=0,
			column=1,
			sticky=E,
			padx=5,
			pady=5
		)
		self.title['ok'].grid(
			row=1,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.title['cancel'].grid(
			row=1,
			column=0,
			sticky=E,
			padx=5,
			pady=5
		)
		self.title['warnlabel'].grid(
			row=1,
			column=1,
			padx=5,
			pady=5
		)

	def app_color(self):
		self.refresh()
		self.blackout()
		self.popup = Toplevel()
		self.popup.title('Application Color')
		self.popup.geometry(
			"+%d+%d" % (
				self.master.winfo_x() + 100,
				self.master.winfo_y() + 100
			)
		)
		self.popup.configure(bg='black')

		self.color = {
			'color_label': Label(
				self.popup,
				fg='white',
				text='App Color:',
				anchor=W,
				bg='black',
				width=18,
				height=1,
				font=self.normal,
				pady=10
			),
			'color': Text(
				self.popup,
				fg='cyan',
				bg='black',
				width=18,
				height=1,
				insertontime=600,
				insertbackground='#33CC00',
				insertofftime=100,
				font=self.normal
			),
			'ok': Button(
				self.popup,
				fg='green',
				bg='black',
				text='Ok',
				font=self.small,
				width=5,
				command=self.add_color
			),
			'cancel': Button(
				self.popup,
				fg='red',
				bg='black',
				text='Cancel',
				font=self.small,
				width=5,
				command=self.popup.destroy
			),
			'warnlabel': Label(
				self.popup,
				fg='red',
				text='',
				bg='black',
				width=18,
				height=1,
				font=self.small,
				pady=10
			)
		}

		self.color['color_label'].grid(
			row=0,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.color['color'].grid(
			row=0,
			column=1,
			sticky=E,
			padx=5,
			pady=5
		)
		self.color['ok'].grid(
			row=1,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.color['cancel'].grid(
			row=1,
			column=0,
			sticky=E,
			padx=5,
			pady=5
		)
		self.color['warnlabel'].grid(
			row=1,
			column=1,
			padx=5,
			pady=5
		)

	def add_menu(self):
		self.refresh()
		self.blackout()
		self.popup = Toplevel()
		self.popup.title('Add Menu')
		self.popup.geometry(
			"+%d+%d" % (
				self.master.winfo_x() + 100,
				self.master.winfo_y() + 100
			)
		)
		self.popup.configure(bg='black')

		self.menu_layout = {
			'menu_title_label': Label(
				self.popup,
				fg='white',
				text='Menu Title:',
				anchor=W,
				bg='black',
				width=18,
				height=1,
				font=self.normal,
				pady=10
			),
			'menu_title': Text(
				self.popup,
				fg='cyan',
				bg='black',
				width=18,
				height=1,
				insertontime=600,
				insertbackground='#33CC00',
				insertofftime=100,
				font=self.normal
			),
			'submenus_label': Label(
				self.popup,
				fg='white',
				text='Comma-separated\nSubmenus:',
				anchor=W,
				justify=LEFT,
				bg='black',
				width=18,
				height=1,
				font=self.normal,
				pady=10
			),
			'submenus': Text(
				self.popup,
				fg='cyan',
				bg='black',
				width=18,
				height=1,
				insertontime=600,
				insertbackground='#33CC00',
				insertofftime=100,
				font=self.normal
			),
			'ok': Button(
				self.popup,
				fg='green',
				bg='black',
				text='Ok',
				font=self.small,
				width=5,
				command=self.add_menu_item
			),
			'cancel': Button(
				self.popup,
				fg='red',
				bg='black',
				text='Cancel',
				font=self.small,
				width=5,
				command=self.popup.destroy
			),
			'warnlabel': Label(
				self.popup,
				fg='red',
				text='',
				bg='black',
				width=18,
				height=1,
				font=self.small,
				pady=10
			)
		}

		self.menu_layout['menu_title_label'].grid(
			row=0,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.menu_layout['menu_title'].grid(
			row=0,
			column=1,
			sticky=E,
			padx=5,
			pady=5
		)
		self.menu_layout['submenus_label'].grid(
			row=1,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.menu_layout['submenus'].grid(
			row=1,
			column=1,
			sticky=E,
			padx=5,
			pady=5
		)
		self.menu_layout['ok'].grid(
			row=2,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.menu_layout['cancel'].grid(
			row=2,
			column=0,
			sticky=E,
			padx=5,
			pady=5
		)
		self.menu_layout['warnlabel'].grid(
			row=3,
			column=1,
			padx=5,
			pady=5
		)

	def window_menu_color(self):
		self.refresh()
		self.blackout()
		self.popup = Toplevel()
		self.popup.title('Window Menu Color')
		self.popup.geometry(
			"+%d+%d" % (
				self.master.winfo_x() + 100,
				self.master.winfo_y() + 100
			)
		)
		self.popup.configure(bg='black')

		self.menu_color = {
			'fg_label': Label(
				self.popup,
				fg='white',
				text='Foreground Color:',
				anchor=W,
				bg='black',
				justify=LEFT,
				width=18,
				height=1,
				font=self.normal,
				pady=10
			),
			'fg': Text(
				self.popup,
				fg='cyan',
				bg='black',
				width=18,
				height=1,
				insertontime=600,
				insertbackground='#33CC00',
				insertofftime=100,
				font=self.normal
			),
			'bg_label': Label(
				self.popup,
				fg='white',
				text='Background Color:',
				anchor=W,
				bg='black',
				justify=LEFT,
				width=18,
				height=1,
				font=self.normal,
				pady=10
			),
			'bg': Text(
				self.popup,
				fg='cyan',
				bg='black',
				width=18,
				height=1,
				insertontime=600,
				insertbackground='#33CC00',
				insertofftime=100,
				font=self.normal
			),
			'ok': Button(
				self.popup,
				fg='green',
				bg='black',
				text='Ok',
				font=self.small,
				width=5,
				command=self.add_menu_color
			),
			'cancel': Button(
				self.popup,
				fg='red',
				bg='black',
				text='Cancel',
				font=self.small,
				width=5,
				command=self.popup.destroy
			),
			'warnlabel': Label(
				self.popup,
				fg='red',
				text='',
				bg='black',
				width=18,
				height=1,
				font=self.small,
				pady=10
			)
		}

		self.menu_color['fg_label'].grid(
			row=0,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.menu_color['fg'].grid(
			row=0,
			column=1,
			sticky=E,
			padx=5,
			pady=5
		)
		self.menu_color['bg_label'].grid(
			row=1,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.menu_color['bg'].grid(
			row=1,
			column=1,
			sticky=E,
			padx=5,
			pady=5
		)
		self.menu_color['ok'].grid(
			row=2,
			column=0,
			sticky=W,
			padx=5,
			pady=5
		)
		self.menu_color['cancel'].grid(
			row=2,
			column=0,
			sticky=E,
			padx=5,
			pady=5
		)
		self.menu_color['warnlabel'].grid(
			row=3,
			column=1,
			padx=5,
			pady=5
		)

	def available_fonts(self):
		self.popup = Toplevel()
		self.popup.title('Available Fonts')
		self.popup.geometry(
			"+%d+%d" % (
				self.master.winfo_x() + 100,
				self.master.winfo_y() + 100
			)
		)
		self.popup.configure(bg='black')

		directions = '<FontName>  <FontSize>\n\n'
		directions += 'FontSize is an even integer between 8 & 64.\n'
		directions += '\nExamples:\nTkTextFont  8\nTkHeadingFont  24'
		scrollbar = Scrollbar(self.popup, activebackground=self.light)

		self.font = {
			'font_directions': Label(
				self.popup,
				fg=self.light,
				text=directions,
				bg='black',
				width=24,
				height=8,
				justify=LEFT,
				font=self.small,
				pady=15
			),
			'fontbox': Listbox(
				self.popup,
				fg='cyan',
				activestyle='none',
				bg='black',
				borderwidth=0,
				width=30,
				yscrollcommand=scrollbar.set,
				height=8,
				highlightcolor='black',
				selectmode='single'
			),
			'ok': Button(
				self.popup,
				fg='green',
				bg='black',
				text='Ok',
				font=self.small,
				width=5,
				command=self.add_current_font
			),
			'cancel': Button(
				self.popup,
				fg='red',
				bg='black',
				text='Cancel',
				font=self.small,
				width=5,
				command=self.popup.destroy
			),
			'fontsize': Scale(
				self.popup,
				from_=8,
				to=64,
				width=6,
				length=113,
				borderwidth=0,
				resolution=2,
				highlightthickness=0,
				troughcolor='black',
				bg='black',
				fg='cyan',
				orient=HORIZONTAL,
				command=self.set_font
			),
			'fontlabel': Label(
				self.popup,
				fg='red',
				text='Font Size:',
				bg='black',
				width=18,
				height=1,
				font=self.small,
				pady=10
			),
			'font': ''
		}

		scrollbar.config(command=self.font['fontbox'].yview)

		self.font['font_directions'].grid(
			row=0,
			sticky=E + W,
			padx=5,
			pady=5
		)
		scrollbar.grid(
			row=1,
			sticky=E,
			padx=5,
			pady=5
		)
		self.font['fontbox'].grid(
			row=2,
			padx=5,
			pady=5,
			sticky=E + W
		)

		self.font['fontbox'].delete(0, END)
		[self.font['fontbox'].insert(END, font) for font in sorted(names())]

		self.font['fontsize'].grid(
			row=3,
			padx=5,
			pady=5,
			sticky=E
		)
		self.font['fontlabel'].grid(
			row=3,
			padx=5,
			pady=5,
			sticky=W
		)
		self.font['ok'].grid(
			row=4,
			sticky=W,
			padx=5,
			pady=5
		)
		self.font['cancel'].grid(
			row=4,
			sticky=E,
			padx=5,
			pady=5
		)
		self.font['fontbox'].select_set(0)

	def message_thread(self, message):
		th = Thread(target=self.status, args=[message])
		th.start()

	def status(self, msg):
		self.set_status(msg)
		self.status_label.update_idletasks()
		sleep(4)
		self.set_status('')

	def warn_thread(self, message, attribute):
		th = Thread(target=self.warn, args=[attribute, message])
		th.start()

	def warn(self, attr, msg):
		if attr != 'name':
			self.set_status(msg)
			for k, v in self.layout.items():
				if 'label' in k and self.layout[k].cget('text') == attr:
					self.layout[k].configure(fg='red')
					self.layout[k].update_idletasks()
					break
			sleep(2)
			self.layout[k].configure(fg=self.light)
		else:
			self.set_status(msg)
			self.widget_name_label.configure(fg='red')
			self.widget_name_label.update_idletasks()
			sleep(2)
			self.widget_name_label.configure(fg=self.light)
		sleep(2)
		self.set_status('')

	def poll_update(self):
		poll_count = 0
		project_json = self.data_path + '.update'
		while not isfile(project_json):
			sleep(.02)
			if poll_count > 250:
				self.updated = False
				break
			poll_count += 1

	def add_new_widget(self):
		self.update('ADD')

	def update(self, action, changes=None):
		piped = True
		if action == 'ADD':
			output_vals = self.review()
			if 'ERROR' in output_vals:
				[(attr, allowable)] = output_vals['ERROR'].items()
				if isinstance(allowable, list):
					msg = 'Error, Allowable values: '
					msg += ', '.join(allowable)
				else:
					msg = 'Error: '
					msg += allowable
				piped = False
				self.warn_thread(msg, attr)
			elif 'row' not in output_vals:
				piped = False
				self.warn_thread('row attributes required', 'row')
			elif 'column' not in output_vals:
				piped = False
				self.warn_thread('column attributes required', 'column')
			else:
				stdout.write('%s|$|%s|$|%s\n' % (action, self.sel, '|:|'.join(
					['%s|@|%s' % (k, v) for k, v in output_vals.items()])))
				stdout.flush()
		elif action == 'EXIT':
			piped = False
			stdout.write('%s\n' % action)
			stdout.flush()
		elif action in ('RESET', 'BUILD'):
			stdout.write('%s\n' % action)
			stdout.flush()
		else:
			# REMOVE, THEME, WRITE, TITLE, APPCOLOR
			# ICON, DIMENSIONS, LOADUSERPROJ, MENU, MENUCOLOR
			stdout.write('%s|$|%s\n' % (action, changes))
			stdout.flush()
		if piped:
			poll = Thread(target=self.poll_update)
			poll.start()
			poll.join()
			if not self.updated:
				self.warn('ERROR', 'Error writing')
				exit(1)
			rename(self.data_path + '.update', self.data_path)
			rename(self.code_path + '.update', self.code_path)
			self.populate_existing_widgets()
			self.populate_code()
			self.refresh()
			self.blackout()

	def valid_color(self, color):
		color = color.lower()
		if color in self.valid_colors:
			return True
		return all([c in self.chars for c in color]) and len(color) == 6

	def review(self):
		int_vals = (
			'row', 'column', 'padx', 'pady', 'width', 'height',
			'borderwidth', 'highlightthickness', 'insertontime',
			'insertofftime', 'selectborderwidth', 'to', 'from',
			'tickinterval', 'sliderlength', 'length'
		)
		color_types = [
			'activeforeground', 'insertbackground', 'activebackground',
			'highlightcolor', 'selectcolor', 'background', 'foreground',
			'selectbackground', 'selectforeground', 'troughcolor'
		]

		color_list = self.valid_colors + ['or 6-digit vex value']

		valid = {
			'activeforeground': color_list,
			'insertbackground': color_list,
			'activebackground': color_list,
			'highlightcolor': color_list,
			'selectcolor': color_list,
			'background': color_list,
			'foreground': color_list,
			'selectbackground': color_list,
			'selectforeground': color_list,
			'troughcolor': color_list,
			'selectmode': ['BROWSE', 'SINGLE', 'MULTIPLE', 'EXTENDED'],
			'relief': ['FLAT', 'RAISED', 'GROOVE', 'SUNKEN', 'RIDGE'],
			'sliderrelief': ['FLAT', 'RAISED', 'GROOVE', 'SUNKEN', 'RIDGE'],
			'anchor': ['N', 'NE', 'E', 'SE', 'S', 'SW', 'W', 'NW', 'CENTER'],
			'activestyle': ['DOTBOX', 'NONE', 'UNDERLINE'],
			'indicatoron': ['True', 'False'],
			'orient': ['HORIZONTAL', 'VERTICAL'],
			'wrap': ['CHAR', 'WORD'],
			'font': "see 'Available Fonts' in File menu",
			'sticky': [
				'E', 'W', 'N', 'S', 'E+W', 'N+S',
				'N+W', 'S+W', 'S+E', 'N+E', 'W+E+N+S'
			]
		}

		widget_dict = {}
		for index, (k, v) in enumerate(self.layout.items()):

			if 'label' not in k:
				continue
			key = self.layout[k].cget('text')
			if key.isdigit():
				continue
			try:
				value = self.layout[k.split('_')[0]].get()
			except Exception:
				value = self.layout[k.split('_')[0]].get('1.0', END).strip()

			if value == -1 or (not value and isinstance(value, str)):
				continue

			if key == 'text' or key in int_vals:
				widget_dict[key] = value
				continue

			if key in ('rowspan', 'columnspan'):
				if value > 0:
					widget_dict[key] = value
				continue

			if key == 'values':
				if ',' not in value:
					return {'ERROR': {key: 'Use a comma-separated list'}}
				widget_dict[key] = ', '.join(
					['\'%s\'' % i.strip() for i in value.split(',') if i])
				continue

			if key == 'justify':
				val = value.upper()
				if val not in ('LEFT', 'CENTER', 'RIGHT'):
					return {'ERROR': {key: 'Use LEFT, CENTER, or RIGHT'}}
				widget_dict[key] = val
				continue

			if key == 'command':
				if value in self.reserved or value in self.available_widgets:
					msg = 'Invalid  Name, %s is reserved' % value
					return {'ERROR': {key: msg}}
				if not value.isidentifier() or iskeyword(value):
					msg = 'Only use valid Python naming conventions'
					return {'ERROR': {key: msg}}
				widget_dict[key] = value
				continue

			if key == 'font':
				try:
					font_type, font_size = value.split()
				except Exception:
					return {'ERROR': {key: valid[key]}}
				sizes = [str(i) for i in range(8, 65, 2)]
				if font_type not in names() or font_size not in sizes:
					return {'ERROR': {key: valid['font']}}
				widget_dict[key] = value
				continue

			if key == 'image':
				if isfile(value):
					widget_dict[key] = value
					continue
				return {'ERROR': {key: 'Invalid image file'}}

			if key == 'show':
				if value != '*':
					msg = 'Leave this field blank, or use * for password entry'
					return {'ERROR': {key: msg}}
				widget_dict[key] = value
				continue

			for sep in '|:|', '|$|', '|@|':
				if not isinstance(value, int) and sep in value:
					return {'ERROR': {key: 'Chars not allowed: %s' % value}}

			check = {
				'activeforeground': self.valid_color(value),
				'insertbackground': self.valid_color(value),
				'activebackground': self.valid_color(value),
				'highlightcolor': self.valid_color(value),
				'selectcolor': self.valid_color(value),
				'background': self.valid_color(value),
				'foreground': self.valid_color(value),
				'selectbackground': self.valid_color(value),
				'selectforeground': self.valid_color(value),
				'troughcolor': self.valid_color(value),
				'relief': value in valid[key],
				'sliderrelief': value in valid[key],
				'selectmode': value in valid[key],
				'anchor': value in valid[key],
				'activestyle': value in valid[key],
				'indicatoron': value in valid[key],
				'orient': value in valid[key],
				'wrap': value in valid[key],
				'sticky': value in valid[key]
			}

			if not check[key]:
				return {'ERROR': {key: valid[key]}}
			widget_dict[key] = value

		name = self.widget_name.get().strip()
		if name in self.reserved or name in self.available_widgets:
			msg = 'Invalid  Name %s is reserved' % name
			return {'ERROR': {'name': msg}}
		if not name or not name.isidentifier() or iskeyword(name):
			msg = 'Invalid  Name, use Python naming conventions'
			return {'ERROR': {'name': msg}}
		if not self.is_existing and name in self.project:
			msg = 'Invalid  Name, name already exists'
			return {'ERROR': {'name': msg}}
		widget_dict['name'] = name

		for k, v in widget_dict.items():
			if k in color_types and v not in self.valid_colors:
				widget_dict[k] = '#%s' % v
		return widget_dict

	def add_icon(self):
		self.refresh()
		self.blackout()
		self.icon = askopenfilename(
			title='16x16 .png',
			filetypes=(("png files", "*.png"),)
		)
		if not self.icon:
			self.message_thread('No icon path entered')
		elif not isfile(self.icon):
			self.icon = ''
			self.message_thread('Invalid icon path entered')
		else:
			self.message_thread(self.icon)
			self.update('ICON', changes=self.icon)

	def add_dimensions(self):
		try:
			xdim = int(self.xydim['xdim'].get('1.0', END).strip())
			ydim = int(self.xydim['ydim'].get('1.0', END).strip())
		except Exception:
			self.xydim['warnlabel'].configure(text='Integers only')
		else:
			if not all([xdim < 2000, ydim < 2000, xdim > 0, ydim > 0]):
				self.xydim['warnlabel'].configure(text='0 - 2000 only')
			else:
				self.popup.destroy()
				self.update('DIMENSIONS', changes='%sx%s' % (xdim, ydim))
				self.popup.destroy()

	def remove(self):
		try:
			name = self.existing_box.get(self.existing_box.curselection())
		except Exception:
			self.message_thread('Select an existing widget')
		else:
			self.update('REMOVE', name)

	def add_style(self):
		if not self.theme:
			self.theme_layout['warnlabel'].configure(text='Select a theme')
		else:
			self.popup.destroy()
			self.update('THEME', self.theme)
			self.theme = ''

	def add_title(self):
		title = self.title['title'].get('1.0', END).strip()
		if not title:
			self.title['warnlabel'].configure(text='Enter a title')
		elif not title.isidentifier() or iskeyword(title):
			msg = 'Use valid Python\nnaming conventions'
			self.title['warnlabel'].configure(text=msg)
		else:
			self.popup.destroy()
			self.update('TITLE', changes=title)
			self.title = None

	def add_menu_item(self):
		menu = self.menu_layout['menu_title'].get('1.0', END).strip()
		subs = self.menu_layout['submenus'].get('1.0', END).strip()

		if not menu:
			self.menu_layout['warnlabel'].configure(text='Missing menu title')
		elif not subs:
			self.menu_layout['warnlabel'].configure(text='Missing submenus')
		else:
			msg = '%s,%%s' % menu
			msg = msg % ','.join([sm.strip() for sm in subs.split(',') if sm])
			self.popup.destroy()
			self.update('MENU', changes=msg)
			self.menu_layout = {}

	def add_color(self):
		color = self.color['color'].get('1.0', END).strip()

		if self.valid_color(color):
			if color not in self.valid_colors:
				color = '#%s' % color
			self.popup.destroy()
			self.update('APPCOLOR', changes=color)
			self.color = {}
		else:
			self.color['warnlabel'].configure(text='Invalid app color')
			self.color['color'].delete('1.0', END)

	def add_menu_color(self):
		fg = self.menu_color['fg'].get('1.0', END).strip()
		bg = self.menu_color['bg'].get('1.0', END).strip()
		if not self.valid_color(fg):
			self.menu_color['warnlabel'].configure(text='Invalid foreground')
		elif not self.valid_color(bg):
			self.menu_color['warnlabel'].configure(text='Invalid background')
		else:
			if fg not in self.valid_colors:
				fg = '#%s' % fg
			if bg not in self.valid_colors:
				bg = '#%s' % bg
			self.popup.destroy()
			self.update('MENUCOLOR', changes='%s|:|%s' % (fg, bg))
			self.menu_color = {}

	def set_font(self, current_font):
		self.font['font'] = current_font

	def add_current_font(self):
		font = self.font['fontbox'].get(self.font['fontbox'].curselection())
		font_size = self.font['font']
		if not font_size:
			font_size = '8'
		for k, v in self.layout.items():
			if 'label' in k and self.layout[k].cget('text') == 'font':
				text_box = k.split('_')[0]
				self.layout[text_box].delete('1.0', END)
				self.layout[text_box].insert(END, '%s %s' % (font, font_size))
		self.popup.destroy()

	def build(self):
		self.message_thread('Executing current build...')
		self.update('BUILD')

	def fill_current_widget(self):
		if not self.sel:
			self.message_thread('Please select a widget')
		elif self.sel == 'Button':
			self.set_name_txt('Button')
			self.display_button()
		elif self.sel == 'Checkbutton':
			self.set_name_txt('Checkbutton')
			self.display_checkbutton()
		elif self.sel == 'Entry':
			self.set_name_txt('Entry')
			self.display_entry()
		elif self.sel == 'Image':
			self.set_name_txt('Image')
			self.display_image()
		elif self.sel == 'Label':
			self.set_name_txt('Label')
			self.display_label()
		elif self.sel == 'Listbox':
			self.set_name_txt('Listbox')
			self.display_listbox()
		elif self.sel == 'Radiobutton':
			self.set_name_txt('Radiobutton')
			self.display_radiobutton()
		elif self.sel == 'Scale':
			self.set_name_txt('Scale')
			self.display_scale()
		elif self.sel == 'Spinbox':
			self.set_name_txt('Spinbox')
			self.display_spinbox()
		elif self.sel == 'Text':
			self.set_name_txt('Text')
			self.display_text()

	def set_current_widget(self):
		self.is_existing = False
		try:
			self.sel = self.new_box.get(self.new_box.curselection())
		except Exception:
			self.sel = ''
		else:
			self.widget_name.configure(state=NORMAL)
			self.fill_current_widget()

	def set_existing_widget(self):
		self.widget_name.configure(state=NORMAL)
		self.is_existing = True
		try:
			name = self.existing_box.get(self.existing_box.curselection())
		except Exception:
			self.message_thread('Select a widget')
		else:
			self.load_project_json()
			try:
				self.sel = self.project[name]['widget']
				if self.sel == 'Spinbox' and 'values' in self.project[name]:
					tmp = self.project[name]['values'].replace('\'', '')
					self.project[name]['values'] = tmp

				for attr, value in self.project[name].items():
					if isinstance(value, str) and value[:1] == '#':
						self.project[name][attr] = value[1:]

				self.fill_current_widget()
				self.fill_existing_widget(name)
			except Exception:
				self.sel = self.project[name]
				if name != 'ICON':
					self.add_menu()
					self.menu_layout['menu_title'].delete('1.0', END)
					self.menu_layout['menu_title'].insert(END, name)
					subs = ','.join([v for k, v in self.sel.items()])
					self.menu_layout['submenus'].delete('1.0', END)
					self.menu_layout['submenus'].insert(END, subs)
				else:
					self.add_icon()
			self.widget_name.configure(state=DISABLED)

	def fill_existing_widget(self, name):
		self.widget_name.delete(0, END)
		self.widget_name.insert(0, name)
		for k, v in self.layout.items():
			if 'label' in k and not self.layout[k].cget('text').isnumeric():
				label = self.layout[k].cget('text')
				vx = k.split('_')[0]
				if label in self.project[name]:
					if int(str(vx[1:])) < 12:
						self.layout[vx].set(self.project[name][label])
					else:
						self.layout[vx].delete('1.0', END)
						self.layout[vx].insert(END, self.project[name][label])

	def fill_layout(self, layout):
		index = 1
		for side in 'left', 'right':
			for k, v in layout[side].items():
				label = 'v%d_label' % index
				self.layout[label].configure(text=k)
				if side != 'left':
					self.layout[label[:-6]].delete('1.0', END)
					self.layout[label[:-6]].insert(END, v)
				index += 1

	def disable(self, key):
		self.layout[key].configure(state=DISABLED)
		target = key.split('_')[0]
		self.layout[target].config(
			fg='black',
			highlightthickness=0,
			state=DISABLED
		)
		if int(target[1:]) < 12:
			self.layout[target].config(sliderrelief=FLAT)

	def blackout(self):
		for k, _ in self.layout.items():
			if 'label' in k:
				self.disable(k)

	def disable_remaining(self):
		for k, _ in self.layout.items():
			if 'label' in k and self.layout[k].cget('text').isnumeric():
				self.disable(k)

	def destroy_theme_box(self):
		self.theme = ''
		self.popup.destroy()

	def quit(self):
		if askyesno('Exit', 'Exit Visipy?'):
			self.set_status('Exiting')
			self.update('EXIT')
			Tk().quit()

	def reset(self):
		if askyesno('Reset Project', 'Clear entire project?'):
			self.set_status('Resetting')
			self.update('RESET')

	def set_theme(self):
		self.theme = self.theme_layout['theme'].get()

	def populate_add_widgets(self):
		[self.new_box.insert(END, w) for w in self.available_widgets]

	def populate_existing_widgets(self):
		self.load_project_json()
		widgets = [key for key in self.project if key not in self.reserved]
		self.existing_box.delete(0, END)
		_ = [self.existing_box.insert(END, widget) for widget in widgets]
		dim = self.project['DIMENSIONS'].get('dimensions').strip().split('x')
		self.xdim.config(text='L : %s' % dim[0])
		self.ydim.config(text='W : %s' % dim[1])

	def populate_code(self):
		self.load_project_code()
		self.build_box.delete(0, 'end')
		self.build_box.insert('end', *self.code.split('\n'))

	def write_file(self):
		file_path = asksaveasfilename(
			title='Write Current Build',
			filetypes=(('all files', '*.*'),)
		)
		if file_path:
			self.message_thread('creating project files')
			self.update('WRITE', file_path)

	def load_user_project(self):
		user_path = askopenfilename()
		if not user_path or not isfile(user_path):
			self.message_thread('Invalid file chosen')
		else:
			self.update('LOADUSERPROJ', changes=user_path)

	def load_project_json(self):
		with open(self.data_path) as file_in:
			self.project = load(file_in)

	def load_project_code(self):
		with open(self.code_path) as code_in:
			self.code = code_in.read()

	def set_status(self, incoming):
		self.status_label.config(text=incoming)
		self.status_label.update_idletasks()

	def set_name_txt(self, incoming):
		if self.widget_type['state'] != 'normal':
			self.widget_type.config(state=NORMAL)
		self.message_thread('Edit %s' % incoming)
		self.widget_type.delete('1.0', END)
		self.widget_type.insert(END, incoming)
		if self.sel in self.available_widgets:
			self.widget_type.config(state=DISABLED)

	def check_rowspan(self, current):
		if current == '0' and self.layout['v3'].get() == 0:
			self.set_status('0 is not allowed for rowspan')
			self.layout['v3_label'].configure(fg='red')
		else:
			self.layout['v3_label'].configure(fg=self.light)
			self.set_status('')

	def check_columnspan(self, current):
		if current == '0' and self.layout['v4'].get() == 0:
			self.set_status('0 is not allowed for columnspan')
			self.layout['v4_label'].configure(fg='red')
		else:
			self.layout['v4_label'].configure(fg=self.light)
			self.set_status('')

	def refresh(self):
		self.clear()
		for _, v in self.layout.items():
			if v['state'] != 'disabled':
				continue
			widget_type = str(v)
			if 'label' in widget_type:
				v.config(
					fg=self.light,
					state=NORMAL
				)
			elif 'text' in widget_type:
				v.config(
					fg='cyan',
					highlightthickness=1,
					state=NORMAL
				)
			elif 'scale' in widget_type:
				v.config(
					sliderrelief=SUNKEN,
					fg='cyan',
					state=NORMAL
				)

	def clear(self):
		self.widget_name.delete(0, END)
		for k, v in self.layout.items():
			if 'label' not in k:
				try:
					self.layout[k].set(-1)
				except Exception:
					self.layout[k].delete('1.0', END)


if __name__ == "__main__":
	root = Tk()
	root.style = Style()
	root.style.theme_use('clam')
	VisiPy(root)
	root.mainloop()
`)
	return bytes.ReplaceAll(vpy, []byte{0x09}, []byte{0x20, 0x20, 0x20, 0x20})
}
