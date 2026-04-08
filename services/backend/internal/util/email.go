package util

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v3"
)

func SendEmail() {
	ctx := context.TODO()
	client := resend.NewClient("re_H2nZKg2c_HbSxJ17YC3MyZvMTUWks6JtU")

	params := &resend.SendEmailRequest{
		From:    "Test <no-reply@mail.launch-date.com>",
		To:      []string{"elve960520@gmail.com"},
		Subject: "XXXXXXXxxx",
		Html:    "<p>it works!</p>",
	}

	sent, err := client.Emails.SendWithContext(ctx, params)

	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Id)
}
