package main

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
)

func main() {
	// Let's create a merge patch from these two documents...
	original := []byte(`{"name": "John", "age": 24, "height": 3.21}`)
	target := []byte(`{"name": "Jane", "age": 24}`)

	patch, err := jsonpatch.CreateMergePatch(original, target)
	if err != nil {
		panic(err)
	}

	// Now lets apply the patch against a different JSON document...

	alternative := []byte(`{"name": "Tina", "age": 28, "height": 3.75}`)
	modifiedAlternative, err := jsonpatch.MergePatch(alternative, patch)

	fmt.Printf("patch document:   %s\n", patch)
	fmt.Printf("updated alternative doc: %s\n", modifiedAlternative)
}