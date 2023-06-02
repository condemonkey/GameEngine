package main

import (
	"game-engine/core"
	"game-engine/core/collision"
	"game-engine/math64"
	"game-engine/math64/vector3"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	g3ncore "github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"math/rand"
	"time"
)

func main() {
	a := app.App()
	scene := g3ncore.NewNode()
	gui.Manager().Set(scene)

	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	camera.NewOrbitControl(cam)

	onResize := func(name string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	rand.Seed(time.Now().UnixNano())

	tree := collision.NewBVTree()

	var entities []*core.Entity

	entitiesCount := 3
	for i := 0; i < entitiesCount; i++ {
		entity := core.NewEntity(nil, "test", "3dx")
		entity.AddComponent(core.NewSphereCollider())
		entity.AddComponent(core.NewTransform())

		pointx := float64(math64.RandInt(0, 10))
		pointy := float64(math64.RandInt(0, 10))
		pointz := float64(math64.RandInt(0, 10))
		entity.Transform().SetPosition(vector3.Vector3{X: pointx, Y: pointy, Z: pointz})

		entities = append(entities, entity)

		entity.Awake()
	}

	for _, entity := range entities {
		tree.AddCollider(entity.Collider())
	}

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	nodes := make(map[*collision.BVTreeNode]*collision.BVTreeNode)
	
	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		for _, entity := range entities {
			entity.Update(int(deltaTime.Milliseconds()))
		}

		tree.Traverse(func(node *collision.BVTreeNode) {
			if node.IsLeaf() {
				t.scene.Add(t.mesh)
			} else {
				aabb := node.FatAABB()
			}

			//snapshots.Snapshots = append(snapshots.Snapshots, collision.Snapshot{
			//	Min:    aabb.Min(),
			//	Max:    aabb.Max(),
			//	Center: aabb.Center(),
			//	Size:   aabb.Size(),
			//	IsLeaf: node.IsLeaf(),
			//	Radius: radius,
			//})
		})

		renderer.Render(scene, cam)
	})
}
