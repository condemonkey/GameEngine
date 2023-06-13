package collision

import (
	"bufio"
	"encoding/json"
	"game-engine/math64/vector3"
	"game-engine/util"
	"os"
	"sync"
)

type BVTreeNode struct {
	pool     *util.PoolNode
	parent   *BVTreeNode
	left     *BVTreeNode
	right    *BVTreeNode
	collider *Collider
	fatAabb  *AABB
}

func (n *BVTreeNode) IsLeaf() bool {
	return n.left == nil
}

func (n *BVTreeNode) IsInFatAABB() bool {
	return n.fatAabb.Contains(n.collider.AABB(0))
}

func (n *BVTreeNode) FatAABB() *AABB {
	return n.fatAabb
}

func (n *BVTreeNode) Collider() *Collider {
	return n.collider
}

func (n *BVTreeNode) UpdateBranchAABB() {
	assert(n.collider == nil && !n.IsLeaf())
	n.fatAabb = AABBUnion(n.left.fatAabb, n.right.fatAabb)
	//fatAabb = AABB::Encapsulate(left->fatAabb, right->fatAabb);
}

func (n *BVTreeNode) UpdateLeafAABB() {
	assert(n.IsLeaf() && n.collider != nil)
	n.fatAabb = n.collider.AABB(FatAABBFactor)
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

func (n *BVTreeNode) Reset() {
	n.fatAabb = nil
	n.collider = nil
	n.left = nil
	n.right = nil
	n.parent = nil
}

// 풀링
func (n *BVTreeNode) Node() *util.PoolNode {
	return n.pool
}

func (n *BVTreeNode) SetNode(node *util.PoolNode) {
	n.pool = node
}

type TreeQuery interface {
	IntersectRangeAsync(origin vector3.Vector3, distance float64, result func(hits int))
	IntersectRange(origin vector3.Vector3, distance float64) int
	IntersectRangeCollidersAsync(origin vector3.Vector3, distance float64, result func(colliders []*Collider))
}

// 스레드 세이프 안함
type BVTree struct {
	//	CollisionUtil::ColliderPairSet colliderPairSet;
	//std::unordered_map<class Collider*, BVTreeNode*> colNodeMap;
	root      *BVTreeNode
	nodes     *util.Pool[*BVTreeNode]
	colliders map[*Collider]*BVTreeNode
	cache     *util.Queue[*BVTreeNode]
	wg        *sync.WaitGroup
}

func NewBVTree() *BVTree {
	tree := &BVTree{
		nodes:     util.NewPool[*BVTreeNode](2048, func() *BVTreeNode { return &BVTreeNode{} }),
		cache:     util.NewQueue[*BVTreeNode](),
		colliders: make(map[*Collider]*BVTreeNode),
		wg:        new(sync.WaitGroup),
	}
	return tree
}

func (b *BVTree) NewNode() *BVTreeNode {
	node := b.nodes.Pop()
	node.Reset()
	return node
}

func (b *BVTree) NewNodeWithAABB(aabb *AABB) *BVTreeNode {
	node := b.NewNode()
	node.fatAabb = aabb
	return node
}

func (b *BVTree) NewNodeWithCollider(collider *Collider) *BVTreeNode {
	node := b.NewNode()
	node.fatAabb = collider.AABB(FatAABBFactor)
	node.collider = collider
	return node
}

func (b *BVTree) AddCollider(collider *Collider) {
	node := b.NewNodeWithCollider(collider)
	b.colliders[collider] = node
	b.addNode(node)
}

func (b *BVTree) RemoveCollider(collider *Collider) {
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
			leftIncrease := AABBUnion(cur.left.fatAabb, newAABB).SurfaceArea() - cur.left.fatAabb.SurfaceArea()
			rightIncrease := AABBUnion(cur.right.fatAabb, newAABB).SurfaceArea() - cur.right.fatAabb.SurfaceArea()
			if leftIncrease > rightIncrease {
				cur = cur.right
			} else {
				cur = cur.left
			}
		}
		//fmt.Println("call addnoe lieaf", rightcnt, lefcnt)
		if cur == b.root {
			// cur is root
			b.root = b.NewNodeWithAABB(AABBUnion(cur.fatAabb, newAABB))
			cur.parent = b.root
			newNode.parent = b.root
			b.root.left = cur
			b.root.right = newNode
		} else {
			// cur is actual leaf, convert cur to branch
			newBranch := b.NewNodeWithAABB(AABBUnion(cur.fatAabb, newNode.fatAabb))
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

func (b *BVTree) UpdateCollider(collider *Collider) bool {
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
	//if node == nil || !node.fatAabb.Raycast(ray, maxDistance, hit) {
	//	return false
	//}
	//if node.IsLeaf() {
	//	hitTmp := &RaycastHit{}
	//	if node.collider.IntersectRay(ray, maxDistance, hitTmp) && hitTmp.distance < hit.distance {
	//		return true
	//	}
	//	return false
	//}
	//return b.raycast(node.left, ray, maxDistance, hit) || b.raycast(node.right, ray, maxDistance, hit)
	return true
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

func (b *BVTree) IntersectRange(origin vector3.Vector3, distance float64) int {
	if b.root == nil {
		return 0
	}
	sphere := NewSphereCollider(distance)
	sphere.SetPosition(origin)
	return b.Intersect(sphere)
}

func (b *BVTree) IntersectRangeAsync(origin vector3.Vector3, distance float64, result func(hits int)) {
	if b.root == nil {
		return
	}
	b.wg.Add(1)
	go func(origin vector3.Vector3, distance float64) {
		defer b.wg.Done()
		result(b.IntersectRange(origin, distance))
	}(origin, distance)
}

func (b *BVTree) IntersectRangeCollidersAsync(origin vector3.Vector3, distance float64, result func(colliders []*Collider)) {
	if b.root == nil {
		return
	}
	b.wg.Add(1)
	go func(origin vector3.Vector3, distance float64) {
		defer b.wg.Done()
		sphere := NewSphereCollider(distance)
		sphere.SetPosition(origin)
		var colliders []*Collider
		b.IntersectColliders(sphere, &colliders)
		result(colliders)
	}(origin, distance)
}

func (b *BVTree) IntersectRangeCollider(origin vector3.Vector3, distance float64, result func(colliders []*Collider)) {
	if b.root == nil {
		return
	}
	sphere := NewSphereCollider(distance)
	sphere.SetPosition(origin)
	var colliders []*Collider
	b.IntersectColliders(sphere, &colliders)
	result(colliders)
}

func (b *BVTree) WaitGroup() *sync.WaitGroup {
	return b.wg
}

func (b *BVTree) Intersect(collider *Collider) int {
	if b.root == nil {
		return 0
	}

	aabb := collider.AABB(0)

	var stack [2048]*BVTreeNode
	stack[0] = b.root
	var top int = 1
	colcount := 0

	for top > 0 {
		top--
		node := stack[top]
		fatAabb := node.fatAabb
		if fatAabb.Intersect(aabb) {
			if node.IsLeaf() {
				if node.collider.IntersectShape(collider) {
					colcount++
				}
			} else {
				stack[top] = node.left
				top++
				stack[top] = node.right
				top++
			}
		}
	}
	return colcount
}

func (b *BVTree) IntersectColliders(collider *Collider, results *[]*Collider) {
	if b.root == nil {
		return
	}

	aabb := collider.AABB(0)

	var stack [2048]*BVTreeNode
	stack[0] = b.root
	var top int = 1

	for top > 0 {
		top--
		node := stack[top]
		fatAabb := node.fatAabb
		if fatAabb.Intersect(aabb) {
			if node.IsLeaf() {
				if node.collider.IntersectShape(collider) {
					*results = append(*results, node.collider)
				}
			} else {
				stack[top] = node.left
				top++
				stack[top] = node.right
				top++
			}
		}
	}
}

func (b *BVTree) IntersectQueue(collider *Collider) int {
	//hits := make([]Collider, 0, 256)

	if b.root == nil {
		return 0
	}
	hits := 0

	aabb := collider.AABB(0)

	queue := b.cache
	queue.Push(b.root)

	for !queue.Empty() {
		node := queue.Pop()
		if node.fatAabb.Intersect(aabb) {
			if node.IsLeaf() {
				if node.collider.IntersectShape(collider) {
					hits++
				}
			} else {
				if node.right != nil {
					queue.Push(node.right)
				}
				if node.left != nil {
					queue.Push(node.left)
				}
			}
		}
	}

	return hits
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

func (b *BVTree) Snapshot() {
	snapshots := Snapshots{}

	b.Traverse(func(node *BVTreeNode) {
		aabb := node.FatAABB()
		radius := float64(0)
		if node.IsLeaf() {
			radius = node.Collider().Shape().(*Sphere).Radius()
		}
		snapshots.Snapshots = append(snapshots.Snapshots, Snapshot{
			Min:    aabb.Min(),
			Max:    aabb.Max(),
			Center: aabb.Center(),
			Size:   aabb.Size(),
			IsLeaf: node.IsLeaf(),
			Radius: radius,
		})
	})

	bytes, _ := json.Marshal(snapshots)
	f, err := os.Create("C:\\Users\\kjk83317\\Desktop\\Unitiy\\Interpolation\\Assets\\Saves\\aabbs.json")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	w.WriteString(string(bytes))
	w.Flush()
}

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

//type Stack[T any] struct {
//	nodes [2048]T
//	index int
//}
//
//func (s *Stack[T]) Empty() bool {
//	return len(s.nodes) == 0
//}
//
//func (s *Stack[T]) Push(item T) {
//	s.nodes[s.index] = item
//	s.index++
//}
//
//func (s *Stack[T]) Pop() T {
//	var item T
//	top := len(s.nodes) - 1
//	item, s.nodes = s.nodes[top], s.nodes[:top]
//	return item
//}
