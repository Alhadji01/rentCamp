package model

import (
	midtrans "github.com/midtrans/midtrans-go"
	"github.com/sirupsen/logrus"
)

type PaymentModelInterface interface {
	InitiatePayment(orderID string, grossAmount int) (*midtrans.SnapResponse, error)
}

type PaymentModel struct {
	midclient midtrans.Client
}

func NewPaymentModel(midclient midtrans.Client) PaymentModelInterface {
	return &PaymentModel{
		midclient: midclient,
	}
}

func (pm *PaymentModel) InitiatePayment(orderID string, grossAmount int) (*midtrans.SnapResponse, error) {
	snapGateway := midtrans.NewSnapGateway(pm.midclient)

	// Buat permintaan pembayaran
	chargeReq := &midtrans.ChargeReq{
		PaymentType: midtrans.PaymentTypeBankTransfer,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:     orderID,
			GrossAmount: grossAmount,
		},
		CustomerDetails: midtrans.CustomerDetails{
			Email: "customer@example.com",
		},
	}

	// Inisiasi pembayaran dengan Midtrans
	snapResp, err := snapGateway.CreateTransaction(chargeReq)
	if err != nil {
		logrus.Error("Payment Model: Error creating transaction:", err)
		return nil, err
	}

	return snapResp, nil
}
