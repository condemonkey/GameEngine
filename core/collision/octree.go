package collision

import (
	"fmt"
	"game-engine/math64/vector3"
)

// FIXME
var (
	CAPACITY = 500
)

type OCNode struct {
	objects  []Collider
	region   *AABB
	children *[8]OCNode
}

// Insert ...
func (n *OCNode) insert(object Collider) bool {
	// Object Bounds doesn't fit in node region => return false
	if !n.region.Contains(object.AABB(0)) {
		return false
	}

	// Number of objects < CAPACITY and children is nil => add in objects
	if len(n.objects) < CAPACITY && n.children == nil {
		n.objects = append(n.objects, object)
		return true
	}

	// Number of objects >= CAPACITY and children is nil => create children,
	// try to move all objects in children
	// and try to add in children otherwise add in objects
	if len(n.objects) >= CAPACITY && n.children == nil {
		n.split()

		objects := n.objects
		n.objects = []Collider{}

		// Move old objects to children
		for i := range objects {
			n.insert(objects[i])
		}
	}

	// Children isn't nil => try to add in children otherwise add in objects
	if n.children != nil {
		for i := range n.children {
			if n.children[i].insert(object) {
				return true
			}
		}
	}
	n.objects = append(n.objects, object)
	return true
}

func (n *OCNode) remove(object Collider) bool {
	// Object outside Bounds
	aabb := object.AABB(0)
	if !aabb.Intersect(n.region) {
		return false
	}

	for i, o := range n.objects {
		// Found it ? delete it and return true
		// o.Equal(object)
		if o.Id() == object.Id() {
			n.objects = append(n.objects[:i], n.objects[i+1:]...)
			n.merge()
			return true
		}
	}

	// If we couldn't remove in current node objects, let's try in children
	if n.children != nil {
		for i := range n.children {
			if n.children[i].remove(object) {
				n.merge()
				return true
			}
		}
	}
	return false
}

func (n *OCNode) getColliding(aabb *AABB) []Collider {
	// If current node region entirely fit inside desired Bounds,
	// No need to search somewhere else => return all objects
	//n.region.Contains(aabb)
	if aabb.Contains(n.region) {
		return n.getAllObjects()
	}
	var objects []Collider
	// If bounds doesn't intersects with region, no collision here => return empty
	if !n.region.Intersect(aabb) {
		return objects
	}
	// return objects that intersects with bounds and its children's objects
	for _, obj := range n.objects {
		if obj.AABB(0).Intersect(aabb) {
			objects = append(objects, obj)
		}
	}
	// No children ? Stop here
	if n.children == nil {
		return objects
	}
	// Get the colliding children
	for _, c := range n.children {
		objects = append(objects, c.getColliding(aabb)...)
	}
	return objects
}

func (n *OCNode) getAllObjects() []Collider {
	var objects []Collider
	objects = append(objects, n.objects...)
	if n.children == nil {
		return objects
	}
	for _, c := range n.children {
		objects = append(objects, c.getAllObjects()...)
	}
	return objects
}

func (n *OCNode) getObjects() []Collider {
	return n.objects
}

// range is already taken
func (n *OCNode) rang(f func(Collider) bool) {
	for _, o := range n.objects {
		if !f(o) {
			return
		}
	}
	if n.children != nil {
		for _, c := range n.children {
			c.rang(f)
		}
	}
}

/* Merge all children into this node - the opposite of Split.
 * Note: We only have to check one level down since a merge will never happen if the children already have children,
 * since THAT won't happen unless there are already too many objects to merge.
 */
func (n *OCNode) merge() bool {
	totalObjects := len(n.objects)
	if n.children != nil {
		for _, child := range n.children {
			if child.children != nil {
				// If any of the *children* have children, there are definitely too many to merge,
				// or the child would have been merged already
				return false
			}
			totalObjects += len(child.objects)
		}
	}
	if totalObjects > CAPACITY {
		return false
	}

	if n.children != nil {
		for i := range n.children {
			curChild := n.children[i]
			numObjects := len(curChild.objects)
			for j := numObjects - 1; j >= 0; j-- {
				curObj := curChild.objects[j]
				n.objects = append(n.objects, curObj)
			}
		}
		// Remove the child nodes (and the objects in them - they've been added elsewhere now)
		n.children = nil
		return true
	}
	return false
}

func (n *OCNode) move(object Collider, pos vector3.Vector3) bool {
	if !n.remove(object) {
		return false
	}
	object.SetCenter(pos)
	// insert 내 AABB(0)에서 업데이트 된 position을 기준으로 aabb를 생성한다.
	// 아래 로직은 필요 없을 듯
	//s := object.Bounds.GetSize().Times(0.5)
	//object.Bounds.Max.X = newPosition[0] + s.X
	//object.Bounds.Max.Y = newPosition[1] + s.Y
	//object.Bounds.Max.Z = newPosition[2] + s.Z
	//
	//object.Bounds.Min.X = newPosition[0] - s.X
	//object.Bounds.Min.Y = newPosition[1] - s.Y
	//object.Bounds.Min.Z = newPosition[2] - s.Z
	return n.insert(object)
}

