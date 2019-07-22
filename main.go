package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/go-ini/ini"
	"os"
	"os/exec"
	"sort"
)

var cfg *ini.File
var err error

func main() {
	app := app.New()
	w := app.NewWindow("App")
	w.SetContent(widget.NewVBox(
		Mbox()...
	),
	)
	titile := cfg.Section("").Key("titile").String()
	w.SetTitle(titile)
	w.ShowAndRun()
}

func Mbox()[]fyne.CanvasObject {
	var box []fyne.CanvasObject
	names := cfg.SectionStrings()
	for _,v := range names {
		if v == "DEFAULT"{
			continue
		}
		lab := append(box,ALab(v))
		box = append(lab,ABox(v))
	}
	return box
}

func ABox(name string)*widget.Box  {
	var bt []fyne.CanvasObject
	maphash := cfg.Section(name).KeysHash()
	funcs := make(map[string]func(),0)
	for k,v := range maphash {
		funcs[k]=Lamda(v)
	}
	//再拿出来排序
	keys := make([]string, len(funcs))
	i := 0
	for k, _ := range funcs {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, v := range keys{
		bt = append(bt,widget.NewButton(v, funcs[v]))
	}

	return widget.NewHBox(bt...)

}

func ALab(name string)*widget.Label  {
	return widget.NewLabel(name)
}

func Lamda(v string)func()  {
	return func() {
		cmd := exec.Command("/bin/bash", "-c", v)
		cmd.Start()
	}
}

func init()  {
	var conf string = "config.ini"
	if len(os.Args)>=2 {
		conf = os.Args[1]
	}
	cfg, err = ini.Load(conf)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

}