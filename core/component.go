package core

type UpdatableComponent interface {
	Update(dt int)
}

type Component interface {
	Name() string
	SetName(name string)
	SetEntity(entity *Entity)
	SetEnable(enable bool)
	Enable() bool
	Tag() string
	Awake()
	Start()
	Release()
	Entity() *Entity
	Transform() *Transform
}

type BaseComponent struct {
	entity *Entity
	name   string
	tag    string
	enable bool
}

func NewBaseComponent() BaseComponent {
	return BaseComponent{}
}

func (b *BaseComponent) SetEnable(enable bool) {
	b.enable = enable
}

func (b *BaseComponent) Enable() bool {
	return b.enable
}

func (b *BaseComponent) SetEntity(entity *Entity) {
	b.entity = entity
}

func (b *BaseComponent) SetName(name string) {
	b.name = name
}

func (b *BaseComponent) Name() string {
	return b.name
}

func (b *BaseComponent) Tag() string {
	return b.tag
}

func (b *BaseComponent) Awake() {
}

func (b *BaseComponent) Start() {
}

func (b *BaseComponent) Release() {
}

func (b *BaseComponent) Entity() *Entity {
	return b.entity
}

func (b *BaseComponent) Transform() *Transform {
	return b.entity.Transform()
}

func (b *BaseComponent) Collider() *Collider {
	return b.entity.Collider()
}
