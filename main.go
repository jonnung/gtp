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

// Show usage text when incoming unknown command.
func helpText() string {
	return "usage: gtp [{number}|list|add|remove|clear]\n\n" +
		"  {number}  Show time based one-time password by specified secret\n" +
		"  list      All registered OTP secrets\n" +
		"  add       Add new OTP secret\n" +
		"  remove    Remove the specified secret\n" +
		"  clear     Clear all secrets\n"
}

func errorOutOfRangeText(index int) string {
	return fmt.Sprintf("ERRORRRRR: Selected number %d is out of range OTP list\n", index)
}

// Create linked text with OTP list.
func getTotalOtpsInfo(otpList []otp) []string {
	listResult := make([]string, len(otpList))
	for seq, otp := range otpList {
		listResult[seq] = fmt.Sprintf("{%d} %s:%s:<secret>", seq+1, otp.Issuer, otp.AccountName)
	}
	return listResult
}

// Write gtp configuration file (overwrite).
func rewriteConfFile(file *os.File, otpList []otp) {
	var jsonConfig []byte
	if len(otpList) > 0 {
		jsonConfig, _ = json.Marshal(otpList)
	} else {
		jsonConfig = []byte("")
	}

	ioutil.WriteFile(file.Name(), jsonConfig, os.FileMode(0644))
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

	if len(otpList) == 0 && command != "add" {
		fmt.Println("¯\\_(ツ)_/¯ Nothing has been registered")
		return
	}

	switch command {
	case "list":
		fmt.Println(strings.Join(getTotalOtpsInfo(otpList), "\n"))

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

		rewriteConfFile(file, otpList)

		fmt.Println("✨ 🔑 ✨ Completed the addition of new OTP ")

	case "remove":
		fmt.Println(strings.Join(getTotalOtpsInfo(otpList), "\n"))
		fmt.Print("\nChoose remove target OTP: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		removedSeq, err := strconv.Atoi(scanner.Text())

		if err != nil || removedSeq > len(otpList) {
			fmt.Println(errorOutOfRangeText(removedSeq))
		} else {
			newOtpList := append(otpList[:removedSeq-1], otpList[removedSeq:]...)
			rewriteConfFile(file, newOtpList)
			fmt.Println("💨 Removal success")
		}

	case "clear":
		fmt.Println(strings.Join(getTotalOtpsInfo(otpList), "\n"))
		fmt.Print("\nDo you want clear all?: [y|N]")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		confirmFlag := scanner.Text()

		if confirmFlag == "y" || confirmFlag == "Y" {
			rewriteConfFile(file, []otp{})
			fmt.Println("🗑 Removed all")
		}

	default:
		if otpSeq, err := strconv.Atoi(command); err != nil {
			fmt.Println(helpText())
		} else {
			if otpSeq < 1 || otpSeq > len(otpList) {
				fmt.Printf(errorOutOfRangeText(otpSeq))
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
