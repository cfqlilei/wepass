package utils

import (
	"sync"
	"testing"
)

/**
 * 雪花ID生成器测试
 * @author 陈凤庆
 * @date 2025-10-01
 */

/**
 * TestSnowflake_NextID 测试生成ID
 */
func TestSnowflake_NextID(t *testing.T) {
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		t.Fatalf("创建雪花ID生成器失败: %v", err)
	}

	// 生成10个ID
	ids := make([]int64, 10)
	for i := 0; i < 10; i++ {
		id, err := sf.NextID()
		if err != nil {
			t.Fatalf("生成ID失败: %v", err)
		}
		ids[i] = id
	}

	// 检查ID是否递增
	for i := 1; i < len(ids); i++ {
		if ids[i] <= ids[i-1] {
			t.Errorf("ID未递增: %d <= %d", ids[i], ids[i-1])
		}
	}
}

/**
 * TestSnowflake_Concurrent 测试并发生成ID
 */
func TestSnowflake_Concurrent(t *testing.T) {
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		t.Fatalf("创建雪花ID生成器失败: %v", err)
	}

	const goroutines = 100
	const idsPerGoroutine = 100

	var wg sync.WaitGroup
	idChan := make(chan int64, goroutines*idsPerGoroutine)

	// 启动多个goroutine并发生成ID
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id, err := sf.NextID()
				if err != nil {
					t.Errorf("生成ID失败: %v", err)
					return
				}
				idChan <- id
			}
		}()
	}

	wg.Wait()
	close(idChan)

	// 检查ID是否唯一
	idMap := make(map[int64]bool)
	for id := range idChan {
		if idMap[id] {
			t.Errorf("发现重复ID: %d", id)
		}
		idMap[id] = true
	}

	expectedCount := goroutines * idsPerGoroutine
	if len(idMap) != expectedCount {
		t.Errorf("ID数量不正确，期望: %d, 实际: %d", expectedCount, len(idMap))
	}
}

/**
 * TestSnowflake_ParseID 测试解析ID
 */
func TestSnowflake_ParseID(t *testing.T) {
	sf, err := NewSnowflake(5, 10)
	if err != nil {
		t.Fatalf("创建雪花ID生成器失败: %v", err)
	}

	id, err := sf.NextID()
	if err != nil {
		t.Fatalf("生成ID失败: %v", err)
	}

	timestamp, datacenterID, workerID, sequence := ParseID(id)

	if datacenterID != 5 {
		t.Errorf("数据中心ID不正确，期望: 5, 实际: %d", datacenterID)
	}

	if workerID != 10 {
		t.Errorf("机器ID不正确，期望: 10, 实际: %d", workerID)
	}

	if timestamp <= 0 {
		t.Errorf("时间戳不正确: %d", timestamp)
	}

	if sequence < 0 {
		t.Errorf("序列号不正确: %d", sequence)
	}
}

/**
 * TestSnowflake_InvalidParams 测试无效参数
 */
func TestSnowflake_InvalidParams(t *testing.T) {
	// 测试无效的数据中心ID
	_, err := NewSnowflake(32, 1)
	if err == nil {
		t.Error("应该返回错误:数据中心ID超出范围")
	}

	// 测试无效的机器ID
	_, err = NewSnowflake(1, 32)
	if err == nil {
		t.Error("应该返回错误:机器ID超出范围")
	}

	// 测试负数
	_, err = NewSnowflake(-1, 1)
	if err == nil {
		t.Error("应该返回错误:数据中心ID为负数")
	}
}

/**
 * TestGlobalSnowflake 测试全局雪花ID生成器
 */
func TestGlobalSnowflake(t *testing.T) {
	err := InitGlobalSnowflake(1, 1)
	if err != nil {
		t.Fatalf("初始化全局雪花ID生成器失败: %v", err)
	}

	// 生成ID (现在返回GUID字符串)
	id1 := GenerateID()
	id2 := GenerateID()

	if len(id1) == 0 {
		t.Errorf("生成的ID无效: %s", id1)
	}

	if id1 == id2 {
		t.Errorf("生成了重复的ID: %s == %s", id1, id2)
	}

	// 验证GUID格式
	if !IsValidGUID(id1) {
		t.Errorf("生成的ID格式无效: %s", id1)
	}

	// 生成字符串ID
	idStr := GenerateIDString()
	if len(idStr) == 0 {
		t.Error("生成的ID字符串为空")
	}

	// 验证GUID格式
	if !IsValidGUID(idStr) {
		t.Errorf("生成的ID字符串格式无效: %s", idStr)
	}
}

/**
 * BenchmarkSnowflake_NextID 性能测试
 */
func BenchmarkSnowflake_NextID(b *testing.B) {
	sf, _ := NewSnowflake(1, 1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = sf.NextID()
	}
}

/**
 * BenchmarkSnowflake_Concurrent 并发性能测试
 */
func BenchmarkSnowflake_Concurrent(b *testing.B) {
	sf, _ := NewSnowflake(1, 1)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = sf.NextID()
		}
	})
}
