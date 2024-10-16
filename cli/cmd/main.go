package main

import (
	"github.com/8thgencore/microservice-chat/cli/cmd/root"
	"github.com/8thgencore/microservice-common/pkg/closer"
)

func main() {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	root.Execute()
}
