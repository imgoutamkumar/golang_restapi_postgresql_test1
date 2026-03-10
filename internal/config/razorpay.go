package config

import (
	razorpay "github.com/razorpay/razorpay-go"
)

var RazorpayClient *razorpay.Client

func InitRazorpay() {
	RazorpayClient = razorpay.NewClient("<YOUR_API_KEY>", "<YOUR_API_SECRET>")
}
