package redis

import "context"

type (
	Recommend struct{}
)

func (r Recommend) AddRecommend(ctx context.Context, k string, v interface{}) error {
	return add(rdbRecommend, ctx, k, v)
}

func (r Recommend) ExistRecommend(ctx context.Context, k string) bool {
	return exist(rdbRecommend, ctx, k)
}

func (r Recommend) GetRecommend(ctx context.Context, k string) (string, error) {
	return get(rdbRecommend, ctx, k)
}
