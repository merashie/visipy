<h3>Visipy - A Tkinter Widget-Templater/GUI Builder for Python</h3>



<img src="https://github.com/rootVIII/visipy/blob/master/example.png" alt="example">




###### Installation
 <ul>
  <li><code>git clone</code> the project or <code>go get github.com/rootVIII/visipy/visipy</code></li>
  <li>run: <code>cd visipy/; go run visipy.go</code></li>
  <li>build: <code>go build visipy.go</code></li>
  <li>If you don't have Go or don't want to build the exe, there is an Ubuntu 64-bit Build in <a href="https://github.com/rootVIII/visipy/tree/master/bin">bin/</a> (Mac OS Build is there too but not tested yet-UI may need adjustments)</li>
</ul> 



###### Requirements
<ul>
  <li><code>Python3</code> (must also be in your path)</li>
  <li>Internet connection (only at application start-up)</li>
  <li>Linux or Mac OS (Not yet tested on Mac)</li>
</ul>



###### Youtube Demo
<ul>
  <li>Also check out the <a href="https://youtu.be/i0pYqYdM1VQ" target="_blank">demo screencast</a> for example operation</li>
</ul>



###### Visipy Features
<ul>
  <li>Allow for rapid creation of a single Class, utilizing Tk() init and <i>N</i> widgets</li>
  <li>Easy Tkinter widget creation, styling, and editing</li>
  <li>Allow for easy, visual manipulation of the Tk grid() geometry manager</li>
  <li>Run your project with the press of a button at any time to view changes</li>
  <li>See the current codebase change as widgets are added, edited, or removed</li>
  <li><b>After templating/GUI building, it's up to you to finish the program (action handlers, fine-tuning etc.) in your own code editor</b></li>
</ul>



###### Widgets
<ul>
  <li>While Visipy doesn't offer all of the Tk widgets and Tk widget attributes, it does include a large portion of them</li>
  <li>Visipy will template an entire initializer, along with any method names for Tk widgets that have the <b>command</b> option</li>
  <li> Work with the following widgets and features:
    <ul>
      <li>Button</li>
      <li>Checkbutton</li>
      <li>Entry</li>
      <li>Image</li>
      <li>Label</li>
      <li>Listbox</li>
      <li>Radiobutton</li>
      <li>Scale</li>
      <li>Spinbox</li>
      <li>Text</li>
      <li>Overall App Dimensions</li>
      <li>Overall App Theme</li>
      <li>Overall App Color</li>
      <li>Window Bar Menus</li>
      <li>16x16 icon.png</li>
    </ul>
  </li>
</ul>



###### End Result
<ul>
  <li>A runnable, Python3 Tkinter GUI (pycodestyle valid unless really long path names are used etc.)</li>
  <li>Method placeholders for TODO/Action handling (such as a Button's associated command)</li>
  <li>The current build/code can be written to disk any time</li>
  <li>A JSON representation is also available - used to reload the project for later editing</li>
</ul>



###### Things to Note:
<ul>
  <li>The current GUI build can be written to a <code>.py</code> file at any time with the Write to File option</li>
  <li>A <code>.project</code> file (JSON) will also be created in the same directory as your <code>.py</code> file</li>
  <li>The <code>.project</code> file may be discarded or saved to reload the project later to continue working on the same project (<b>do not</b> edit the JSON file)</li>
  <li>The current build/GUI should be runnable at all times, easing the creation of your application</li>
  <li>This does not mean however that your app is going to look as intended</li>
  <li>The following link describes many widgets, and their available attributes: <a href="http://effbot.org/tkinterbook/tkinter-classes.htm" target="_blank">tkinter book</a></li>
  <li>If you are unsure what to put for a widget's attribute, put any value and try running Update; allowable values will be suggested for you if something invalid is found</li>
  <li>Hitting the run button after each widget addition, edit, or removal is a good way to double-check your GUI along with examining the current build window each time code is updated</li>
</ul>



<hr>
This project was developed on Ubuntu 18.04.4 LTS