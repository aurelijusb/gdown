package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type file struct {
	Path string
	Size int64
}

type archiveUploaded struct {
	ArchiveId string `json:"archiveId"`
	Checksum  string `json:"checksum"`
	Location  string `json:"location"`
}

type archiveToDownload struct {
	Type        string
	ArchiveId   string
	Description string
}

type summary struct {
	Path               string
	Size               int64
	UploadedAt         time.Time
	JobOutput          json.RawMessage
	DownloadParameters archiveToDownload
}

func showErrorMessage(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(-1)
	}
}

func getElements(folder string) []file {
	folder = strings.TrimRight(folder, "/ ")
	result := []file{}
	files, err := ioutil.ReadDir(folder)
	showErrorMessage(err)
	for _, meta := range files {
		result = append(result, file{folder + "/" + meta.Name(), meta.Size()})
	}
	return result
}

func uploadFiles(vaultName string, files []file, outputFilename string) {
	uploaded := []summary{}

	for _, f := range files {
		var response bytes.Buffer
		cmd := exec.Command("aws", "glacier", "upload-archive", "--account-id", "-", "--vault-name", vaultName, "--body", f.Path, "--archive-description", f.Path)
		cmd.Stdout = &response
		cmd.Stderr = &response
		fmt.Printf("Uploading: %s\n", f.Path)
		err := cmd.Run()
		if err != nil {
			showError(err)
		}

		if response.Len() > 0 {
			fmt.Printf("Response from AWS:\n%s\n", response.Bytes())
		}
		uploadedMeta := archiveUploaded{}
		err = json.Unmarshal(response.Bytes(), &uploadedMeta)
		showError(err)

		now := time.Now()
		uploaded = append(uploaded, summary{
			Path:       f.Path,
			Size:       f.Size,
			UploadedAt: now,
			JobOutput:  json.RawMessage(response.Bytes()),
			DownloadParameters: archiveToDownload{
				Type:        "archive-retrieval",
				ArchiveId:   uploadedMeta.ArchiveId,
				Description: f.Path + " | " + now.String(),
			},
		})
	}

	fmt.Printf("Storing summary: %s\n", outputFilename)
	output, err := json.MarshalIndent(uploaded, " ", " ")
	showError(err)
	err = ioutil.WriteFile(outputFilename, output, 0644)
	showError(err)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage:")
		fmt.Println(" " + os.Args[0] + " FOLDER VAULT_NAME SUMMARY_FILE")
		os.Exit(1)
	}
	folder := os.Args[1]
	vaultName := os.Args[2]
	summaryFile := os.Args[3]

	files := getElements(folder)
	fmt.Printf("Will upload %d files\n", len(files))

	uploadFiles(vaultName, files, summaryFile)
}
