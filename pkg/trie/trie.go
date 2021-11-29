package trie

// This is an opaque wrapper type to abstract over a ternary search trie.
type Trie struct {
	root *trieNode
}

type trieNode struct {
	val      rune
	terminal bool

	left  *trieNode
	right *trieNode
	next  *trieNode
}

// Insert a key, which is a slice of rune's (chars).
func (t *Trie) Insert(key []rune) {
	t.root = t.root.Insert(key)
}

// Given a prefix slice of rune characters, this returns a collection of all
// full matching strings.
func (t *Trie) Search(prefix []rune) [][]rune {
	prefix_node := t.root.Match(prefix)
	if prefix_node == nil {
		return nil
	} else {
		return prefix_node.Collect(prefix)
	}
}

func (t *trieNode) Insert(key []rune) *trieNode {
	if len(key) == 0 {
		return t
	}

	if t == nil {
		t = &trieNode{
			key[0],
			len(key) == 1,

			nil,
			nil,
			nil,
		}
	}

	if len(key) == 0 {
		t.terminal = true
	} else if t.val == key[0] {
		t.next = t.next.Insert(key[1:])
	} else if t.val < key[0] {
		t.right = t.right.Insert(key)
	} else if t.val > key[0] {
		t.left = t.left.Insert(key)
	}

	return t
}

// Return the trieNode that matches a prefix.
func (t *trieNode) Match(key []rune) *trieNode {
	if t == nil {
		return nil
	}

	if len(key) == 0 {
		return t
	} else if t.val == key[0] {
		return t.next.Match(key[1:])
	} else if t.val < key[0] {
		return t.right.Match(key)
	} else {
		return t.left.Match(key)
	}
}

// Gather all complete children of this key, prepending a given prefix.
func (t *trieNode) Collect(prefix []rune) [][]rune {
	if t == nil {
		return make([][]rune, 0, 0)
	}

	// the capacities here are actually just wild guesses, so benchmarks
	// TODO i guess
	pool := make([][]rune, 0, 32)

	pool = append(pool, t.left.Collect(prefix)...)

	if t.terminal {
		pool = append(pool, append(prefix, t.val))
	}

	pool = append(pool, t.next.Collect(append(prefix, t.val))...)
	pool = append(pool, t.right.Collect(prefix)...)

	return pool
}
