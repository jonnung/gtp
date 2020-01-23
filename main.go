package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func helpText() string {
	return "usage: gtp [{number}|list|add|remove|clear]\n\n" +
		"\t{number}\t등록된 OTP의 Passcode를 출력\n" +
		"\tlist\t등록된 OTP 전체 목록\n" +
		"\tadd\t신규 OTP Secret 추가\n" +
		"\tremove\t선택된 OTP 삭제\n" +
		"\tclear\t전체 OTP 삭제\n"
}

func main() {
	var otpList []otp
	var file *os.File

	userHomeDir, _ := os.UserHomeDir()
	configFile := userHomeDir + "/.gtplist"

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Check exists file for practice purpose
		f, _ := os.Create(configFile)
		file = f
	} else {
		f, _ := os.OpenFile(configFile, os.O_RDWR, os.FileMode(0644))
		file = f
	}

	defer file.Close()
	// gtpList, err := ioutil.ReadFile(configFile)
	gtpList, err := ioutil.ReadAll(file)
	if err != nil {
		panic("Error reading configuration file")
	}

	if len(gtpList) > 0 {
		// otpListString := []byte(`[{"Issuer": "Sample", "AccountName": "jonnung", "Secret": "VOLFSSWKAUJRINVWNJNV57QZL74Y5627"}]`)
		if err := json.Unmarshal(gtpList, &otpList); err != nil {
			panic(err)
		}
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println(helpText())
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
		newOtp := otp{}
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Step 1/3) Issuer: ")
		scanner.Scan()
		newOtp.Issuer = scanner.Text()

		fmt.Print("Step 2/3) Account Name: ")
		scanner.Scan()
		newOtp.AccountName = scanner.Text()

		fmt.Print("Step 3/3) Secret: ")
		scanner.Scan()
		newOtp.Secret = scanner.Text()

		otpList = append(otpList, newOtp)

		jsonConfig, _ := json.Marshal(otpList)
		file.WriteAt(jsonConfig, 0)

		fmt.Println("✨ 🔑 ✨ Completed the addition of new OTP ")

	default:
		if otpSeq, err := strconv.Atoi(command); err != nil {
			fmt.Println(helpText())
		} else {
			if otpSeq < 1 || otpSeq > len(otpList) {
				fmt.Printf("ERRORRRRR: Selected number %d is out of range set OTP list.\n\n", otpSeq)
				fmt.Println(helpText())
				return
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
