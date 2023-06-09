package core

import (
	"game-engine/util"
	"reflect"
	"strings"
)

type EntityAttribute int64

const (
	IsActive         EntityAttribute = 0
	NeedDestroy                      = 1
	IsTransformDirty                 = 2
)

type Entity struct {
	name             string
	tag              string
	id               int
	components       []Component
	updateComponents []UpdatableComponent
	active           bool
	attributes       *util.Bitset64
	level            Level

	transform *Transform
	collider  *Collider
	vision    *Vision
	network   *NetworkSession
}

var ids int = 0

func NewEntity(level Level, name, tag string) *Entity {
	ids++
	return &Entity{
		attributes: util.NewBitset(3),
		level:      level,
		id:         ids,
		name:       name,
		tag:        tag,
	}
}

//func (e *Entity) SetParent() *Transform {
//	return e.transform
//}
//
//func (e *Entity) ParentTransform() *Transform {
//	return e.parent
//}

func (e *Entity) NetworkSession() *NetworkSession {
	return e.network
}

func (e *Entity) Transform() *Transform {
	return e.transform
}

func (e *Entity) Collider() *Collider {
	return e.collider
}

func (e *Entity) Vision() *Vision {
	return e.vision
}

func (e *Entity) Active() bool {
	return e.active
}

func (e *Entity) SetAttribute(attr EntityAttribute, flag bool) {
	if flag {
		e.attributes.Set(uint64(attr))
	} else {
		e.attributes.Clear(uint64(attr))
	}
}

func (e *Entity) Attribute(attr EntityAttribute) bool {
	return e.attributes.Test(uint64(attr))
}

func (e *Entity) AddComponent(comp Component) Component {
	switch comp.(type) {
	case UpdatableComponent:
		e.updateComponents = append(e.updateComponents, comp.(UpdatableComponent))
	}

	switch comp.(type) {
	case *Collider:
		e.collider = comp.(*Collider)
		break
	case *Transform:
		e.transform = comp.(*Transform)
		break
	case *Vision:
		e.vision = comp.(*Vision)
		break
	case *NetworkSession:
		e.network = comp.(*NetworkSession)
		break
	}

	e.components = append(e.components, comp)

	names := strings.Split(reflect.TypeOf(comp).String(), ".")
	comp.SetName(names[len(names)-1])
	comp.SetEntity(e)

	return comp
}

func (e *Entity) FindComponent(name string) Component {
	for _, comp := range e.components {
		if comp.Name() == name {
			return comp
		}
	}
	return nil
}

func (e *Entity) FindComponentTag(tag string) Component {
	for _, comp := range e.components {
		if comp.Tag() == tag {
			return comp
		}
	}
	return nil
}

func (e *Entity) SetActive(enable bool) {
	e.active = enable
}

func (e *Entity) SetName(name string) {
	e.name = name
}

func (e *Entity) SetId(id int) {
	e.id = id
}

func (e *Entity) SetTag(tag string) {
	e.tag = tag
}

func (e *Entity) Tag() string {
	return e.tag
}

func (e *Entity) Level() Level {
	return e.level
}

func (e *Entity) SetLevel(level Level) Level {
	e.level = level
	return e.level
}

func (e *Entity) Id() int {
	return e.id
}

func (e *Entity) Name() string {
	return e.name
}

//	func (e *Entity) Awake() {
//		e.awake()
//	}
func (e *Entity) awake() {
	for _, comp := range e.components {
		comp.Awake()
		comp.SetEnable(true)
	}
	e.active = true
}

func (e *Entity) update(dt int) {
	if !e.active {
		return
	}
	for _, comp := range e.updateComponents {
		comp.Update(dt)
	}
}

func (e *Entity) finalUpdate(dt int) {
	if !e.active {
		return
	}
	for _, comp := range e.updateComponents {
		comp.FinalUpdate(dt)
	}
}

//func (e *Entity) Update(dt int) {
//	if !e.active {
//		return
//	}
//	for _, comp := range e.updateComponents {
//		comp.Update(dt)
//	}
//}

func (e *Entity) start() {
	for _, comp := range e.components {
		comp.Start()
	}
}

func (e *Entity) release() {
	for _, comp := range e.components {
		comp.Release()
	}
}
