package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/go-ini/ini"
)

var cfg *ini.File
var path string
var err error

func main() {
	_ = os.Setenv("FYNE_FONT", path+"/miniHei.ttf")
	app := app.New()
	w := app.NewWindow("App")
	f, err := os.Open(path + "/icon.png")
	bit, _ := ioutil.ReadAll(f)
	if err == nil {
		w.SetIcon(&fyne.StaticResource{
			StaticName:    "ico,ico",
			StaticContent: bit,
		})
	}

	w.SetContent(widget.NewVBox(
		Mbox()...,
	),
	)
	titile := cfg.Section("").Key("titile").String()
	w.SetTitle(titile)
	w.ShowAndRun()
}

func Mbox() []fyne.CanvasObject {
	var box []fyne.CanvasObject
	names := cfg.SectionStrings()
	for _, v := range names {
		if v == "DEFAULT" {
			continue
		}
		lab := append(box, ALab(v))
		box = append(lab, ABox(v))
	}
	return box
}

func ABox(name string) *widget.Box {
	var bt []fyne.CanvasObject
	keys2 := cfg.Section(name).KeyStrings()
	// 根据key拿出value
	for _, k := range keys2 {
		v := cfg.Section(name).Key(k).String()
		bt = append(bt, widget.NewButton(k, Lamda(v)))
	}

	return widget.NewHBox(bt...)

}

func ALab(name string) *widget.Label {
	return widget.NewLabel(name)
}

func Lamda(v string) func() {
	return func() {
		if runtime.GOOS == "linux" {
			cmd := exec.Command("/bin/bash", "-c", v)
			cmd.Start()
		} else {
			c := exec.Command("bash", "-c", v)
			c.Run()
		}
	}
}

func init() {
	path, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	var conf string = path + "/config.ini"
	cfg, err = ini.Load(conf)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

}
