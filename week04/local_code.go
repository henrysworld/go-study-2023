package week04

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("我也不知发生什么了，反正是跟 code 有关")
)

// 编译器会在编译的时候，把 set_code 的代码放进来这个 luaSetCode 变量里
//
////go:embed lua/set_code.lua
//var luaSetCode string
//
////go:embed lua/verify_code.lua
//var luaVerifyCode string

type ILocalCodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type item struct {
	value      interface{}
	expiration int64
}
type LocalCodeCache struct {
	data  sync.Map
	mutex sync.RWMutex
	ttl   time.Duration
}

// NewLocalCodeCache Go 的最佳实践是返回具体类型
func NewLocalCodeCache() *LocalCodeCache {
	c := &LocalCodeCache{
		ttl: 10 * time.Second,
	}

	go c.cleanup()

	return c
}

func (c *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.setCodeAndGetInt(ctx, "", []string{c.key(biz, phone)}, code)
	if err != nil {
		return err
	}
	switch res {
	case 0:
		// 毫无问题
		return nil
	case -1:
		// 发送太频繁
		return ErrCodeSendTooMany
	//case -2:
	//	return
	default:
		// 系统错误
		return errors.New("系统错误")
	}
}

func (c *LocalCodeCache) setCodeAndGetInt(ctx context.Context, script string, keys []string, args ...interface{}) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	key := keys[0]
	cntKey := key + ":cnt"
	val := args[0]
	v, ok := c.data.Load(key)
	if ok {
		return -2, nil
	}

	//ttl := time.Now().Add(c.ttl).UnixNano()

	ttl := v.(item).expiration - time.Now().UnixMilli()
	if ttl < 540*1000 {
		vX := item{}
		vX.value = val
		vX.expiration = time.Now().Add(10 * time.Minute).UnixMilli()
		c.data.Store(key, vX)

		vC := item{}
		vC.value = 3
		vC.expiration = time.Now().Add(10 * time.Minute).UnixMilli()
		c.data.Store(cntKey, vC)
		return 0, nil
	}

	return -1, nil
}

func (c *LocalCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.verifyAndGetInt(ctx, "", []string{c.key(biz, phone)}, inputCode)
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		// 正常来说，如果频繁出现这个错误，你就要告警，因为有人搞你
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
		//default:
		//	return false, ErrUnknownForCode
	}
	return false, ErrUnknownForCode
}

func (c *LocalCodeCache) verifyAndGetInt(ctx context.Context, script string, keys []string, args ...interface{}) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	key := keys[0]
	expetedCode := args[0]
	cntKey := key + ":cnt"

	cnt, ok := c.data.Load(cntKey)
	if ok {
		return -1, nil
	}
	if cnt.(int32) <= 0 {
		return -1, nil
	}
	code, ok := c.data.Load(key)
	if ok {
		return -1, nil
	}

	if expetedCode == code {
		c.data.Store(cntKey, -1)
		return 0, nil
	}

	curCnt, ok := c.data.Load(cntKey)
	if ok {
		return -1, nil
	}

	tCurCnt := curCnt.(int8)
	tCurCnt--
	c.data.Store(cntKey, tCurCnt)
	return -2, nil
}

func (c *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *LocalCodeCache) cleanup() {
	for {
		time.Sleep(c.ttl)
		c.data.Range(func(key, value interface{}) bool {
			item := value.(*item)
			if time.Now().UnixNano() > item.expiration {
				c.data.Delete(key)
			}
			return true
		})
	}
}

// LocalCodeCache 假如说你要切换这个，你是不是得把 lua 脚本的逻辑，在这里再写一遍？
//type LocalCodeCache struct {
//	client redis.Cmdable
//}
