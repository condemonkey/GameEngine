package collision

import (
	"game-engine/util"
)

//type NodePool *util.Pool[*BVTreeNode]

type BVTreeNode struct {
	pool     *util.PoolNode
	parent   *BVTreeNode
	left     *BVTreeNode
	right    *BVTreeNode
	collider Collider // Collider
	aabb     AABB
}

func (b *BVTree) NewNodeWithEncapsulate(a0 *AABB, a1 *AABB) *BVTreeNode {
	node := b.allocator.Pop()
	node.aabb = *AABBEncapsulate(a0, a1)
	return node
}

func (b *BVTree) NewNodeWithAABB(aabb *AABB) *BVTreeNode {
	node := b.allocator.Pop()
	node.aabb = *aabb
	return node
}

func (b *BVTree) NewNodeWithShape(collider Collider) *BVTreeNode {
	node := b.allocator.Pop()
	node.aabb = *collider.FatAABB()
	node.collider = collider
	return node
}

func (n *BVTreeNode) IsLeaf() bool {
	return n.left == nil
}

func (n *BVTreeNode) ContainsFatter() bool {
	return n.aabb.Contains(n.collider.AABB())
}

func (n *BVTreeNode) UpdateBranchAABB() {
	assert(n.collider == nil && !n.IsLeaf())
	n.aabb = *AABBEncapsulate(&n.left.aabb, &n.right.aabb)
	//aabb = AABB::Encapsulate(left->aabb, right->aabb);
}

func (n *BVTreeNode) UpdateLeafAABB() {
	assert(n.IsLeaf() && n.collider != nil)
	n.aabb = *n.collider.FatAABB()
	//aabb = collider->GetFatAABB();
}

func (n *BVTreeNode) SwapOutChild(oldChild *BVTreeNode, newChild *BVTreeNode) {
	assert(oldChild == n.left || oldChild == n.right)
	if oldChild == n.left {
		n.left = newChild
		n.left.parent = n
	} else {
		n.right = newChild
		n.right.parent = n
	}
}

// 풀링
func (p *BVTreeNode) Node() *util.PoolNode {
	return p.pool
}

func (p *BVTreeNode) SetNode(node *util.PoolNode) {
	p.pool = node
}

// 스레드 세이프 안함
type BVTree struct {
	//	CollisionUtil::ColliderPairSet colliderPairSet;
	//std::unordered_map<class Collider*, BVTreeNode*> colNodeMap;
	root      *BVTreeNode
	allocator *util.Pool[*BVTreeNode]
	colNodes  map[Collider]*BVTreeNode
	cache     *cache
}

func NewBVTree() *BVTree {
	tree := &BVTree{
		allocator: util.NewPool[*BVTreeNode](2048, func() *BVTreeNode { return &BVTreeNode{} }),
		cache:     &cache{},
		colNodes:  make(map[Collider]*BVTreeNode),
	}
	return tree
}

func (b *BVTree) AddCollider(collider Collider) {
	newNode := b.NewNodeWithShape(collider)
	b.colNodes[collider] = newNode
	b.addNode(newNode)
}

func (b *BVTree) RemoveCollider(collider Collider) {
	node := b.colNodes[collider]
	assert(node != nil)
	b.removeNode(node, true)
	delete(b.colNodes, collider)
}

func (b *BVTree) addNode(newNode *BVTreeNode) {
	newAABB := newNode.aabb

	if b.root == nil {
		b.root = newNode
		b.root.parent = nil
	} else {
		cur := b.root
		lefcnt := 0
		rightcnt := 0
		for !cur.IsLeaf() {
			//cnt++
			leftIncrease := AABBEncapsulate(&cur.left.aabb, &newAABB).SurfaceArea() - cur.left.aabb.SurfaceArea()
			rightIncrease := AABBEncapsulate(&cur.right.aabb, &newAABB).SurfaceArea() - cur.right.aabb.SurfaceArea()
			if leftIncrease > rightIncrease {
				cur = cur.right
				rightcnt++
			} else {
				cur = cur.left
				lefcnt++
			}
		}
		//fmt.Println("call addnoe lieaf", rightcnt, lefcnt)
		if cur == b.root {
			// cur is root
			b.root = b.NewNodeWithEncapsulate(&cur.aabb, &newAABB)
			cur.parent = b.root
			newNode.parent = b.root
			b.root.left = cur
			b.root.right = newNode
		} else {
			// cur is actual leaf, convert cur to branch
			newBranch := b.NewNodeWithEncapsulate(&cur.aabb, &newNode.aabb)
			newBranch.parent = cur.parent
			cur.parent.SwapOutChild(cur, newBranch)
			cur.parent = newBranch
			newNode.parent = newBranch
			newBranch.left = cur
			newBranch.right = newNode

			parent := newBranch.parent
			for parent != nil {
				parent.UpdateBranchAABB()
				parent = parent.parent
			}
		}
	}
}

func (b *BVTree) removeNode(node *BVTreeNode, deleteNode bool) {
	assert(node.IsLeaf())

	if node == b.root {
		b.root = nil
	} else if node.parent == b.root {
		var newRoot *BVTreeNode

		if node == b.root.left {
			newRoot = b.root.right
		} else {
			newRoot = b.root.left
		}

		b.allocator.Push(b.root)
		b.root = newRoot
		b.root.parent = nil
	} else {
		parent := node.parent
		grandParent := parent.parent

		assert(grandParent != nil)
		assert(node == parent.left || node == parent.right)

		if node == parent.left {
			grandParent.SwapOutChild(parent, parent.right)
		} else {
			grandParent.SwapOutChild(parent, parent.left)
		}

		b.allocator.Push(parent)

		cur := grandParent
		for cur != nil {
			cur.UpdateBranchAABB()
			cur = cur.parent
		}
	}

	if deleteNode {
		b.allocator.Push(node)
	}
}

func (b *BVTree) Update() {
	if b.root == nil {
		return
	}
	b.cache.reset()

	// fatAABB범위 밖의 collider들을 검출 (리프까지 순회함)
	b.updateNodes(b.root)

	for i := 0; i < b.cache.count; i++ {
		node := b.cache.nodes[i]
		// 삭제
		b.removeNode(node, false)
	}

	for i := 0; i < b.cache.count; i++ {
		node := b.cache.nodes[i]
		// 부모 박스 크기변경
		node.UpdateLeafAABB()
		// 노드 새로 추가
		b.addNode(node)
	}
}

func (b *BVTree) updateNodes(node *BVTreeNode) {
	if node.IsLeaf() {
		if !node.ContainsFatter() {
			b.cache.push(node)
		}
	} else {
		if node.left != nil {
			b.updateNodes(node.left)
		}
		if node.right != nil {
			b.updateNodes(node.right)
		}
	}
}

type cache struct {
	nodes [2048]*BVTreeNode
	count int
}

func (c *cache) reset() {
	c.count = 0
}

func (c *cache) push(node *BVTreeNode) {
	c.nodes[c.count] = node
	c.count++
}

func (c *cache) empty() bool {
	return c.count == 0
}

func (c *cache) front() *BVTreeNode {
	return c.nodes[0]
}
