package main

import (
	"fmt"
	"game-engine/core/collision"
	"game-engine/core/geo"
	"game-engine/math64"
	"game-engine/math64/vector3"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	min, max := -2000, 2000
	ecount := 5000

	tree := collision.NewBVTree()

	var entities []collision.Collider
	for i := 0; i < ecount; i++ {
		collider := collision.NewCollider(geo.NewSphere(0.5)) // 30
		pointx := float64(math64.RandInt(min, max))
		pointy := float64(math64.RandInt(min, max))
		pointz := float64(math64.RandInt(min, max))
		collider.SetPosition(vector3.Vector3{X: pointx, Y: pointy, Z: pointz})
		collider.SetCollisionHandle(func(hit collision.Collider) {
		})
		entities = append(entities, collider)
	}

	collider := collision.NewCollider(geo.NewSphere(200)) // 30
	collider.SetPosition(vector3.Zero)

	entities = append(entities, collider)

	now := time.Now()
	for _, entity := range entities {
		tree.AddCollider(entity)
	}
	fmt.Println(time.Now().Sub(now).Milliseconds())

	for i := 0; i < 30; i++ {
		now = time.Now()
		wg := new(sync.WaitGroup)
		for j := 0; j < 5000; j++ {
			wg.Add(1)
			go func() {
				tree.Intersect(collider)
				wg.Done()
			}()
		}
		wg.Wait()
		fmt.Println("detect time", time.Now().Sub(now).Milliseconds())
		//t.Log("collision count", count)
		time.Sleep(time.Millisecond * 30)
	}
}
