package filter

import (
	"bufio"
	"fmt"
	"io"
	"modtest/gostudy/lesson2/mercury/util"
	"os"
)

var (
	trie *util.Trie
)

//初始化敏感词库
func Init(filename string) (err error) {
	trie = util.NewTrie()
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		word, errRet := reader.ReadString('\n')
		if errRet == io.EOF {
			fmt.Println("文件末尾")
			return
		}

		if errRet != nil {
			err = errRet
			fmt.Println("读取文件失败:", err)
			return
		}
		err = trie.Add(word, nil)
		if err != nil {
			fmt.Println("添加节点失败:", err)
			return
		}
	}
}
