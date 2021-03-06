package main

import (
    "os"
    "net/http"
    "io"
    "fmt"
    "time"
    "io/ioutil"
    "encoding/json"

    "github.com/smartystreets/go-aws-auth"
    "strings"
)

// Job – Structure as in "aws glacier describe-job" output
type Job struct {
    CompletionDate time.Time
    VaultARN string
    RetrievalByteRange string
    SHA256TreeHash string
    Action string
    JobDescription string
    ArchiveId string
    JobId string
    StatusMessage string
    StatusCode string
    Completed bool
    SNSTopic string
    Tier string
    ArchiveSHA256TreeHash string
    ArchiveSizeInBytes int
}

// readJobDescription - read file generated from "aws glacier describe-job"
func readJobDescription(fileName string) Job {
    data, err := ioutil.ReadFile(fileName)
    showError(err)
    result := Job{}
    json.Unmarshal(data, &result)
    return result
}

// downloadFile - downloads file from AWS Glacier
func downloadFile(jobFileName string, outputFilename string) (err error) {
    job := readJobDescription(jobFileName)

    credentials := readUserCredentials()
    region := strings.Split(job.VaultARN, ":")[3]
    vaultName :=  strings.Split(job.VaultARN, "/")[1]
    config := struct {
        AccountId string
        Region string
        AccessKeyId string
        SecretAccessKey string
    }{
        AccountId: strings.Split(job.VaultARN, ":")[4],
        Region: region,
        AccessKeyId: credentials["aws_access_key_id"],
        SecretAccessKey: credentials["aws_secret_access_key"],
    }

    client := &http.Client{}
    url := fmt.Sprintf("https://glacier.%s.amazonaws.com/%s/vaults/%s/jobs/%s/output", region, config.AccountId, vaultName, job.JobId)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }
    req.Header.Add("x-amz-glacier-version", "2012-06-01")

    signedReq := awsauth.Sign(req, awsauth.Credentials{
        AccessKeyID: config.AccessKeyId,
        SecretAccessKey: config.SecretAccessKey,
    })
    if signedReq == nil {
        return fmt.Errorf("Bad signedReq: %#v", signedReq)
    }
    resp, err := client.Do(signedReq)

    if err != nil {
        return err
    }
    defer resp.Body.Close()

    out, err := os.Create(outputFilename)
    if err != nil  {
        return err
    }
    defer out.Close()


    buf := make([]byte, 1024*1024*100) // 100MB
    _, err = io.CopyBuffer(out, resp.Body, buf)
    if err != nil  {
        return err
    }

    return nil
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage:")
        fmt.Println(" " + os.Args[0] + " JOB_DESCRIPTION_FILE OUTPUT_FILE")
        os.Exit(1)
    }
    jobDescription := os.Args[1]
    outputFile := os.Args[2]
    fmt.Printf("Loading job description from %s\n", jobDescription)
    fmt.Println("Downloading...")
    err := downloadFile(jobDescription, outputFile)
    showError(err)
    fmt.Printf("Finished. Saved to %s\n", outputFile)
}