package http

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func Post(url string, data map[string]interface{}) {

	client := resty.New()
	// POST Map, default is JSON content type. No need to set one
	resp, err := client.R().
		SetBody(data).Post(url)

	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	// fmt.Println("  Status     :", resp.Status())
	// fmt.Println("  Proto      :", resp.Proto())
	// fmt.Println("  Time       :", resp.Time())
	// fmt.Println("  Received At:", resp.ReceivedAt())
	// fmt.Println("  Body       :\n", resp)
	// fmt.Println()
}
