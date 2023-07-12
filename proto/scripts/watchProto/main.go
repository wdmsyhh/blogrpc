package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	cmdPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	protoPath := path.Join(cmdPath, "../../")
	notifyCh := make(chan int64)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	WatchDir(watcher, protoPath)
	go watchFiles(watcher, notifyCh)
	go watchTime(notifyCh, protoPath)
	select {}
}

// if the conditions are met, execute the shell script
func execCmd(protoPath string) {
	cmd := exec.Command("/bin/bash", "-e", protoPath+"/genpb")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Execute Command failed:" + err.Error())
		return
	}
	fmt.Println("Execute Command finished.")
}

// handle folder files changed event
func watchFiles(watcher *fsnotify.Watcher, ch chan int64) {
	for {
		select {
		case ev := <-watcher.Events:
			{
				isNotify := false

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					if strings.Contains(ev.Name, ".proto") {
						log.Println("delete : ", ev.Name)
						isNotify = true
						err := watcher.Remove(ev.Name)
						fmt.Printf("remove watch: %s, err: %vn", ev.Name, err)
					}
				}

				if ev.Op&fsnotify.Create == fsnotify.Create {
					log.Println("create : ", ev.Name)
					if strings.Contains(ev.Name, ".proto") {
						isNotify = true
					}
					file, err := os.Stat(ev.Name)
					if err == nil && strings.Contains(file.Name(), ".proto") {
						watcher.Add(ev.Name)
						fmt.Println("add watch : ", ev.Name)
					}
				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					if strings.Contains(ev.Name, ".proto") {
						log.Println("write : ", ev.Name)
						isNotify = true
					}
				}

				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					if "" != ev.Name {
						if strings.Contains(ev.Name, ".proto") {
							log.Println("rename : ", ev.Name)
							isNotify = true
							err := watcher.Remove(ev.Name)
							fmt.Printf("remove watch: %s, err: %vn", ev.Name, err)
						}
					}
				}

				if isNotify {
					ch <- time.Now().Unix()
				}
			}
		case err := <-watcher.Errors:
			{
				log.Println("watcher error : ", err)
				return
			}
		}
	}
}

// if folder event met, execute the shell script
func watchTime(ch chan int64, protoPath string) {
	var timer *time.Timer
	for {
		select {
		case <-ch:
			{
				if nil != timer {
					log.Printf("reset timer")
					timer.Stop()
				}
				timer = time.NewTimer(1 * time.Second)
				go func() {
					<-timer.C
					execCmd(protoPath)
				}()
			}
		}
	}
}

// watch the folder and sub folders
func WatchDir(watcher *fsnotify.Watcher, dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		p, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		err = watcher.Add(p)
		if err != nil {
			return err
		}
		return nil
	})
}
