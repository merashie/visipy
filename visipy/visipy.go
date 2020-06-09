/*

BSD 3-Clause License
Copyright (c) 2020, James Colley All rights reserved. Redistribution
and use in source and binary forms, with or without modification, are
permitted provided that the following conditions are met:

Redistributions of source code must retain the above copyright notice,
this list of conditions and the following disclaimer. Redistributions
in binary form must reproduce the above copyright notice, this list of
conditions and the following disclaimer in the documentation and/or
other materials provided with the distribution. Neither the name of the
copyright holder nor the names of its contributors may be used to
endorse or promote products derived from this software without specific
prior written permission. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT
HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/rootVIII/visipy/control"
	"github.com/rootVIII/visipy/utils"
)

func main() {
	var bootstrap = &utils.Bootstrap{
		IsPython3:  false,
		HaveImgs:   true,
		HaveGUI:    true,
		HaveConfig: true,
	}
	tmp := fmt.Sprintf("visipy-%s", time.Now().Format("20060102150405"))
	bootstrap.TempPath = fmt.Sprintf("/var/tmp/%s/", tmp)
	tmpdir := bootstrap.TempPath[:len(bootstrap.TempPath)-(len(tmp)+1)]

	// Clean up any leftover/old project files in case the
	// application exited early during previous usage.
	contents, err := ioutil.ReadDir(tmpdir)
	if err != nil {
		bootstrap.ErrorExit(fmt.Sprintf("Unable to read:\n%s", tmpdir))
	}
	for _, file := range contents {
		if !strings.Contains(file.Name(), "visipy-") {
			continue
		}
		os.RemoveAll(fmt.Sprintf("%s%s", tmpdir, file.Name()))
	}
	os.Mkdir(bootstrap.TempPath, 0700)

	urls := [3]string{
		"https://raw.githubusercontent.com/visipy/vpy/master/icon.ico",
		"https://raw.githubusercontent.com/visipy/vpy/master/icon.png",
		"https://raw.githubusercontent.com/visipy/vpy/master/icon.gif",
	}

	ch := make(chan struct{})
	for _, url := range urls {
		go bootstrap.GetImgs(url, ch)
	}
	go bootstrap.CreateProjectFile("project.json", bootstrap.GetInitialJSON(), ch)
	go bootstrap.CreateProjectFile("gui.py", utils.GetGui(), ch)
	go bootstrap.CheckPython(ch)
	for i := 0; i < 6; i++ {
		<-ch
	}

	bootstrap.MasterLightOffChecklist()

	var visipy control.Controller
	visipy = &control.AppParser{
		Executable: bootstrap.ExePy,
		VisiPath:   bootstrap.TempPath + "gui.py",
		Project:    bootstrap.TempPath + "project",
	}

	visipy.RunVisipy()
	_ = os.RemoveAll(bootstrap.TempPath)
}
