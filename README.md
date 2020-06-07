<h3>Visipy - A Tkinter Widget-Templater/GUI Builder for Python</h3>



<img src="https://github.com/rootVIII/visipy/blob/master/example.png" alt="example">




###### Installation
 <ul>
  <li style="list-style-type:square"><code>git clone</code> the project or <code>go get github.com/rootVIII/visipy/visipy</code>  and navigate to project root</li>
  <li style="list-style-type:square">run: <code>go run visipy/visipy.go</code></li>
  <li style="list-style-type:square">build: <code>go build visipy/visipy.go</code></li>
  <li style="list-style-type:square">If you don't have Go or don't want to build the exe, there is an Ubuntu 64-bit Build in <a href="https://github.com/rootVIII/visipy/tree/master/bin">bin/</a> (Mac OS Build coming soon)</li>
</ul> 



###### Requirements
<ul>
  <li style="list-style-type:square"><code>Python3</code> (must also be in your path)</li>
  <li style="list-style-type:square">Internet connection (only at application start-up)</li>
  <li style="list-style-type:square">Linux or Mac OS (Not yet tested on Mac)</li>
</ul>



###### Visipy Features
<ul>
  <li style="list-style-type:square">Allow for rapid creation of a single Class, utilizing Tk() init and <i>N</i> widgets</li>
  <li style="list-style-type:square">Easy Tkinter widget creation, styling, and editing</li>
  <li style="list-style-type:square">Allow for easy, visual manipulation of the Tk grid() geometry manager</li>
  <li style="list-style-type:square">Run your project with the press of a button at any time to view changes</li>
  <li style="list-style-type:square">See the current codebase change as widgets are added, edited, or removed</li>
  <li style="list-style-type:square"><b>After templating/GUI building, it's up to you to finish the program (action handlers, fine-tuning etc.) in your own code editor</b></li>
</ul>



###### Widgets
<ul>
  <li style="list-style-type:square">While Visipy doesn't offer all of the Tk widgets and Tk widget attributes, it does include a large portion of them</li>
  <li style="list-style-type:square">Visipy will template an entire initializer, along with any method names for Tk widgets that have the <b>command</b> option</li>
  <li style="list-style-type:square"> Work with the following widgets and features:
    <ul>
      <li style="list-style-type:square">Button</li>
      <li style="list-style-type:square">Checkbutton</li>
      <li style="list-style-type:square">Entry</li>
      <li style="list-style-type:square">Image</li>
      <li style="list-style-type:square">Label</li>
      <li style="list-style-type:square">Listbox</li>
      <li style="list-style-type:square">Radiobutton</li>
      <li style="list-style-type:square">Scale</li>
      <li style="list-style-type:square">Spinbox</li>
      <li style="list-style-type:square">Text</li>
      <li style="list-style-type:square">Overall App Dimensions</li>
      <li style="list-style-type:square">Overall App Theme</li>
      <li style="list-style-type:square">Overall App Color</li>
      <li style="list-style-type:square">Window Bar Menus</li>
      <li style="list-style-type:square">16x16 icon.png</li>
    </ul>
  </li>
</ul>



###### End Result
<ul>
  <li style="list-style-type:square">A runnable, Python3 Tkinter GUI (pycodestyle valid unless really long path names are used etc.)</li>
  <li style="list-style-type:square">Method placeholders for TODO/Action handling (such as a Button's associated command)</li>
  <li style="list-style-type:square">The current build/code can be written to disk any time</li>
  <li style="list-style-type:square">A JSON representation is also available - used to reload the project for later editing</li>
</ul>



###### Things to Note:
<ul>
  <li style="list-style-type:square">This project has been tested Linux (Ubuntu & Plasma desktop environments). It "should" run on Mac as well, but the GUI may or may not look proportionate until testing/updates can be done (the same goes for other Linux desktop environments).</li>
  <li style="list-style-type:square">When the current build is written to a <code>.py</code> file (current code build), a <code>.project</code> file (JSON) is also created in the same directory</li>
  <li style="list-style-type:square">The <code>.project</code> file may be discarded or saved to reload the project later (<b>do not</b> edit the JSON file)</li>
  <li style="list-style-type:square">The project should run with the Run button even if the settings you have are a bit off (or a lot off)</li>
  <li style="list-style-type:square">This does not mean however that your app is going to look as intended</li>
  <li style="list-style-type:square">The following link describes many widgets, and some of the options that are available: <a href="http://effbot.org/tkinterbook/tkinter-classes.htm" target="_blank">tkinter book</a></li>
  <li style="list-style-type:square">Hitting the run button after each widget addition, edit, or removal is a good way to double-check your GUI as well as examining the current build window each time code is updated</li>
  <li style="list-style-type:square">The ultimate goal is to have the current build be runnable at all times, and to ease the creation of a Python GUI</li>
</ul>



###### Try the tutorial
<ul>
  <li style="list-style-type:square">The included <a href="https://github.com/rootVIII/visipy/blob/master/tutorial.pdf" target="_blank">tutorial</a> is simple, but it can give a better feel for overall application usage</li>
  <li style="list-style-type:square">The tutorial example was made using Kubuntu. Therefore if the tutorial is followed on another system, some widget settings may need to be altered to produce the same result.</li>
</ul>


<hr>
This project was developed on Ubuntu 18.04.4 LTS (Kubuntu)