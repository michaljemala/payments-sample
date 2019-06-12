// +build ignore

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(http.Dir("../../api"), vfsgen.Options{
		PackageName:  "doc",
		VariableName: "openapi",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
