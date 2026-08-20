//go:build ignore

package main

import "github.com/gongt/go/pkg/logger"

func main() {
	logger.DLogF("test", "wow, such %s", "doge")
}
