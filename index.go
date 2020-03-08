package main

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	files, err := Unzip("1.csv.zip", "output-folder")
	if err != nil {
		log.Fatal(err)
	}

	// parseLocation(files[0])

	fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))

}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()

		if err != nil {
			return filenames, err
		}

		type Domain struct {
			DomainName   string `json:"domainName"`
			RegistarName string `json:"registarName"`
		}

		reader := csv.NewReader(rc)
		var domain []Domain
		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}

			if err != nil {
				fmt.Println(err)
				// return
			}
			if len(line) > 1 {
				domain = append(domain, Domain{
					DomainName:   line[0],
					RegistarName: line[1],
				})
			}
		}
		// fmt.Println((domain))
		domainJson, err := json.Marshal(domain)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		fmt.Println(string(domainJson))
		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
