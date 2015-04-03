package mars

import (
	"fmt"

	"code.google.com/p/go-uuid/uuid"
	"github.com/garyburd/redigo/redis"
)

const (
	coreKeyFormat        = "core:%s:detail"
	warriorKeyFormat     = "warrior:%s:detail"
	warriorListKeyFormat = "core:%s:warriors"
)

type Mars struct {
	r redisHelper
}

type Model interface {
	Key() string
}

type Core struct {
	Id   string `redis:"id"`
	Size int32  `redis:"size"`
}

type Warrior struct {
	Token string `redis:"id"`
	Name  string `redis:"name"`
}

func (c Core) Key() string {
	return fmt.Sprintf(coreKeyFormat, c.Id)
}

func (w Warrior) Key() string {
	return fmt.Sprintf(warriorKeyFormat, w.Token)
}

func NewMars(conn redis.Conn) Mars {
	return Mars{redisHelper{conn}}
}

func coreKey(id string) string {
	return fmt.Sprintf(coreKeyFormat, id)
}

func (m Mars) CreateCore(size int32) Core {
	core := Core{
		Id:   uuid.New(),
		Size: size,
	}

	if _, err := m.r.do(&SetDetailCmd{core}); err != nil {
		panic(err)
	}

	return core
}

func (m Mars) LookupCore(id string) (Core, error) {
	core := Core{
		Id: id,
	}

	v, err := redis.Values(m.r.do(GetDetailCmd{core}))
	if err != nil {
		return core, fmt.Errorf("Core not found")
	}

	if err := redis.ScanStruct(v, &core); err != nil {
		return core, fmt.Errorf("Core corrupted")
	}

	return core, nil
}

type WarriorList struct {
	core Core `redis:"-"`
	list []Warrior
}

func (wl WarriorList) Key() string {
	return fmt.Sprintf(warriorListKeyFormat, wl.core.Id)
}

func (wl WarriorList) List() []string {
	l := make([]string, len(wl.list))
	for i, w := range wl.list {
		l[i] = w.Token
	}
	return l
}

func (m Mars) AddWarrior(name, coreId string) (w Warrior) {
	c := Core{
		Id: coreId,
	}
	if !m.r.exists(&c) {
		fmt.Printf("Could not find core %v\n", c)
		return
	}

	w = Warrior{
		Token: uuid.New(),
		Name:  name,
	}

	if _, err := m.r.do(&SetDetailCmd{&w}); err != nil {
		panic(err)
	}

	wl := WarriorList{
		core: c,
		list: []Warrior{w},
	}

	if _, err := m.r.do(&AppendListCmd{wl}); err != nil {
		panic(err)
	}

	return
}
