package requests_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/carlmjohnson/requests"
)

func ExampleReplayString() {
	const res = `HTTP/1.1 200 OK

An example response.`

	var s string
	const expected = `An example response.`
	if err := requests.
		URL("http://response.example").
		Client(&http.Client{
			Transport: requests.ReplayString(res),
		}).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println(s == expected)
	// Output:
	// true
}

func ExamplePermitURLTransport() {
	cl := *http.DefaultClient
	// Wrap an existing transport or use nil for http.DefaultTransport
	cl.Transport = requests.PermitURLTransport(cl.Transport, `^http://example\.com/?`)
	var s string
	if err := requests.
		URL("http://example.com").
		Client(&cl).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println(strings.Contains(s, "Example Domain"))

	if err := requests.
		URL("http://unauthorized.example.com").
		Client(&cl).
		ToString(&s).
		Fetch(context.Background()); err != nil {
		fmt.Println(err) // unauthorized subdomain not allowed
	}
	// Output:
	// true
	// Get "http://unauthorized.example.com": requested URL not permitted by regexp: ^http://example\.com/?
}
