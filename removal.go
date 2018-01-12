package removal

import (
	"os"
	"github.com/botherder/go-files"
	log "github.com/sirupsen/logrus"
	"github.com/mattn/go-colorable"
)

type Removal struct {
	Name 			string
	ProcessNames 	[]string
	Paths 			[]string
}

func NewRemoval(name string, processNames, paths []string) *Removal {
	return &Removal{
		Name: name,
		ProcessNames: processNames,
		Paths: paths,
	}
}

func (r *Removal) KillProcesses() {
	log.Info("Looking for associated processes and trying to kill them...")

	for _, processName := range r.ProcessNames {
		pids, err := FindProcessByName(processName)
		if err != nil {
			log.WithFields(log.Fields{
				"name": processName,
			}).Info("Failed to look for processes: ", err.Error())
			continue
		}

		if len(pids) == 0 {
			log.Warning("Have not found any processes with name ", processName)
			continue
		}

		for _, pid := range(pids) {
			log.WithFields(log.Fields{
				"name": processName,
				"pid": pid,
			}).Info("Found a process matching the name")

			err = KillProcessByPid(pid)
			if err != nil {
				log.WithFields(log.Fields{
					"name": processName,
					"pid": pid,
				}).Error("Failed to kill process: ", err.Error())
				continue
			}

			err = WaitForProcessToDie(pid)
			if err != nil {
				continue
			}

			exePath, err := GetProcessExe(pid)

			log.WithFields(log.Fields{
				"name": processName,
				"pid": pid,
				"path": exePath,
			}).Info("Trying to remove executable for process")

			err = os.Remove(exePath)
			if err != nil {
				log.WithFields(log.Fields{
					"name": processName,
					"pid": pid,
					"path": exePath,
				}).Error("Unable to remove executable: ", err.Error())
			}
		}
	}
}

func (r *Removal) RemoveFiles() {
	log.Info("Trying to remove all associated files and folders...")

	for _, path := range r.Paths {
		path = files.ExpandWindows(path)
		if _, err := os.Stat(path); err == nil {
			err = os.RemoveAll(path)
			if err != nil {
				log.WithFields(log.Fields{
					"path": path,
				}).Error("Unable to remove file or folder: ", err.Error())
			} else {
				log.WithFields(log.Fields{
					"path": path,
				}).Info("Removed file or folder")
			}
		} else {
			log.WithFields(log.Fields{
				"path": path,
			}).Warning("The specified path does not appear to exist")
		}
	}
}

func (r *Removal) Run() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())

	log.Info("Launching removal tool for ", r.Name)

	r.KillProcesses()
	r.RemoveFiles()

	log.Info("Completed")

	log.Info("Press Enter to finish ...")
	var b = make([]byte, 1)
	os.Stdin.Read(b)
}
