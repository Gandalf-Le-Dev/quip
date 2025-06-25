package main

import (
	"fmt"

	"github.com/Gandalf-Le-Dev/quip/internal/pkg/utils/nanoid"
)

func main() {
	// Simple usage
	id, err := nanoid.New()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generated id: %s\n", id)

	// Custom length
	id, err = nanoid.New(5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generated id: %s\n", id)

	// Custom alphabet
	id, err = nanoid.Generate("abcdefg", 10)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generated id: %s\n", id)

	// Custom non ascii alphabet
	id, err = nanoid.Generate("こちんにабдежиклмнは你好喂שלום😯😪🥱😌😛äöüß", 10)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generated id: %s\n", id)

	fmt.Printf("Generated id: %s\n", nanoid.Must())
	fmt.Printf("Generated id: %s\n", nanoid.MustGenerate("🚀💩🦄🤖", 4))
}