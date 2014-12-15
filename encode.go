//go:build encode
// +build encode

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/yimiaoxiehou/ltcodes/lt"
)

// Usage 显示程序使用说明
func Usage() {
	fmt.Println("encode <filename> <blockSize> [seed]")
}

func main() {
	// 参数解析
	flag.Parse()
	if flag.NArg() < 2 || flag.NArg() > 3 {
		Usage()
		return
	}
	filename := flag.Arg(0) // 获取输入文件名
	blockSize, err := strconv.Atoi(flag.Arg(1)) // 获取块大小参数

	if err != nil {
		Usage()
		return
	}
	
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())
	seed := rand.Uint32() // 生成随机种子
	
	// 如果提供了种子参数，使用指定的种子
	if flag.NArg() == 3 {
		tmp, err := strconv.ParseUint(flag.Arg(2), 10, 32)
		if err != nil {
			Usage()
			return
		}
		seed = uint32(tmp)
	}

	// 获取文件信息
	stats, err := os.Lstat(filename)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	fSize := stats.Size() // 文件大小

	// 打开输入文件
	f, err := os.Open(filename)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	// 确保文件关闭
	defer func() {
		if err = f.Close(); err != nil {
			fmt.Errorf(err.Error())
		}
	}()

	// 创建LT编码器
	encoder := lt.NewEncoder(f, uint64(fSize), uint32(blockSize), seed)

	// 持续编码并输出编码块
	for err == nil {
		nextBlock := encoder.NextCodedBlock() // 获取下一个编码块
		_, err = os.Stdout.Write(nextBlock.Pack()) // 写入标准输出
	}
}
