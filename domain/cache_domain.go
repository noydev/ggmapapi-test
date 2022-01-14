package domain

//go:generate mockgen -destination=../mocks/domain/mockCacheDomain.go -package=domain github.com/noydev/ggmapapi-test/domain CacheDomain
type CacheDomain interface {
	Set(key string, value interface{}) error
	Get(key string, dest interface{}) error
}
