package tarUtil

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Untar(sourcefile, outputDir string) (string, string, error) {

	dirName, fileName := filepath.Split(sourcefile)
	log.Println("File Name " + fileName)
	log.Println("Directory Name " + dirName)
	absoluteFileName := ""
	customerName := ""
	requestID := ""

	file, err := os.Open(sourcefile)

	if err != nil {
		log.Println(err)
		return "", "", err
	}

	defer file.Close()

	var fileReader io.ReadCloser = file

	// just in case we are reading a tar.gz file, add a filter to handle gzipped file
	if strings.HasSuffix(fileName, ".tgz") {
		if fileReader, err = gzip.NewReader(file); err != nil {
			log.Println(err)
			return "", "", err
		}

		absoluteFileName = strings.Split(fileName, ".")[0]
		strArray := strings.Split(absoluteFileName, "_")
		if len(strArray) < 2 {
			log.Println("Invalid file name syntax. it should be customerName_requestId.zip")
			return "", "", errors.New("Invalid file name syntax. it should be customerName_requestId.zip")
		} else {
			customerName = strArray[0]
			requestID = strArray[1]
			log.Println("Customer Name " + customerName)
			log.Println("Request Id " + requestID)
		}
		defer fileReader.Close()
	} else {
		log.Println("Invalid file")
		return "", "", errors.New("Invalid File")
	}
	tarBallReader := tar.NewReader(fileReader)

	// Extracting tarred files

	rootDir := ""
	outputDir = outputDir + "/" + customerName + "/"
	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			return "", "", err
		}

		// get the individual filename and extract to the output directory
		filename := outputDir + header.Name
		if len(rootDir) == 0 {
			rootDir = filename
		}
		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			log.Println("Creating directory :", filename)
			err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer

			if err != nil {
				log.Println(err)
				return "", "", err
			}

		case tar.TypeReg:
			// handle normal file
			//log.Println("Untarring :", filename)
			writer, err := os.Create(filename)

			if err != nil {
				log.Println(err)
				return "", "", err
			}

			io.Copy(writer, tarBallReader)

			err = os.Chmod(filename, os.FileMode(header.Mode))

			if err != nil {
				log.Println(err)
				return "", "", err
			}

			writer.Close()
		default:
			log.Printf("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
	}
	return rootDir, customerName, nil
}
