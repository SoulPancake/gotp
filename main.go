package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type OTP struct {
	OTP string `json:"otp"`
}

func generateOTP(secret string) string {
	b32secret := base32.StdEncoding.EncodeToString([]byte(secret))
	currentTime := time.Now().Unix() / 30
	h := hmac.New(sha1.New, []byte(b32secret))
	_, _ = h.Write([]byte(strconv.FormatInt(currentTime, 10)))
	otp := fmt.Sprintf("%06d", int(h.Sum(nil)[19])&0x7fffffff%1000000)
	return otp
}

func generateOTPHandler(w http.ResponseWriter, r *http.Request) {
	secretKey := os.Getenv("secret")
	otp := generateOTP(secretKey)
	otpResponse := OTP{OTP: otp}
	json.NewEncoder(w).Encode(otpResponse)
}

func main() {
	http.HandleFunc("/otp", generateOTPHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
