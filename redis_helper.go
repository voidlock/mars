package mars

import "github.com/garyburd/redigo/redis"

type Record interface {
	Key() string
}

type List interface {
	Record
	List() []string
}

type redisHelper struct {
	c redis.Conn
}

func (h redisHelper) exists(r Record) bool {
	exists, _ := redis.Bool(h.c.Do("EXISTS", redis.Args{}.Add(r.Key())...))
	return exists
}

func (h redisHelper) do(cmd RedisCmd) (interface{}, error) {
	return h.c.Do(cmd.Command(), cmd.Args()...)
}

type RedisCmd interface {
	Command() string
	Args() redis.Args
}

type GetDetailCmd struct {
	Record
}

func (c GetDetailCmd) Command() string {
	return "HGETALL"
}

func (c GetDetailCmd) Args() redis.Args {
	return redis.Args{}.Add(c.Key())
}

type SetDetailCmd struct {
	Record
}

func (c SetDetailCmd) Command() string {
	return "HMSET"
}

func (c SetDetailCmd) Args() redis.Args {
	return redis.Args{}.Add(c.Key()).AddFlat(c.Record)
}

type AppendListCmd struct {
	List
}

func (c AppendListCmd) Command() string {
	return "RPUSH"
}

func (c AppendListCmd) Args() redis.Args {
	return redis.Args{}.Add(c.Key()).AddFlat(c.List.List())
}
