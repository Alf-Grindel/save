package utils

import "testing"

func TestMinDistance(t *testing.T) {
	tag1 := []string{"chengdu", "wuhan", "guangzhou", "beijing"}
	tag2 := []string{"chengdu", "wuhan", "nanjing", "beijing"}
	tag3 := []string{"zhejiang", "wenzhou", "hangzhou", "xian"}
	tag4 := []string{"chengdu", "wuhan", "guangzhou", "beijing", "xian"}

	t.Log(MinDistance(tag1, tag2))
	t.Log(MinDistance(tag1, tag3))
	t.Log(MinDistance(tag1, tag4))
	
}
