package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
)

type otp struct {
	Issuer      string
	AccountName string
	Secret      string
}

func help() string {
	return "$ gtp [{number}|list|add|remove|clear]"
}

func main() {
	var otpList []otp
	// todo: Check `.gtplist` file in $HOME directory
	// todo: Parse JSON data from $HOME/.gtplist
	otpListString := []byte(`[{"Issuer": "Sample", "AccountName": "jonnung", "Secret": "VOLFSSWKAUJRINVWNJNV57QZL74Y5627"}]`)
	if err := json.Unmarshal(otpListString, &otpList); err != nil {
		panic(err)
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println(help())
		return
	}
	command := args[1]

	switch command {
	case "list":
		listResult := make([]string, len(otpList))
		for seq, otp := range otpList {
			listResult[seq] = fmt.Sprintf("{%d} %s:%s:<secret>", seq+1, otp.Issuer, otp.AccountName)
		}
		fmt.Println(strings.Join(listResult, "\n"))
	case "add":
		// todo: Add new otp information by Stdin
		fmt.Println("add")
	default:
		if otpSeq, err := strconv.Atoi(command); err != nil {
			fmt.Println(help())
		} else {
			if otpSeq < 1 || otpSeq > len(otpList) {
				panic("Out of range")
			}

			secret := otpList[otpSeq-1].Secret
			passcode, err := totp.GenerateCode(secret, time.Now().UTC())
			if err != nil {
				panic(err)
			}
			fmt.Println(passcode)
		}
	}
}
