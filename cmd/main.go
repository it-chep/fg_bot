package main

import (
	"context"

	"fg_bot/internal"
)

func main() {
	ctx := context.Background()
	internal.New(ctx).Run(ctx)
}
