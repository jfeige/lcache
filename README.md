# lcache
golang实现一个基于内存的key-value缓存，模仿redis，简单实现了几个命令,最先想用gob，将数据编码存储，但是在解码时，碰到了一个问题，一直没有解决，所以对存储的数据没有经过编码

## 安装:

go get github.com/jfeige/lcache


## golang map

存储采用map来实现

## api列表

* Set(key string, value interface{}) error
* Setex(key string, value interface{}, ttl int64) error
* Get(key string, tp ...interface{}) (interface{}, error)
* Hmset(key string, args ...interface{}) error
* Hgetall(key string) (map[string]interface{}, error)
* Zadd(key string, values ...interface{}) error
* Zrange(key string, start, length int) ([]interface{}, error)
* Keys() int
* Expire(key string, ttl int64)

## 使用方法:

```
  	cache := NewCache()

	//set
	cache.Set("first", "first value")

	cache.Set("second", 100)

	cache.Set("third", 100.10)
	//setex
	cache.Setex("four", "four value", 3)

	cache.Setex("five", 9.5, 3)

	//hash
	user := make(map[string]interface{})
	user["name"] = "Jack"
	user["age"] = 28
	user["address"] = "山东省菏泽曹县财神庙街"

	args := make([]interface{}, 0)
	for k, v := range user {
		args = append(args, k, v)
	}

	cache.Hmset("myuser", args...)

	//取hash
	mymap, err := cache.Hgetall("myuser")

	fmt.Printf("myuser:%v,err:%v\n", mymap, err)


	//列表
	users := make([]interface{},0)
	users = append(users,"Jack")
	users = append(users,"Tom")
	users = append(users,"Lucy")

	cache.Zadd("myusers",users...)

	cache.Expire("myusers",3)

	//取列表
	myslice,err := cache.Zrange("myusers",0,-1)

	fmt.Printf("myusers:%v,err:%v\n", myslice, err)


	fmt.Printf("当前缓存key的数量:%d\n",cache.Keys())

	//过期测试,休眠5秒钟
	time.Sleep(5*time.Second)


	four,err := cache.Get("four")

	fmt.Printf("four:%v,err:%v\n", four, err)

	five,err := cache.Get("five")

	fmt.Printf("five:%v,err:%v\n", five, err)

	second,err := cache.Get("second")

	fmt.Printf("second:%v,err:%v\n", second, err)


	//列表已过期
	myslice,err = cache.Zrange("myusers",0,-1)

	fmt.Printf("myusers:%v,err:%v\n", myslice, err)

	fmt.Printf("当前缓存key的数量:%d\n",cache.Keys())
  
```