// Splits the OCNode into eight children.
func (n *OCNode) split() {
	subAabbs := n.region.OctSplit()
	n.children = &[8]OCNode{}
	for i := range subAabbs {
		n.children[i] = OCNode{region: subAabbs[i]}
	}
}

/* * * * * * * * * * * * * * * * * Debugging * * * * * * * * * * * * * * * * */
func (n *OCNode) getNodes() []OCNode {
	var nodes []OCNode
	nodes = append(nodes, *n)
	if n.children != nil {
		for _, c := range n.children {
			nodes = append(nodes, c.getNodes()...)
		}
	}
	return nodes
}

// GetRegion is used for debugging visualisation outside octree package
func (n *OCNode) GetRegion() *AABB {
	return n.region
}

func (n *OCNode) getHeight() int {
	if n.children == nil {
		return 1
	}
	max := 0
	for _, c := range n.children {
		h := c.getHeight()
		if h > max {
			max = h
		}
	}
	return max + 1
}

func (n *OCNode) getNumberOfNodes() int {
	if n.children == nil {
		return 1
	}
	sum := len(n.children)
	for _, c := range n.children {
		nb := c.getNumberOfNodes()
		sum += nb
	}
	return sum
}

func (n *OCNode) getNumberOfObjects() int {
	if n.children == nil {
		return len(n.objects)
	}
	sum := len(n.objects)
	for _, c := range n.children {
		n := c.getNumberOfObjects()
		sum += n
	}
	return sum
}

func (n *OCNode) toString(verbose bool) string {
	var s string
	s = ",\nobjects: [\n"
	if verbose {
		for _, o := range n.objects {
			s += fmt.Sprintf("%v,\n", o)
		}
	} else {
		s += fmt.Sprintf("%v objects,\n", len(n.objects))
	}
	s += "]\n,children: [\n"
	if verbose {
		if n.children != nil {
			for _, c := range n.children {
				s += fmt.Sprintf("%v,\n", c.toString(verbose))
			}
		}
	} else {
		s += fmt.Sprintf("%v children,\n", len(n.children))
	}
	s += "],\n"
	return s
}

type Octree struct {
	root *OCNode
}

// NewOctree is a Octree constructor for ease of use
func NewOctree(size vector3.Vector3) *Octree {
	return &Octree{
		root: &OCNode{region: NewAABBWithSize(vector3.Zero, size)},
	}
}

// Insert a object in the Octree, TODO: bool or object return?
func (o *Octree) Insert(object Collider) bool {
	return o.root.insert(object)
}

// Move object to a new Bounds, pass a pointer because we want to modify the passed object data
func (o *Octree) Move(object Collider, pos vector3.Vector3) bool {
	return o.root.move(object, pos)
}

// Remove object
func (o *Octree) Remove(object Collider) bool {
	return o.root.remove(object)
}

// GetColliding returns an array of objects that intersect with the specified bounds, if any.
// Otherwise returns an empty array.
func (o *Octree) GetColliding(aabb *AABB) []Collider {
	return o.root.getColliding(aabb)
}

// GetAllObjects return all objects, the returned array is sorted in the DFS order
func (o *Octree) GetAllObjects() []Collider {
	return o.root.getAllObjects()
}

// Range based on https://golang.org/src/sync/map.go?s=9749:9805#L296
// Range calls f sequentially for each object present in the octree.
// If f returns false, range stops the iteration.
func (o *Octree) Range(f func(collider Collider) bool) {
	o.root.rang(f)
}

// Get will try to find a specific object based on an id
//func (o *Octree) Get(id uint64, box *AABB) Collider {
//	objs := o.GetColliding(box)
//	for _, obj := range objs {
//
//		//if id == obj.ID() {
//		//	return &obj
//		//}
//	}
//	return nil
//}

// GetSize returns the size of the Octree (cubic volume)
func (o *Octree) GetSize() int64 {
	s := o.root.region.Size()
	return int64(s.X)
}

// GetNodes flatten all the nodes into an array, the returned array is sorted in the DFS order
func (o *Octree) GetNodes() []OCNode {
	return o.root.getNodes()
}

// getHeight debug function
func (o *Octree) getHeight() int {
	return o.root.getHeight()
}

// getNumberOfNodes debug function
func (o *Octree) getNumberOfNodes() int {
	return o.root.getNumberOfNodes()
}

// getNumberOfObjects debug function
func (o *Octree) getNumberOfObjects() int {
	return o.root.getNumberOfObjects()
}

// getUsage ...
func (o *Octree) getUsage() float64 {
	return float64(o.getNumberOfObjects()) / float64(o.getNumberOfNodes()*CAPACITY)
}

func (o *Octree) toString(verbose bool) string {
	return fmt.Sprintf("Octree: {\n%v\n}", o.root.toString(verbose))
}
