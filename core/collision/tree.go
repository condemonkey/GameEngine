package collision

import (
	"game-engine/math64/vector3"
	"game-engine/util"
	g3ncore "github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
)

type Graphic struct {
	scene *g3ncore.Node
	mesh  *graphic.Mesh
	box   *graphic.Mesh
}

type BVTreeNode struct {
	pool     *util.PoolNode
	parent   *BVTreeNode
	left     *BVTreeNode
	right    *BVTreeNode
	collider Collider // Collider
	fatAabb  AABB
	graphic  graphic.Graphic
}

func (b *BVTree) NewNodeWithAABB(aabb *AABB) *BVTreeNode {
	node := b.nodes.Pop()
	node.collider = nil
	node.fatAabb = *aabb
	return node
}

func (b *BVTree) NewNodeWithCollider(collider Collider) *BVTreeNode {
	node := b.nodes.Pop()
	node.fatAabb = *collider.FatAABB()
	node.collider = collider
	return node
}

func (n *BVTreeNode) IsLeaf() bool {
	return n.left == nil
}

func (n *BVTreeNode) IsInFatAABB() bool {
	return n.fatAabb.Contains(n.collider.AABB())
}

func (n *BVTreeNode) FatAABB() AABB {
	return n.fatAabb
}

func (n *BVTreeNode) Collider() Collider {
	return n.collider
}

func (n *BVTreeNode) UpdateBranchAABB() {
	assert(n.collider == nil && !n.IsLeaf())
	n.fatAabb = *AABBEncapsulate(&n.left.fatAabb, &n.right.fatAabb)
	//fatAabb = AABB::Encapsulate(left->fatAabb, right->fatAabb);
}

