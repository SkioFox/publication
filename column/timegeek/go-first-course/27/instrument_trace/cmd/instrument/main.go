package main

import (
	"flag"
	"fmt"
	"github.com/SkioFox/publication/column/timegeek/go-first-course/27/instrument_trace/instrumenter"
	"github.com/SkioFox/publication/column/timegeek/go-first-course/27/instrument_trace/instrumenter/ast"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	wrote bool
)

func init() {
	flag.BoolVar(&wrote, "w", false, "write result to (source) file instead of stdout")
}

func usage() {
	fmt.Println("instrument [-w] xxx.go")
	flag.PrintDefaults()
}

/*
*
instrument 使用标准库的 flag 包实现对命令行参数（这里是 -w）的解析，通过 os.Args 获取待注入的 Go 源文件路径。
在完成对命令行参数个数与值的校验后， instrument 程序声明了一个 instrumenter.Instrumenter 接口类型变量 ins，
然后创建了一个实现了 Instrumenter 接口类型的 ast.instrumenter 类型的实例，并赋值给变量 ins。
*/
func main() {
	fmt.Println("start main")
	fmt.Println(os.Args)
	flag.Usage = usage
	flag.Parse() // 解析命令行参数

	if len(os.Args) < 2 { // 对命令行参数个数进行校验
		usage()
		return
	}

	var file string
	if len(os.Args) == 3 {
		file = os.Args[2]
	}

	if len(os.Args) == 2 {
		file = os.Args[1]
	}
	if filepath.Ext(file) != ".go" {
		usage()
		return
	}

	var ins instrumenter.Instrumenter // 声明instrumenter.Instrumenter接口类型变量
	// 创建以ast方式实现Instrumenter接口的ast.instrumenter实例
	ins = ast.New("github.com/SkioFox/publication/column/timegeek/go-first-course/27/instrument_trace", "trace", "Trace")
	newSrc, err := ins.Instrument(file)
	if err != nil {
		panic(err)
	}

	if newSrc == nil {
		// add nothing to the source file. no change
		fmt.Printf("no trace added for %s\n", file)
		return
	}

	if !wrote {
		fmt.Println(string(newSrc))
		return
	}

	// write to the source file
	if err = ioutil.WriteFile(file, newSrc, 0666); err != nil {
		fmt.Printf("write %s error: %v\n", file, err)
		return
	}
	fmt.Printf("instrument trace for %s ok\n", file)
}
