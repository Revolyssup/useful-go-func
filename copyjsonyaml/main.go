package main

import (
	"fmt"

	"github.com/Revolyssup/useful-go-func/copyjsonyaml/pkg"
)

//Main function to test out stuff.
func main() {
	fp := "./test.json"
	dp := "./"
	copier := pkg.NewCopier(fp, pkg.JSON, dp)
	err := copier.Copy(func(env pkg.Env) (string, []pkg.KeyValPair) {
		filename := "newfile.json"
		pairs := []pkg.KeyValPair{
			{
				Key: "info.x",
				Val: env.ParentDirectoryName,
			},
		}
		return filename, pairs
	})
	if err != nil {
		fmt.Println("BHAIIII: ", err.Error())
	}
}
