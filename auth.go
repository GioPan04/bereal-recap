package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GioPan04/bereal"
)

const CONFIG_FILE = "config.json"

// Login or restore previous session. If the previuos session is expired then it will be renewed
func InitConfig() *Config {
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

		if time.Now().After(config.RefreshAt.Add(time.Duration(config.Session.Expiration * 1000000000))) {
			fmt.Println("Refreshing tokens...")
			refresh(config)
		}
	}

	return config
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
