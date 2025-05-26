package utils

// https://blog.csdn.net/DBC_121/article/details/104198838

func MinDistance(tagList1, tagList2 []string) int {
	n := len(tagList1)
	m := len(tagList2)

	if n*m == 0 {
		return n + m
	}

	// 初始化二维切片 d[n+1][m+1]
	d := make([][]int, n+1)
	for i := range d {
		d[i] = make([]int, m+1)
	}

	// 边界条件初始化
	for i := 0; i <= n; i++ {
		d[i][0] = i
	}
	for j := 0; j <= m; j++ {
		d[0][j] = j
	}

	// 动态规划计算
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			left := d[i-1][j] + 1
			down := d[i][j-1] + 1
			leftDown := d[i-1][j-1]
			if tagList1[i-1] != tagList2[j-1] {
				leftDown += 1
			}
			d[i][j] = min(left, min(down, leftDown))
		}
	}

	return d[n][m]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
