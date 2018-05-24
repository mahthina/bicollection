package main

import (
	"bicollection/collectInfo"
	"bicollection/tarUtil"
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fsnotify/fsnotify"
	Stack "github.com/golang-collections/collections/stack"
)

const inputDir = "/opt/AutomaticLogAnalysis/inputs"
const outputDir = "/opt/AutomaticLogAnalysis/outputs/"

var aggregateTargets = map[string]string{
	"broker":       "logfile.*",
	"client":       "logfile.*",
	"controller":   "logfile.*",
	"eventmgr":     "logfile.*",
	"idaccessmgr":  "logfile.*",
	"inframgr":     "logfile.*",
	"mysql":        "mysqld.log",
	"web_cloudmgr": "apache-tomcat/logs/.*.log",
}

// main
func main() {

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				//log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Print("File create notification received for :", event.Name)
					outputFile, customerName, err := tarUtil.Untar(event.Name, outputDir)
					if err != nil {
						log.Print("Error while processing " + event.Name + " : " + err.Error())
					} else {
						os.Remove(event.Name)
						summaryFileName := outputFile + "/summary.txt"
						diagsFileName := outputFile + "/diags.txt"
						jsonOutputfileName := outputFile + "/data.json"
						collectInfo.CollectInformation(customerName, summaryFileName, diagsFileName, jsonOutputfileName)
						aggLogs(outputFile)
					}
				}
				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(inputDir); err != nil {
		log.Println("ERROR", err)
	}

	<-done
}

func unzip(fileFullName string) (string, error) {

	dirName, fileName := filepath.Split(fileFullName)
	log.Println("File Name " + fileName)
	log.Println("Directory Name " + dirName)
	absoluteFileName := ""
	customerName := ""
	requestID := ""
	if strings.HasSuffix(fileName, ".zip") {
		log.Println("Valid zip file")
		absoluteFileName = strings.Split(fileName, ".")[0]
		strArray := strings.Split(absoluteFileName, "_")
		if len(strArray) < 2 {
			log.Println("Invalid file name syntax. it should be customerName_requestId.zip")
			return "", errors.New("Invalid file name syntax. it should be customerName_requestId.zip")
		} else {
			customerName = strArray[0]
			requestID = strArray[1]
			log.Println("Customer Name " + customerName)
			log.Println("Request Id " + requestID)
		}
	} else {
		log.Println("Invalid file")
		return "", errors.New("Invalid File")
	}
	log.Println("Unzip in progress for " + fileFullName)
	dest := outputDir + "/" + customerName + "/" + requestID
	log.Println("Full Name " + fileFullName)
	r, err := zip.OpenReader(fileFullName)

	if err != nil {
		return "", err
	}

	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, os.ModePerm)

			if err != nil {
				log.Fatal(err)
				return "", err
			}

			f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())

			if err != nil {
				return "", err
			}

			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return "", err
			}
		}
	}
	log.Println("Unzip process completed " + fileFullName)
	log.Println("Files are unzipped successfully on " + dest + "/" + absoluteFileName)
	return dest + "/" + absoluteFileName, nil
}

func aggLogs(outputFiles string) ([]string, error) {
	log.Println("Log aggregation initiated  for " + outputFiles)
	src := outputFiles + "/logs" // can be more flexible
	dest := src + "_agg"

	var filenames []string

	fil, err := os.Open(src)
	if err != nil {
		return filenames, err
	}

	folderNames, _ := fil.Readdirnames(-1)
	for _, folderName := range folderNames {
		if val, ok := aggregateTargets[folderName]; ok {
			curPath := src + "/" + folderName
			destPath := dest + "/" + folderName
			os.MkdirAll(destPath, os.ModePerm)
			stack := new(Stack.Stack)
			stack.Push(curPath)
			if err != nil {
				return filenames, err
			}

			for stack.Len() > 0 {
				pop_path := stack.Pop().(string)
				open_pop_file, err := os.Open(pop_path)
				if err != nil {
					return filenames, err
				}
				file_infos, err := open_pop_file.Readdir(-1)
				if err != nil {
					return filenames, err
				}

				for _, file_info := range file_infos {
					if file_info.IsDir() {
						stack.Push(pop_path + "/" + file_info.Name())
					} else {
						if match, _ := regexp.MatchString(curPath+"/"+val, pop_path+"/"+file_info.Name()); match {
							// copy log file to logs_agg
							in, err := os.Open(pop_path + "/" + file_info.Name())
							if err != nil {
								return filenames, err
							}
							defer in.Close()

							out, err := os.Create(destPath + "/" + file_info.Name())
							if err != nil {
								return filenames, err
							}
							defer out.Close()

							_, err = io.Copy(out, in)
							if err != nil {
								return filenames, err
							}

							filenames = append(filenames, destPath+"/"+file_info.Name())
						}
					}
				}
			}
		}
	}
	log.Println("Log aggregation completed  on " + outputFiles)
	return filenames, nil
}
