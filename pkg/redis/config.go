package redis_db

import (
	"context"
	"crypto/tls"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

type IConnectionPool interface {
	client(index int) *redis.ClusterClient
	uninit()
	checkClients()
	isUseClusterMode() bool
}

// connectionPool manage connection connectionPool for redis
type connectionPool struct {
	hosts          []string
	password       string
	useClusterMode bool
	useTLS         bool
	clients        map[int]*redis.ClusterClient
	mu             sync.Mutex
	ctx            context.Context
}

var pool IConnectionPool
var initPoolOnce sync.Once

func (p *connectionPool) client(index int) *redis.ClusterClient {
	p.mu.Lock()
	defer p.mu.Unlock()

	if c, ok := p.clients[index]; ok {
		return c
	}

	opts := redis.ClusterOptions{
		Addrs:         p.hosts,
		Password:      p.password,
		RouteRandomly: false,
		ReadOnly:      true,
		NewClient: func(opt *redis.Options) *redis.Client {
			if !p.useClusterMode {
				opt.DB = index
			}
			return redis.NewClient(opt)
		},
	}

	if !p.useClusterMode {
		opts.ClusterSlots = func(ctx context.Context) ([]redis.ClusterSlot, error) {
			nodes := []redis.ClusterNode{}

			//FIRST HOST MUST BE PRIMARY, ANOTHER ONE MUST BE READ REPLICA
			for _, host := range p.hosts {
				nodes = append(nodes, redis.ClusterNode{
					Addr: host,
				})
			}

			slots := []redis.ClusterSlot{
				{
					Start: 0,
					End:   16384 + 1, //End ตัวเลขตรงนี้จะต้อง มากกว่าหรือเท่ากับ 16383 มันคือตัวเลขที่ go-redis ใช้ในการหาว่า cmd ครั้งนี้จะไปรันที่ slot ไหน โดยโอกาสที่จะเกิดมีตั้งแต่ 0 ถึง 16383
					Nodes: nodes,
				},
			}
			return slots, nil
		}
	}

	if p.useTLS {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	p.clients[index] = redis.NewClusterClient(&opts)

	return p.clients[index]
}

func (p *connectionPool) uninit() {
	for _, c := range p.clients {
		c.Close()
	}
	for k := range p.clients {
		delete(p.clients, k)
	}

}

// checkClients ping each clients and close if error
func (p *connectionPool) checkClients() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for k, c := range p.clients {
		if err := c.Ping(p.ctx).Err(); err != nil {
			c.Close()
			delete(p.clients, k)
		}
	}

}

func (p *connectionPool) isUseClusterMode() bool {
	return p.useClusterMode
}

// Init initialize redis connection pool, must call this function before use
func Init(hosts []string, password string, useClusterMode bool, useTLS bool) {
	initPoolOnce.Do(func() {
		// //MAKE SURE WE LOAD ONLY 1 HOST IF NOT CLUSTER MODE
		// if !useClusterMode {
		// 	hosts = []string{hosts[0]}
		// }

		pool = &connectionPool{
			hosts:          hosts,
			password:       password,
			useClusterMode: useClusterMode,
			useTLS:         useTLS,
			clients:        make(map[int]*redis.ClusterClient),
			ctx:            context.Background(),
		}
	})

}

// Uninit cleanup
func Uninit() {
	if pool == nil {
		return
	}

	pool.uninit()
}

// Client create and return redis client
func Client(index int) *redis.ClusterClient {
	if pool == nil {
		log.Fatal("redis connection pool not init")
		return nil
	}

	if pool.isUseClusterMode() {
		index = 0
	}

	return pool.client(index)
}
