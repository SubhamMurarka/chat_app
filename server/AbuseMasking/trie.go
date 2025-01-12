package AbuseMasking

//implementing trie data structure

type Node struct {
	links [26]*Node
	flag  bool
}

func (n *Node) containsKey(ch byte) bool {
	return n.links[ch-'a'] != nil
}

func (n *Node) get(ch byte) *Node {
	return n.links[ch-'a']
}

func (n *Node) put(ch byte, node *Node) {
	n.links[ch-'a'] = node
}

func (n *Node) isEnd() bool {
	return n.flag
}

func (n *Node) setEnd() {
	n.flag = true
}

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	return &Trie{root: &Node{}}
}

func (t *Trie) insert(word string) {
	node := t.root
	for i := 0; i < len(word); i++ {
		if !node.containsKey(word[i]) {
			node.put(word[i], &Node{})
		}
		node = node.get(word[i])
	}
	node.setEnd()
}

func (t *Trie) search(word string) bool {
	node := t.root
	for i := 0; i < len(word); i++ {
		if !node.containsKey(word[i]) {
			return false
		}
		node = node.get(word[i])
	}
	return node.isEnd()
}
