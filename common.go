package main

import (
    "fmt"
    "os"
    "os/user"
    "io/ioutil"
    "strings"
)

// showError – helper for fast error printing
func showError(err error) {
    if err != nil {
        fmt.Printf("Error: %s\n", err.Error())
        os.Exit(-1)
    }

}


// readUserCredentials - AWS Glacier requires to sign requests. Assuming you have done "aws configure"
func readUserCredentials() map[string]string {
    usr, err := user.Current()
    showError(err)

    content, err := ioutil.ReadFile(usr.HomeDir + "/.aws/credentials")
    showError(err)
    result := map[string]string{}
    keys := []string{"aws_access_key_id", "aws_secret_access_key"}
    for _, line := range strings.Split(string(content), "\n") {
        line = strings.TrimSpace(line)
        for _, key := range keys {
            if strings.Contains(line, key) {
                result[key] = strings.Split(line, " = ")[1]
            }
        }
    }

    content, err = ioutil.ReadFile(usr.HomeDir + "/.aws/config")
    showError(err)
    keys = []string{"region"}
    for _, line := range strings.Split(string(content), "\n") {
        line = strings.TrimSpace(line)
        for _, key := range keys {
            if strings.Contains(line, key) {
                result[key] = strings.Split(line, " = ")[1]
            }
        }
    }

    return result
}