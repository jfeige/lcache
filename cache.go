package lcache

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

var (
	ErrNotExists = errors.New("lcache: key not exists!")
	ErrDataType  = errors.New("lcache: wrong data type!")
	ErrArgsCount = errors.New("lcache: wrong args count!")
	ErrOutRange  = errors.New("lcache: index out of range")
)

type Item struct {
	Key          string
	Object       interface{}
	BytesObj     []byte
	Expiration   int64       //过期时长时间（当前秒＋过期时长秒）1540304304
	ExpireTicket *time.Timer //过期定时器
}

type Cache struct {
	sync.RWMutex
	items map[string]Item
}

func NewCache() *Cache {
	cache := new(Cache)
	cache.items = make(map[string]Item)
	return cache
}

//检测对象是否过期
func (this Item) checkExpire(cache *Cache) {
	for {
		select {
		case <-this.ExpireTicket.C:
			cache.delExpired(this)
			return
		}
	}
}

//重置对象过期时间
func (this Item) resetExpire(ttl int64, cache *Cache) {
	if this.ExpireTicket == nil {
		this.ExpireTicket = time.NewTimer(time.Duration(ttl) * time.Second)
		go this.checkExpire(cache)
	} else {
		this.ExpireTicket.Reset(time.Duration(ttl) * time.Second)
	}

}

//仅用于字符串和整形,浮点
func (this *Cache) Set(key string, value interface{}) error {
	this.Lock()
	defer this.Unlock()
	item := Item{
		Key:        key,
		Object:     value,
		Expiration: 0,
	}
	switch reflect.ValueOf(value).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.String, reflect.Float32, reflect.Float64:
		item.Object = value
	default:
		return ErrDataType
	}
	this.items[key] = item
	return nil
}

//仅用于字符串和整形,浮点
func (this *Cache) Setex(key string, value interface{}, ttl int64) error {
	this.Lock()
	defer this.Unlock()
	item := Item{
		Key: key,
	}
	switch reflect.ValueOf(value).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.String, reflect.Float32, reflect.Float64:
		item.Object = value
	default:
		return ErrDataType
	}
	if ttl > 0 {
		item.Expiration = int64(ttl)
		item.ExpireTicket = time.NewTimer(time.Duration(ttl) * time.Second)
		go item.checkExpire(this)
	}
	this.items[key] = item
	return nil
}

//仅用于字符串和整形,浮点
func (this *Cache) Get(key string, tp ...interface{}) (interface{}, error) {
	this.Lock()
	defer this.Unlock()
	item, ok := this.items[key]
	if !ok {
		return nil, ErrNotExists
	}
	switch tp := item.Object.(type) {
	case int, int32, int64, string, float32, float64:
		return tp, nil
	default:
		return nil, ErrDataType
	}

	return nil, nil
}

//用map实现哈希
func (this *Cache) Hmset(key string, args ...interface{}) error {
	this.Lock()
	defer this.Unlock()
	//args: key,value,key,value...
	if len(args)%2 != 0 {
		return ErrArgsCount
	}
	//如果key已经存在，则判断类型是否map[string]interface{}类型
	item, ok := this.items[key]
	if ok {
		switch reflect.TypeOf(item.Object).String() {
		case "map[string]interface {}":
			tmp_map := item.Object.(map[string]interface{})
			for i := 0; i < len(args); i += 2 {
				tmp_key := args[i].(string)
				tmp_value := args[i+1]
				tmp_map[tmp_key] = tmp_value
			}
			item.Object = tmp_map
			this.items[key] = item
			return nil
		default:
			return ErrDataType
		}
	}
	tmp_map := make(map[string]interface{}, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		tmp_key := args[i].(string)
		tmp_value := args[i+1]
		tmp_map[tmp_key] = tmp_value
	}
	item = Item{
		Key:        key,
		Object:     tmp_map,
		Expiration: 0,
	}
	this.items[key] = item

	return nil
}

//取hash数据
func (this *Cache) Hgetall(key string) (map[string]interface{}, error) {
	this.Lock()
	defer this.Unlock()
	item, ok := this.items[key]
	if !ok {
		return nil, ErrNotExists
	}
	switch reflect.TypeOf(item.Object).String() {
	case "map[string]interface {}":
		return item.Object.(map[string]interface{}), nil
	default:
		return nil, ErrDataType
	}
}

//slice列表
func (this *Cache) Zadd(key string, values ...interface{}) error {
	this.Lock()
	defer this.Unlock()
	item, ok := this.items[key]
	if ok {
		switch reflect.TypeOf(item.Object).String() {
		case "[]interface {}":
			tmp_slice := item.Object.([]interface{})
			for _, v := range tmp_slice {
				tmp_slice = append(tmp_slice, v)
			}
			item.Object = tmp_slice
			return nil
		default:
			return ErrDataType
		}
		return nil
	}

	tmp_slice := make([]interface{}, 0)
	for _, v := range values {
		fmt.Println("---", v)
		tmp_slice = append(tmp_slice, v)
	}

	item = Item{
		Key:        key,
		Object:     tmp_slice,
		Expiration: 0,
	}
	this.items[key] = item

	return nil
}

//slice列表取数据(升序) -1表示取全部数据
func (this *Cache) Zrange(key string, start, length int) ([]interface{}, error) {
	this.Lock()
	defer this.Unlock()
	item, ok := this.items[key]
	if !ok {
		return nil, ErrNotExists
	}
	fmt.Println(item.Object)
	switch reflect.TypeOf(item.Object).String() {
	case "[]interface {}":
		ret_slice := make([]interface{}, 0)
		tmp_slice := item.Object.([]interface{})
		end := len(tmp_slice)
		if length > 0 {
			if start > len(tmp_slice) || (start+length) > len(tmp_slice) {
				return nil, ErrOutRange
			}
			end = start + length
		}

		for i := start; i < end; i++ {
			ret_slice = append(ret_slice, tmp_slice[i])
		}
		return ret_slice, nil
	default:
		return nil, ErrNotExists
	}
}

//返回缓存中key的数量
func (this *Cache) Keys() int {
	this.Lock()
	defer this.Unlock()
	return len(this.items)
}

//删除过期的数据
func (this *Cache) delExpired(item Item) {
	this.Lock()
	defer this.Unlock()
	delete(this.items, item.Key)
}

//重置过期时间
func (this *Cache) Expire(key string, ttl int64) {
	this.Lock()
	defer this.Unlock()
	item, ok := this.items[key]
	if ok {
		item.resetExpire(ttl, this)
	}
}
