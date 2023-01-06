package util

import (
	"fmt"
	"testing"
)

func TestTrie(t *testing.T) {
	trie := NewTrie()
	trie.Add("黄色", nil)
	trie.Add("绿色", nil)
	trie.Add("蓝色", nil)
	result, str := trie.Check("我们这里有一个黄色的灯泡, 它存在了很久", "***")
	fmt.Printf("result %#v, str:%v\n", result, str)
}
