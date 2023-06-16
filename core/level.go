package core

import (
	"fmt"
	"game-engine/core/collision"
	"sync"
	"time"
)

type Level interface {
	NewEntity(name, tag string)
	AddEntity(entity *Entity)
	RemoveEntity(entity *Entity)
	RelocateEntity(entity *Entity) bool
	FindEntity(id int) *Entity
	SpaceQuery() collision.TreeQuery
	Create() Level
	Load()
	OnLoad()
	OnDestroy()
	SaveSnapshot()
	Name() string
	setName(name string)
	Run() chan bool
	Stop()
	//setBaseLevel(level *BaseLevel)
}

type BaseLevel struct {
	Level
	entities  map[int]*Entity
	tree      *collision.BVTree
	name      string
	wg        *sync.WaitGroup
	closeChan chan bool
}

func NewLevel[T any]() Level {
	obj := any(new(T))
	level := obj.(Level)
	return level.Create()
}

func NewBaseLevel(name string) *BaseLevel {
	b := &BaseLevel{}
	b.name = name
	b.wg = new(sync.WaitGroup)
	b.closeChan = make(chan bool)

	b.tree = collision.NewBVTree(b.wg)
	b.entities = make(map[int]*Entity)

	return b
}

func (b *BaseLevel) Name() string {
	return b.name
}

func (b *BaseLevel) CreateEntity(name, tag string) *Entity {
	entity := NewEntity(b, name, tag)
	entity.AddComponent(NewSphereCollider(0.5)) // 디폴트.
	entity.AddComponent(NewTransform())
	return entity
}

func (b *BaseLevel) CreateEmptyEntity(name, tag string) *Entity {
	entity := NewEntity(b, name, tag)
	return entity
}

func (b *BaseLevel) FindEntity(id int) *Entity {
	return b.entities[id]
}

func (b *BaseLevel) AddEntity(entity *Entity) {
	// 모든 컴포넌트 초기화 및 활성화
	entity.awake()

	// 공간에 추가
	b.tree.AddCollider(entity.Collider())

	// 완료
	b.entities[entity.Id()] = entity
}

func (b *BaseLevel) RemoveEntity(entity *Entity) {
	delete(b.entities, entity.Id())
	b.tree.RemoveCollider(entity.Collider())
}

func (b *BaseLevel) RelocateEntity(entity *Entity) bool {
	return b.tree.RelocateCollider(entity.Collider())
}

func (b *BaseLevel) SpaceQuery() collision.TreeQuery {
	return b.tree
}

func (b *BaseLevel) Run() chan bool {
	closeChan := make(chan bool)
	go func() {
		tickRate := time.Millisecond * 60 // 0.1초마다

		//end := time.Now()
		ticker := time.NewTicker(tickRate)

		//var now int64
		var delta int64
		start := time.Now()

		for {
			select {
			case <-ticker.C:
				now := time.Now()
				delta = now.Sub(start).Milliseconds()
				start = now

				// 트리 balance 업데이트
				//relocates := b.tree.Update()
				//step1 := int(float64(time.Now().UnixNano()-now) / 1000000)
				//if step1 > 10 {
				//	fmt.Println("tree slow update", step1, relocates)
				//}

				b.update(int(delta))
				step2 := time.Since(now).Milliseconds() //int(float64(time.Now().UnixNano()-now) / 1000000)
				if step2 > 60 {
					fmt.Println("entity slow update", step2)
				}
			case <-b.closeChan:
				closeChan <- true
				return
			}
			//start := time.Now()
			//deltaTick := int(start.Sub(end).Milliseconds())
			//b.update(deltaTick)
			//end = time.Now()
			//
			//updateTick := int(end.Sub(start).Milliseconds())
			//
			//fmt.Println(deltaTick, updateTick)
			//
			//time.Sleep(time.Millisecond * 100)
		}
	}()
	return closeChan
}

func (b *BaseLevel) update(dt int) {
	for _, entity := range b.entities {
		entity.update(dt)
		if entity.Attribute(IsTransformDirty) {
			b.tree.RelocateCollider(entity.Collider())
			entity.SetAttribute(IsTransformDirty, false)
		}
	}

	for _, entity := range b.entities {
		entity.finalUpdate(dt)
	}

	b.wg.Wait()
}

func (b *BaseLevel) Stop() {
	b.closeChan <- true
}

//func (b *BaseLevel) Load() {
//}
//
//func (b *BaseLevel) OnLoad() {
//}
//
//func (b *BaseLevel) OnDestroy() {
//}

func (b *BaseLevel) SaveSnapshot() {
	b.tree.Snapshot()
}
