package config

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func MidtransConfig(serverKey string) *snap.Client {
	midtrans.ServerKey = serverKey
	midtrans.Environment = midtrans.Sandbox

	var snapClient snap.Client
	snapClient.New(serverKey, midtrans.Environment)
	snapClient.ServerKey = midtrans.ServerKey
	snapClient.Env = midtrans.Environment

	return &snapClient
}
