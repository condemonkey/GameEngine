package core

type Level interface {
	NewEntity(name, tag string)
	AddEntity(entity *Entity)
	RemoveEntity(entity *Entity)
	Load()
	OnLoad()
	OnDestroy()
}

type BaseLevel struct {
	Level
	entity   *Entity
	entities map[int]*Entity
}

func NewBaseLevel() *BaseLevel {
	b := &BaseLevel{}
	b.entity = NewEntity(b, "level", "", 0)
	b.entities = make(map[int]*Entity)
	return b
}

func (b *BaseLevel) AddLevelComponent(comp Component) {
	b.entity.AddComponent(NewTransform())
}

func (b *BaseLevel) CreateEntity(name, tag string) *Entity {
	entity := NewEntity(b, name, tag, 0)
	return entity
}

func (b *BaseLevel) AddEntity(entity *Entity) {
	b.entities[entity.Id()] = entity
	entity.awake()
}

func (b *BaseLevel) RemoveEntity(entity *Entity) {
	delete(b.entities, entity.Id())
}

func (b *BaseLevel) Update(dt int) {
	b.entity.update(dt)
	for _, entity := range b.entities {
		if entity.Attribute(IsTransformDirty) {
			entity.SetAttribute(IsTransformDirty, false)
		}
		entity.update(dt)
	}
}

func (b *BaseLevel) Load() {
}

func (b *BaseLevel) OnLoad() {
	b.entity.awake()
}

func (b *BaseLevel) OnDestroy() {
}
