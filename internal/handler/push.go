package handler

import (
	"fmt"
	"net/http"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/spikycham/feedme/internal/model"
	"github.com/spikycham/feedme/pkg/network"
)

var (
	VapidPublicKey  = ""
	VapidPrivateKey = ""
)

var (
	SubMerchant *webpush.Subscription
	SubCustomer *webpush.Subscription
)

type PushHanlder struct{}

func NewPushHandler() *PushHanlder { return &PushHanlder{} }

type ResponseGetVAPIDPublicKey struct {
	Key string `json:"key"`
}

func (h *PushHanlder) GetVAPIDPublicKey(w http.ResponseWriter, r *http.Request) error {
	if VapidPublicKey == "" {
		network.Error(w, http.StatusInternalServerError)
		return fmt.Errorf("public vapid key not generated yet")
	}
	network.Write(w, &ResponseGetVAPIDPublicKey{
		Key: VapidPublicKey,
	})
	return nil
}

type RequestSubscribePush struct {
	Role         model.UserRole       `json:"role"`
	Subscription webpush.Subscription `json:"subscription"`
}

func (h *PushHanlder) SubscribePush(w http.ResponseWriter, r *http.Request) error {
	var sub RequestSubscribePush
	if err := network.Read(r, &sub); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	if sub.Role == model.UserRoleMerchant {
		SubMerchant = &sub.Subscription
	} else {
		SubCustomer = &sub.Subscription
	}

	network.WriteEmpty(w, http.StatusNoContent)
	return nil
}

func SendPush(sub *webpush.Subscription, payload []byte) error {
	if sub == nil {
		return fmt.Errorf("no subscripbed merchant")
	}

	resp, err := webpush.SendNotification(
		payload,
		sub,
		&webpush.Options{
			Subscriber:      "mailto:450139220@qq.com",
			VAPIDPublicKey:  VapidPublicKey,
			VAPIDPrivateKey: VapidPrivateKey,
			TTL:             60,
		},
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("push failed: %s", resp.Status)
	}

	return nil
}

func GetVAPIDKeys() error {
	private, public, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		return err
	}
	VapidPublicKey = public
	VapidPrivateKey = private
	return nil
}
