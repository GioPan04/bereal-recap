package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GioPan04/bereal"
)

const CONFIG_FILE = "config.json"

func main() {
	var config *Config

	if _, err := os.Stat(CONFIG_FILE); errors.Is(err, os.ErrNotExist) {
		fmt.Println("You haven't logged in yet. Enter your phone number to login.")
		config = login()
		config.Save(CONFIG_FILE)
	} else {
		config, err = LoadConfig(CONFIG_FILE)
		if err != nil {
			panic(err)
		}
	}

	if time.Now().After(config.RefreshAt.Add(time.Duration(config.Session.Expiration * 1000000000))) {
		fmt.Println("Refreshing tokens...")
		refresh(config)
	}

	memories, err := config.Session.GetMemories()
	if err != nil {
		fmt.Println("Fetching memories")
		panic(err)
	}
	res, err := json.Marshal(memories.Data[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
}

func refresh(config *Config) {
	err := config.Session.RefreshSession()
	if err != nil {
		fmt.Println("Error in session refresh")
		panic(err)
	}

	now := time.Now()
	config.RefreshAt = &now

	config.Save(CONFIG_FILE)
}

func login() *Config {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Phone (with +): ")
	phone, _ := reader.ReadString('\n')
	phone = strings.TrimSuffix(phone, "\n")

	otp_session, err := bereal.SendOtp(phone)
	if err != nil {
		panic(err)
	}

	fmt.Print("Enter OTP Code: ")
	otp, _ := reader.ReadString('\n')
	otp = strings.TrimSuffix(otp, "\n")

	session, err := bereal.VerifyOtp(otp, otp_session)
	if err != nil {
		panic(err)
	}

	now := time.Now()

	return &Config{
		Session:   session,
		RefreshAt: &now,
	}
}
