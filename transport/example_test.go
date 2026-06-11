package transport_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
)

func ExampleCodec() {
	codec := &stringCodec{}

	// Decode a request.
	req, err := codec.Decode(context.Background(), strings.NewReader("hello"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(req)

	// Encode a response.
	var buf bytes.Buffer
	if err := codec.Encode(context.Background(), &buf, "world"); err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())

	// Output:
	// hello
	// world
}
