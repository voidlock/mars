package mars

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func dial(t *testing.T) redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		t.Fatal("Could not connect to local redis")
	}

	return c
}

func cleanup(conn redis.Conn) {
	conn.Do("FLUSHDB")
	conn.Close()
}

func TestCreateCore(t *testing.T) {
	c := dial(t)
	defer func() {
		cleanup(c)
	}()

	mars := NewMars(c)
	core := mars.CreateCore(8000)

	key := coreKey(core.Id)

	v, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		panic(err)
	}

	result := Core{}

	if err := redis.ScanStruct(v, &result); err != nil {
		panic(err)
	}

	if core != result {
		t.Error("Expected did not match actual: ", core, result)
	}
}

func TestLookupCore(t *testing.T) {
	c := dial(t)
	defer func() {
		cleanup(c)
	}()

	mars := NewMars(c)
	expected := mars.CreateCore(8000)
	actual, err := mars.LookupCore(expected.Id)
	if err != nil {
		t.Fatal(err)
	}

	if expected != actual {
		t.Error("Expected did not match actual: ", expected, actual)
	}
}

func TestAddWarrior(t *testing.T) {
	c := dial(t)
	// defer func() {
	// 	cleanup(c)
	// }()

	mars := NewMars(c)
	core := mars.CreateCore(8000)

	warrior := mars.AddWarrior("Imp", core.Id)

	if warrior.Token == "" {
		t.Error("Expected a token for warrior")
	}
}
