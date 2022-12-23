package algorithm

type Priority interface {
	Less(i, j interface{}) bool
	Equal(i, j interface{}) bool
}

type BSTree struct {
	data  interface{}
	Left  *BSTree
	Right *BSTree
}

// NewBSTNode Attention: data need implement Priority
func NewBSTNode(data interface{}) *BSTree {
	return &BSTree{
		data: data,
	}
}

func (p *BSTree) Insert(newNode *BSTree) *BSTree {
	if newNode == nil {
		return p
	}
	if p == nil {
		return NewBSTNode(newNode.data)
	}

	if p.data == nil || newNode.data == nil {
		panic("can't less nil data")
	}

	if p.data.(Priority).Less(p.data, newNode.data) {
		if p.Right == nil {
			p.Right = newNode
		} else {
			p.Right.Insert(newNode)
		}
	} else {
		if p.Left == nil {
			p.Left = newNode
		} else {
			p.Left.Insert(newNode)
		}
	}
	return p
}

// 这里也要考虑被目标节点的情况，分为这几种
//
// 要删除的节点是叶子节点，不用顾虑直接删除
// 要删除的节点左子树不为空，用前驱节点替换当前节点。
// 要删除的节点右子树不为空，用后继节点替换当前节点。
//
// 其实这里还有一种思路，
//
// 要删除的节点是叶子节点，不用顾虑直接删除
// 左右子树有一个为空，用不为空的子节点替换当前节点。
// 左右子树都不为空，前驱节点或者后续节点替换当前节点。
func (p *BSTree) Remove(newNode *BSTree) *BSTree {
	if p.Left != nil {
		p.Left = p.Left.Remove(newNode)
	}

	if p.Right != nil {
		p.Right = p.Right.Remove(newNode)
	}

	// ATTENTION golang 无法修改 receiver 本身，需要 return
	if p.data.(Priority).Equal(p.data, newNode.data) {
		if p.Left == nil && p.Right == nil {
			return nil
		}

		if p.Left == nil || p.Right == nil {
			if p.Left != nil {
				return p.Left
			} else {
				return p.Right
			}
		}

		if p.Left != nil && p.Right != nil {
			p.Right.Insert(p.Left)
			return p.Right
		}
	}

	return p
}

func (p *BSTree) Has(newNode *BSTree) bool {
	if p == nil {
		return false
	}

	if p.data == nil || newNode.data == nil {
		panic("can't less nil data")
	}

	if p.Left != nil {
		found := p.Left.Has(newNode)
		if found {
			return true
		}
	}

	if p.Right != nil {
		found := p.Right.Has(newNode)
		if found {
			return true
		}
	}

	if p.data.(Priority).Equal(p.data, newNode.data) {
		return true
	}
	return false
}

func (p *BSTree) Minimum() *BSTree {
	if p.Left == nil {
		return p
	} else {
		return p.Left.Minimum()
	}
}

func (p *BSTree) Maximum() *BSTree {
	if p.Right == nil {
		return p
	} else {
		return p.Right.Maximum()
	}
}

func (p *BSTree) Data() interface{} {
	return p.data
}
