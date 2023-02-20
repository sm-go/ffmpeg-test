package main

import (
	"context"
	"fmt"
)

func main() {
	p := fmt.Println
	p("Go Context Tutorial")
	ctx := context.Background()
	ctx = enrichContext(ctx)
	doSomething(ctx)
}

func enrichContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "request-id", "12345")
}

func doSomething(ctx context.Context) {
	rid := ctx.Value("request-id")
	fmt.Println(rid)
}
