package game

import (
	"sync"
	"unicode/utf8"
)

type Trie struct {
	Root *TrieNode
	Ver  int //禁言库版本号
	Mu   *sync.RWMutex
}

type TrieNode struct {
	Children map[rune]*TrieNode
	End      bool
}

func NewTrie() Trie {
	var r Trie
	r.Root = NewTrieNode()
	r.Mu = new(sync.RWMutex)
	return r
}

func NewTrieNode() *TrieNode {
	n := new(TrieNode)
	n.Children = make(map[rune]*TrieNode)
	return n
}

func (this *Trie) Inster(txt string) {
	if len(txt) < 1 {
		return
	}
	node := this.Root
	key := []rune(txt)
	for i := 0; i < len(key); i++ {
		if _, exists := node.Children[key[i]]; !exists {
			node.Children[key[i]] = NewTrieNode()
		}
		node = node.Children[key[i]]
	}

	node.End = true
}

func (this *Trie) Replace(txt string) string {
	if len(txt) < 1 {
		return txt
	}
	node := this.Root
	key := []rune(txt)
	var chars []rune = nil
	slen := len(key)
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[key[i]]; exists {
			node = node.Children[key[i]]
			for j := i + 1; j < slen; j++ {
				if _, exists := node.Children[key[j]]; exists {
					node = node.Children[key[j]]
					if node.End == true {
						if chars == nil {
							chars = key
						}
						for t := i; t <= j; t++ {
							c, _ := utf8.DecodeRuneInString("*")
							chars[t] = c
						}
						i = j
						node = this.Root
						break
					}
				}
			}
			node = this.Root
		}
	}
	if chars == nil {
		return txt
	} else {
		return string(chars)
	}
}

//判断是否是敏感词
func (this *Trie) IsSensitiveWord(txt string) bool {
	if len(txt) < 1 {
		return true
	}
	this.Mu.RLock()
	defer this.Mu.RUnlock()

	node := this.Root
	key := []rune(txt)
	slen := len(key)
	for i := 0; i < slen; i++ {
		if _, exists := node.Children[key[i]]; exists {
			node = node.Children[key[i]]
			for j := i + 1; j < slen; j++ {
				if _, exists := node.Children[key[j]]; exists {
					node = node.Children[key[j]]
					if node.End == true {
						return true
					}
				}
			}
			node = this.Root
		}
	}
	return false
}
