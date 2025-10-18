package utils

import (
	"fmt"
	"sync"
	"time"
)

/**
 * 雪花ID生成器
 * @author 陈凤庆
 * @date 2025-10-01
 * @description 基于Twitter Snowflake算法的分布式ID生成器
 */

const (
	// 时间戳占用位数
	timestampBits = 41
	// 数据中心ID占用位数
	datacenterBits = 5
	// 机器ID占用位数
	workerBits = 5
	// 序列号占用位数
	sequenceBits = 12

	// 最大值
	maxDatacenterID = -1 ^ (-1 << datacenterBits) // 31
	maxWorkerID     = -1 ^ (-1 << workerBits)     // 31
	maxSequence     = -1 ^ (-1 << sequenceBits)   // 4095

	// 位移量
	workerShift     = sequenceBits                               // 12
	datacenterShift = sequenceBits + workerBits                  // 17
	timestampShift  = sequenceBits + workerBits + datacenterBits // 22

	// 起始时间戳 (2024-01-01 00:00:00 UTC)
	epoch int64 = 1704067200000
)

/**
 * Snowflake 雪花ID生成器
 */
type Snowflake struct {
	mu           sync.Mutex
	timestamp    int64
	datacenterID int64
	workerID     int64
	sequence     int64
}

var (
	// 全局雪花ID生成器实例
	globalSnowflake *Snowflake
	once            sync.Once
)

/**
 * NewSnowflake 创建雪花ID生成器
 * @param datacenterID 数据中心ID (0-31)
 * @param workerID 机器ID (0-31)
 * @return *Snowflake 雪花ID生成器实例
 * @return error 错误信息
 */
func NewSnowflake(datacenterID, workerID int64) (*Snowflake, error) {
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, fmt.Errorf("数据中心ID必须在0-%d之间", maxDatacenterID)
	}
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("机器ID必须在0-%d之间", maxWorkerID)
	}

	return &Snowflake{
		timestamp:    0,
		datacenterID: datacenterID,
		workerID:     workerID,
		sequence:     0,
	}, nil
}

/**
 * InitGlobalSnowflake 初始化全局雪花ID生成器
 * @param datacenterID 数据中心ID (0-31)
 * @param workerID 机器ID (0-31)
 * @return error 错误信息
 */
func InitGlobalSnowflake(datacenterID, workerID int64) error {
	var err error
	once.Do(func() {
		globalSnowflake, err = NewSnowflake(datacenterID, workerID)
	})
	return err
}

/**
 * GetGlobalSnowflake 获取全局雪花ID生成器
 * @return *Snowflake 雪花ID生成器实例
 */
func GetGlobalSnowflake() *Snowflake {
	if globalSnowflake == nil {
		// 如果未初始化,使用默认值初始化
		_ = InitGlobalSnowflake(0, 0)
	}
	return globalSnowflake
}

/**
 * NextID 生成下一个ID
 * @return int64 雪花ID
 * @return error 错误信息
 */
func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano() / 1000000 // 毫秒时间戳

	if now < s.timestamp {
		return 0, fmt.Errorf("时钟回拨,拒绝生成ID")
	}

	if now == s.timestamp {
		// 同一毫秒内,序列号递增
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// 序列号溢出,等待下一毫秒
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 不同毫秒,序列号重置为0
		s.sequence = 0
	}

	s.timestamp = now

	// 组装ID
	id := ((now - epoch) << timestampShift) |
		(s.datacenterID << datacenterShift) |
		(s.workerID << workerShift) |
		s.sequence

	return id, nil
}

/**
 * NextIDString 生成下一个ID的字符串形式
 * @return string 雪花ID字符串
 * @return error 错误信息
 */
func (s *Snowflake) NextIDString() (string, error) {
	id, err := s.NextID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", id), nil
}

/**
 * ParseID 解析雪花ID
 * @param id 雪花ID
 * @return timestamp 时间戳
 * @return datacenterID 数据中心ID
 * @return workerID 机器ID
 * @return sequence 序列号
 */
func ParseID(id int64) (timestamp, datacenterID, workerID, sequence int64) {
	timestamp = (id >> timestampShift) + epoch
	datacenterID = (id >> datacenterShift) & maxDatacenterID
	workerID = (id >> workerShift) & maxWorkerID
	sequence = id & maxSequence
	return
}

/**
 * GenerateID 生成GUID(兼容原有代码)
 * @return string GUID字符串
 * @description 20251001 陈凤庆 改为生成GUID，解决JavaScript精度丢失问题
 */
func GenerateID() string {
	return GenerateGUID()
}

/**
 * GenerateIDString 生成GUID字符串(兼容原有代码)
 * @return string GUID字符串
 * @description 20251001 陈凤庆 改为生成GUID，解决JavaScript精度丢失问题
 */
func GenerateIDString() string {
	return GenerateGUID()
}
