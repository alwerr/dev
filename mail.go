package dev

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
)

func SendMail(msg string) {
	ctx := context.TODO()
	client := resend.NewClient("re_7ook6bbe_LrW18KUD5VUfAGUVg6yTpd9Q")

	params := &resend.SendEmailRequest{
		From:    "onboarding@postrive.com",
		To:      []string{"alwer2640@gmail.com"},
		Subject: "hello world",
		Html:    msg,
	}

	sent, err := client.Emails.SendWithContext(ctx, params)

	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Id)
}
