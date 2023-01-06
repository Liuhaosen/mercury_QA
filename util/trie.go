package util

import "strings"

type Node struct {
	char   rune           //rune 表示一个utf8字符, 单个节点
	Data   interface{}    //关联敏感词数据
	parent *Node          //节点的父节点.
	Depth  int            //节点深度
	childs map[rune]*Node //当前节点下一层的子节点
	term   bool           //是否敏感词的末节点. true:是, false:否.  如果查询敏感词term=true, 说明命中了完整的敏感词
}

type Trie struct {
	root *Node //根节点
	size int   //大小
}

//节点构造函数
func NewNode() *Node {
	return &Node{
		childs: make(map[rune]*Node, 32),
	}
}

//trie树构造函数
func NewTrie() *Trie {
	return &Trie{
		root: NewNode(),
	}
}

//向trie树添加敏感词
//key : 敏感词, data: 关联敏感词数据
func (p *Trie) Add(key string, data interface{}) (err error) {
	key = strings.TrimSpace(key)
	node := p.root
	runes := []rune(key)
	for _, r := range runes {
		result, ok := node.childs[r]
		if !ok {
			//如果没有该节点, 那么创建新节点, 作为当前节点的子节点
			result = NewNode()
			result.Depth = node.Depth + 1
			result.char = r
			// result.parent = node
			node.childs[r] = result
		}
		node = result
	}
	node.term = true //到当前节点, 命中了完整敏感词
	node.Data = data
	return
}

//查询敏感词节点
func (p *Trie) findNode(key string) (result *Node) {
	node := p.root
	for _, v := range key {
		ret, ok := node.childs[v]
		if !ok {
			//如果子节点没命中,直接返回
			return
		}
		//如果命中, 当前节点就是该节点, 继续循环
		node = ret
	}
	result = node //这里存放的就是敏感词的最后一个字符
	return
}

//查询相同前缀的节点
func (p *Trie) collectNode(node *Node) (result []*Node) {
	if node == nil {
		return
	}

	if node.term {
		//如果查到完整的敏感词, 插入结果返回即可
		result = append(result, node)
		return
	}

	var queue []*Node
	queue = append(queue, node)
	for i := 0; i < len(queue); i++ {
		if queue[i].term {
			//如果是一个完整的敏感词
			result = append(result, queue[i])
			continue
		}

		//循环子节点
		for _, v1 := range queue[i].childs {
			queue = append(queue, v1)
		}
	}
	return
}

//前缀搜索  key:前缀
func (p *Trie) PrefixSearch(key string) (result []*Node) {
	node := p.findNode(key)
	if node == nil {
		return
	}
	result = p.collectNode(node)
	return
}

//检测敏感词
//text = "我们都喜欢色情片" 提交文本
//replace = "***" 替换字符
func (p *Trie) Check(text, replace string) (result string, hit bool) {
	chars := []rune(text)
	if p.root == nil {
		return
	}
	var left []rune
	node := p.root
	start := 0
	for index, v := range chars {
		ret, ok := node.childs[v]
		if !ok {
			left = append(left, chars[start:index+1]...)
			start = index + 1
			node = p.root
			continue
		}

		node = ret
		if ret.term {
			hit = true
			node = p.root
			left = append(left, ([]rune(replace))...)
			start = index + 1
			continue
		}
	}

	result = string(left)
	return
}
