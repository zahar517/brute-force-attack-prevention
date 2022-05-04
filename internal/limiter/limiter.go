package limiter

import "time"

const (
	loginBucketName    = "login"
	passwordBucketName = "password"
	ipBucketName       = "ip"
)

type LeakyBucket struct {
	buckets  map[string]*bucket
	interval int64
	ticker   *time.Ticker
}

func New(loginLimit, passwordLimit, ipLimit, interval int64) *LeakyBucket {
	return &LeakyBucket{
		interval: interval,
		buckets: map[string]*bucket{
			loginBucketName:    newBucket(loginLimit),
			passwordBucketName: newBucket(passwordLimit),
			ipBucketName:       newBucket(ipLimit),
		},
	}
}

func (l *LeakyBucket) Start() {
	go l.tick()
}

func (l *LeakyBucket) Stop() {
	if l.ticker != nil {
		l.ticker.Stop()
	}
}

func (l *LeakyBucket) Add(login, password, ip string) bool {
	in := map[string]string{
		loginBucketName:    login,
		passwordBucketName: password,
		ipBucketName:       ip,
	}

	success := true

	for name, value := range in {
		if ok := l.buckets[name].addKey(value); !ok {
			success = false
		}
	}

	return success
}

func (l *LeakyBucket) Reset(login, ip string) {
	in := map[string]string{
		loginBucketName: login,
		ipBucketName:    ip,
	}

	for name, value := range in {
		l.buckets[name].resetKey(value)
	}
}

func (l *LeakyBucket) resetAll() {
	for _, name := range []string{loginBucketName, passwordBucketName, ipBucketName} {
		l.buckets[name].resetBucket()
	}
}

func (l *LeakyBucket) tick() {
	l.ticker = time.NewTicker(time.Second * time.Duration(l.interval))

	for range l.ticker.C {
		l.resetAll()
	}
}
