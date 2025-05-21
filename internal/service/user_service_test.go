package service

import (
	"context"
	"github.com/Alf_Grindel/save/conf"
	"github.com/Alf_Grindel/save/internal/dal"
	"github.com/Alf_Grindel/save/internal/model/basic/user"
	"testing"
)

func BenchmarkSearchUserByTags(b *testing.B) {
	conf.LoadConfig()
	dal.Init()

	tags := &user.SearchUserByTagsReq{
		Tags: []string{"shanghai", "guangzhou"},
	}

	ctx := context.Background()

	service := NewUserService(ctx)

	b.Run("in memory", func(b *testing.B) {
		for range b.N {
			service.SearchUserByTags(tags)
		}
	})

	b.Run("in sql", func(b *testing.B) {
		for range b.N {
			service.SearchUserByTagsBySQL(tags)
		}
	})

}
