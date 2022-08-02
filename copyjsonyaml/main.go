package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Revolyssup/useful-go-func/copyjsonyaml/pkg"
)

var (
	Filepath string
	RootDir  string
	Filename string
	Key      string
	// KeyVal   []pkg.KeyValPair
)

//Main function to test out stuff.
func main() {
	copier := pkg.NewCopier(Filepath, pkg.JSON, RootDir)
	err := copier.Copy(func(env pkg.Env) (string, []pkg.KeyValPair) {
		filename := Filename

		return filename, []pkg.KeyValPair{
			{Key: pkg.Key(Key),
				Val: env.ParentDirectoryName},
		}
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

const usage = `
	USAGE:
	copyjsonyaml <filepath> <root-dir> <filename> <key-value>...
`

func init() {
	if len(os.Args) < 2 {
		log.Fatal("no args passed: " + usage)
	}
	if os.Args[1] == "" {
		log.Fatal("please pass filepath as first argument " + usage)
	}
	Filepath = os.Args[1]
	if os.Args[2] == "" {
		log.Fatal("please pass root directory as second argument " + usage)
	}
	RootDir = os.Args[2]
	if os.Args[3] == "" {
		log.Fatal("please pass file name as third argument " + usage)
	}

	Filename = os.Args[3]
	if os.Args[4] == "" {
		log.Fatal("please pass file name as third argument " + usage)
	}

	Key = os.Args[4]
	// for _, arg := range os.Args[4:] {
	// 	kv := strings.Split(arg, "=")
	// 	if len(kv) < 2 {
	// 		log.Fatal("please pass key val pair as \"key=value\"" + usage)
	// 	}
	// 	k := kv[0]
	// 	v := kv[1]
	// 	// KeyVal = append(KeyVal, pkg.KeyValPair{
	// 	// 	Key: pkg.Key(k),
	// 	// 	Val: pkg.Val(v),
	// 	// })
	// }
}
