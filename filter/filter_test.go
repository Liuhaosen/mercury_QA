package filter

import (
	"fmt"
	"testing"
)

func TestReplace(t *testing.T) {
	err := Init("../data/filter.dat.txt")
	if err != nil {
		t.Errorf("load filter data failed, err : %#v\n", err)
		return
	}

	result, isReplace := Replace("裸体喜欢小黄片乱伦", "***")
	fmt.Printf("isReplace:%#v, str:%v\n", isReplace, result)
}

func TestAdd(t *testing.T) {
	err := Init("../data/filter.dat.txt")
	if err != nil {
		fmt.Println("初始化失败, err:", err)
		return
	}
	fmt.Println("初始化成功")
}

func TestCheck(t *testing.T) {
	err := Init("../data/filter.dat.txt")
	if err != nil {
		fmt.Println("初始化失败, err:", err)
		return
	}
	fmt.Println("初始化成功")
	result, str := trie.Check("裸体喜欢小黄片乱伦", "***")
	fmt.Println("result :", result, "str : ", str)
}