func (n *BVTreeNode) UpdateLeafAABB() {
	assert(n.IsLeaf() && n.collider != nil)
	n.fatAabb = *n.collider.FatAABB()
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
func (n *BVTreeNode) Node() *util.PoolNode {
	return n.pool
}

func (n *BVTreeNode) SetNode(node *util.PoolNode) {
	n.pool = node
}

// 스레드 세이프 안함
type BVTree struct {
	//	CollisionUtil::ColliderPairSet colliderPairSet;
	//std::unordered_map<class Collider*, BVTreeNode*> colNodeMap;
	root      *BVTreeNode
	nodes     *util.Pool[*BVTreeNode]
	colliders map[Collider]*BVTreeNode
	cache     *util.Queue[*BVTreeNode]
}

func NewBVTree() *BVTree {
	tree := &BVTree{
		nodes:     util.NewPool[*BVTreeNode](2048, func() *BVTreeNode { return &BVTreeNode{} }),
		cache:     util.NewQueue[*BVTreeNode](),
		colliders: make(map[Collider]*BVTreeNode),
	}
	return tree
}

func (b *BVTree) AddCollider(collider Collider) {
	node := b.NewNodeWithCollider(collider)
	b.colliders[collider] = node
	b.addNode(node)
}

func (b *BVTree) RemoveCollider(collider Collider) {
	node := b.colliders[collider]
	assert(node != nil)
	b.removeNode(node, true)
	delete(b.colliders, collider)
}

func (b *BVTree) addNode(newNode *BVTreeNode) {
	newAABB := newNode.fatAabb

	if b.root == nil {
		b.root = newNode
		b.root.parent = nil
	} else {
		cur := b.root
		for !cur.IsLeaf() {
			leftIncrease := AABBEncapsulate(&cur.left.fatAabb, &newAABB).SurfaceArea() - cur.left.fatAabb.SurfaceArea()
			rightIncrease := AABBEncapsulate(&cur.right.fatAabb, &newAABB).SurfaceArea() - cur.right.fatAabb.SurfaceArea()
			if leftIncrease > rightIncrease {
				cur = cur.right
			} else {
				cur = cur.left
			}
		}
		//fmt.Println("call addnoe lieaf", rightcnt, lefcnt)
		if cur == b.root {
			// cur is root
			b.root = b.NewNodeWithAABB(AABBEncapsulate(&cur.fatAabb, &newAABB))
			cur.parent = b.root
			newNode.parent = b.root
			b.root.left = cur
			b.root.right = newNode
		} else {
			// cur is actual leaf, convert cur to branch
			newBranch := b.NewNodeWithAABB(AABBEncapsulate(&cur.fatAabb, &newNode.fatAabb))
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

		b.nodes.Push(b.root)
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

		b.nodes.Push(parent)
		cur := grandParent
		for cur != nil {
			cur.UpdateBranchAABB()
			cur = cur.parent
		}
	}

	if deleteNode {
		b.nodes.Push(node)
	}
}

func (b *BVTree) Update() {
	if b.root == nil {
		return
	}

	// 캐시 리셋
	//b.cache.reset()

	// fatAABB범위 밖의 collider들을 검출 (리프까지 순회함)
	//b.updateNodes(b.root)
	//for i := 0; i < b.cache.count; i++ {
	//	node := b.cache.nodes[i]
	//	// 삭제
	//	b.removeNode(node, false)
	//}
	//
	//for i := 0; i < b.cache.count; i++ {
	//	node := b.cache.nodes[i]
	//	// 부모 박스 크기변경
	//	node.UpdateLeafAABB()
	//	// 노드 새로 추가
	//	b.addNode(node)
	//}

	// stack
	var relocates []*BVTreeNode

	b.cache.Push(b.root)
	for !b.cache.Empty() {
		cur := b.cache.Pop()
		if cur.left != nil {
			b.cache.Push(cur.left)
		}
		if cur.right != nil {
			b.cache.Push(cur.right)
		}
		if cur.IsLeaf() && !cur.IsInFatAABB() {
			relocates = append(relocates, cur)
		}
	}

	for _, node := range relocates {
		b.removeNode(node, false)
	}

	for _, node := range relocates {
		node.UpdateLeafAABB()
		b.addNode(node)
	}
}

func (b *BVTree) RelocateCollider(collider Collider) bool {
	node := b.colliders[collider]
	if node == nil {
		panic("")
	}
	if node.IsLeaf() && !node.IsInFatAABB() {
		b.removeNode(node, false)
		node.UpdateLeafAABB()
		b.addNode(node)
	}
	return true
}

func (b *BVTree) Raycast(ray *Ray, maxDistance float64, hit *RaycastHit) bool {
	return b.raycast(b.root, ray, maxDistance, hit)
}

func (b *BVTree) raycast(node *BVTreeNode, ray *Ray, maxDistance float64, hit *RaycastHit) bool {
	if node == nil || !node.fatAabb.Raycast(ray, maxDistance, hit) {
		return false
	}
	if node.IsLeaf() {
		hitTmp := &RaycastHit{}
		if node.collider.Raycast(ray, maxDistance, hitTmp) && hitTmp.distance < hit.distance {
			return true
		}
		return false
	}
	return b.raycast(node.left, ray, maxDistance, hit) || b.raycast(node.right, ray, maxDistance, hit)
}

func (b *BVTree) updateNodes(node *BVTreeNode) {
	if node.IsLeaf() {
		if !node.IsInFatAABB() {
			b.cache.Push(node)
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

// 모든 노드 순회
func (b *BVTree) Traverse(call func(node *BVTreeNode)) {
	if b.root == nil {
		return
	}
	b.cache.Push(b.root)
	for !b.cache.Empty() {
		cur := b.cache.Pop()
		call(cur)
		if cur.left != nil {
			b.cache.Push(cur.left)
		}
		if cur.right != nil {
			b.cache.Push(cur.right)
		}
	}
}

//func (b *BVTree) Snapshot() {
//	snapshots := Snapshots{}
//
//	b.Traverse(func(node *BVTreeNode) {
//		aabb := node.fatAabb
//		radius := float64(0)
//		if node.IsLeaf() {
//			radius = node.collider.(*core.Collider).InternalShape().(*core.Sphere).Radius()
//		}
//		snapshots.Snapshots = append(snapshots.Snapshots, Snapshot{
//			Min:    aabb.min,
//			Max:    aabb.max,
//			Center: aabb.Center(),
//			Size:   aabb.Size(),
//			IsLeaf: node.IsLeaf(),
//			Radius: radius,
//		})
//	})
//
//	bytes, _ := json.Marshal(snapshots)
//	f, err := os.Create("C:\\Users\\kjk83317\\Desktop\\Unitiy\\Interpolation\\Assets\\Saves\\aabbs.json")
//	if err != nil {
//		panic(err)
//	}
//	w := bufio.NewWriter(f)
//	w.WriteString(string(bytes))
//	w.Flush()
//
//	//s.space.SaveNodes("C:\\Users\\kjk83317\\Desktop\\Unitiy\\Interpolation\\Assets\\Saves\\aabbs.json")
//}

type Snapshots struct {
	Snapshots []Snapshot
}

type Snapshot struct {
	Min    vector3.Vector3
	Max    vector3.Vector3
	Center vector3.Vector3
	Size   vector3.Vector3
	Radius float64
	IsLeaf bool
}

//type cache struct {
//	nodes [2048]*BVTreeNode
//	count int
//}
//
//func (c *cache) reset() {
//	c.count = 0
//}
//
//func (c *cache) push(node *BVTreeNode) {
//	c.nodes[c.count] = node
//	c.count++
//}
//
//func (c *cache) empty() bool {
//	return c.count == 0
//}
//
//func (c *cache) front() *BVTreeNode {
//	return c.nodes[0]
//}
