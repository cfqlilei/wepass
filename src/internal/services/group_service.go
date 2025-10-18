package services

import (
	"fmt"
	"time"

	"wepassword/internal/database"
	"wepassword/internal/models"
	"wepassword/internal/utils"
)

/**
 * 分组服务
 * @author 陈凤庆
 * @description 管理密码分组的增删改查操作
 */

/**
 * GroupService 分组服务
 */
type GroupService struct {
	dbManager *database.DatabaseManager
}

/**
 * NewGroupService 创建新的分组服务
 * @param dbManager 数据库管理器
 * @return *GroupService 分组服务实例
 */
func NewGroupService(dbManager *database.DatabaseManager) *GroupService {
	return &GroupService{
		dbManager: dbManager,
	}
}

/**
 * GetAllGroups 获取所有分组
 * @return []models.Group 分组列表
 * @return error 错误信息
 */
func (gs *GroupService) GetAllGroups() ([]models.Group, error) {
	if !gs.dbManager.IsOpened() {
		return nil, fmt.Errorf("数据库未打开")
	}

	db := gs.dbManager.GetDB()
	// 20251001 陈凤庆 删除created_by和updated_by字段
	// 20251002 陈凤庆 删除parent_id字段，修复ORDER BY语句
	rows, err := db.Query(`
		SELECT id, name, icon, sort_order, created_at, updated_at
		FROM groups
		ORDER BY sort_order, name
	`)
	if err != nil {
		return nil, fmt.Errorf("查询分组失败: %w", err)
	}
	defer rows.Close()

	// 20251002 陈凤庆 初始化为非nil的空切片，确保前端接收到[]而不是null
	groups := make([]models.Group, 0)
	for rows.Next() {
		var group models.Group
		err := rows.Scan(
			&group.ID, &group.Name, &group.Icon, &group.SortOrder,
			&group.CreatedAt, &group.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描分组数据失败: %w", err)
		}
		groups = append(groups, group)
	}

	return groups, nil
}

/**
 * GetGroupByID 根据ID获取分组
 * @param id 分组ID
 * @return *models.Group 分组信息
 * @return error 错误信息
 * @modify 20251001 陈凤庆 ID参数改为string类型，使用GUID
 */
func (gs *GroupService) GetGroupByID(id string) (*models.Group, error) {
	if !gs.dbManager.IsOpened() {
		return nil, fmt.Errorf("数据库未打开")
	}

	db := gs.dbManager.GetDB()
	group := &models.Group{}
	// 20251001 陈凤庆 删除created_by和updated_by字段
	// 20251002 陈凤庆 删除parent_id字段
	err := db.QueryRow(`
		SELECT id, name, icon, sort_order, created_at, updated_at
		FROM groups
		WHERE id = ?
	`, id).Scan(
		&group.ID, &group.Name, &group.Icon, &group.SortOrder,
		&group.CreatedAt, &group.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("查询分组失败: %w", err)
	}

	return group, nil
}

/**
 * CreateGroup 创建新分组
 * @param name 分组名称
 * @return models.Group 创建的分组
 * @return error 错误信息
 * @modify 20251002 陈凤庆 删除parentID参数，不需要层级结构
 */
func (gs *GroupService) CreateGroup(name string) (models.Group, error) {
	if !gs.dbManager.IsOpened() {
		return models.Group{}, fmt.Errorf("数据库未打开")
	}

	if name == "" {
		return models.Group{}, fmt.Errorf("分组名称不能为空")
	}

	db := gs.dbManager.GetDB()
	now := time.Now()

	// 获取最大排序号
	var maxSortOrder int
	err := db.QueryRow(`
		SELECT COALESCE(MAX(sort_order), 0)
		FROM groups
	`).Scan(&maxSortOrder)
	if err != nil {
		return models.Group{}, fmt.Errorf("获取排序号失败: %w", err)
	}

	// 生成新的GUID
	newID := utils.GenerateGUID()

	// 插入新分组
	_, err = db.Exec(`
		INSERT INTO groups (id, name, icon, sort_order, created_at, updated_at)
		VALUES (?, ?, 'fa-folder', ?, ?, ?)
	`, newID, name, maxSortOrder+1, now, now)
	if err != nil {
		return models.Group{}, fmt.Errorf("创建分组失败: %w", err)
	}

	// 20251003 陈凤庆 取消自动创建默认类型，用户可以手动创建需要的类型
	// if err := gs.dbManager.CreateDefaultTypeForGroup(newID); err != nil {
	// 	return models.Group{}, fmt.Errorf("创建默认类型失败: %w", err)
	// }

	// 返回创建的分组
	group := models.Group{
		ID:        newID,
		Name:      name,
		Icon:      "fa-folder",
		SortOrder: maxSortOrder + 1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return group, nil
}

/**
 * UpdateGroup 更新分组
 * @param group 分组信息
 * @return error 错误信息
 */
func (gs *GroupService) UpdateGroup(group models.Group) error {
	if !gs.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	db := gs.dbManager.GetDB()
	now := time.Now()

	// 20251002 陈凤庆 删除parent_id和updated_by字段
	_, err := db.Exec(`
		UPDATE groups
		SET name = ?, icon = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`, group.Name, group.Icon, group.SortOrder, now, group.ID)

	if err != nil {
		return fmt.Errorf("更新分组失败: %w", err)
	}

	return nil
}

/**
 * DeleteGroup 删除分组
 * @param id 分组ID
 * @return error 错误信息
 * @modify 20251001 陈凤庆 ID参数改为string类型，使用GUID
 * @modify 20251002 陈凤庆 增加删除前验证：检查分组下是否有账号和类别
 */
func (gs *GroupService) DeleteGroup(id string) error {
	if !gs.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	// 检查是否是默认分组（可以通过名称或特殊标识来判断）
	group, err := gs.GetGroupByID(id)
	if err != nil {
		return err
	}

	// 20251002 陈凤庆 删除parent_id检查，只检查名称
	if group.Name == "默认" {
		return fmt.Errorf("默认分组不能删除")
	}

	db := gs.dbManager.GetDB()

	// 20251002 陈凤庆 删除前验证：检查分组下是否有类别
	var typeCount int
	err = db.QueryRow("SELECT COUNT(*) FROM types WHERE group_id = ?", id).Scan(&typeCount)
	if err != nil {
		return fmt.Errorf("检查分组下的类别失败: %w", err)
	}

	if typeCount > 0 {
		return fmt.Errorf("该分组下还有 %d 个类别，请先删除相关类别", typeCount)
	}

	// 20251002 陈凤庆 删除前验证：检查分组下是否有账号
	var accountCount int
	err = db.QueryRow("SELECT COUNT(*) FROM accounts WHERE typeid IN (SELECT id FROM types WHERE group_id = ?)", id).Scan(&accountCount)
	if err != nil {
		return fmt.Errorf("检查分组下的账号失败: %w", err)
	}

	if accountCount > 0 {
		return fmt.Errorf("该分组下还有 %d 个账号，请先删除相关账号", accountCount)
	}

	// 删除分组（此时分组下已经没有类别和账号）
	// 20251002 陈凤庆 不使用:=赋值，先声明变量类型
	var deleteErr error
	_, deleteErr = db.Exec("DELETE FROM groups WHERE id = ?", id)
	if deleteErr != nil {
		return fmt.Errorf("删除分组失败: %w", deleteErr)
	}

	return nil
}

/**
 * SearchGroups 搜索分组
 * @param keyword 搜索关键词
 * @return []models.Group 搜索结果
 * @return error 错误信息
 */
func (gs *GroupService) SearchGroups(keyword string) ([]models.Group, error) {
	if !gs.dbManager.IsOpened() {
		return nil, fmt.Errorf("数据库未打开")
	}

	db := gs.dbManager.GetDB()
	// 20251002 陈凤庆 删除parent_id字段
	rows, err := db.Query(`
		SELECT id, name, icon, sort_order, created_at, updated_at
		FROM groups
		WHERE name LIKE ?
		ORDER BY sort_order, name
	`, "%"+keyword+"%")
	if err != nil {
		return nil, fmt.Errorf("搜索分组失败: %w", err)
	}
	defer rows.Close()

	// 20251002 陈凤庆 初始化为非nil的空切片，确保前端接收到[]而不是null
	groups := make([]models.Group, 0)
	for rows.Next() {
		var group models.Group
		err := rows.Scan(
			&group.ID, &group.Name, &group.Icon, &group.SortOrder,
			&group.CreatedAt, &group.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描分组数据失败: %w", err)
		}
		groups = append(groups, group)
	}

	return groups, nil
}

/**
 * RenameGroup 重命名分组
 * @param id 分组ID
 * @param newName 新的分组名称
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增分组重命名功能
 */
func (gs *GroupService) RenameGroup(id string, newName string) error {
	if !gs.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	if id == "" {
		return fmt.Errorf("分组ID不能为空")
	}

	if newName == "" {
		return fmt.Errorf("分组名称不能为空")
	}

	// 检查分组是否存在
	group, err := gs.GetGroupByID(id)
	if err != nil {
		return fmt.Errorf("分组不存在: %w", err)
	}

	// 检查是否是默认分组
	if group.Name == "默认" {
		return fmt.Errorf("默认分组不能重命名")
	}

	db := gs.dbManager.GetDB()
	now := time.Now()

	// 更新分组名称
	_, err = db.Exec(`
		UPDATE groups
		SET name = ?, updated_at = ?
		WHERE id = ?
	`, newName, now, id)

	if err != nil {
		return fmt.Errorf("重命名分组失败: %w", err)
	}

	return nil
}

/**
 * MoveGroupLeft 将分组向左移动一位
 * @param id 分组ID
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增分组左移功能
 */
func (gs *GroupService) MoveGroupLeft(id string) error {
	if !gs.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	if id == "" {
		return fmt.Errorf("分组ID不能为空")
	}

	db := gs.dbManager.GetDB()

	// 获取当前分组的排序号
	var currentSortOrder int
	err := db.QueryRow(`
		SELECT sort_order
		FROM groups
		WHERE id = ?
	`, id).Scan(&currentSortOrder)
	if err != nil {
		return fmt.Errorf("获取分组排序号失败: %w", err)
	}

	// 查找左边的分组（排序号小于当前分组的最大排序号）
	var leftGroupID string
	var leftSortOrder int
	err = db.QueryRow(`
		SELECT id, sort_order
		FROM groups
		WHERE sort_order < ?
		ORDER BY sort_order DESC
		LIMIT 1
	`, currentSortOrder).Scan(&leftGroupID, &leftSortOrder)
	if err != nil {
		// 没有找到左边的分组，说明已经是最左边了
		return fmt.Errorf("分组已经在最左边，无法继续左移")
	}

	// 交换两个分组的排序号
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()

	// 更新当前分组的排序号
	_, err = tx.Exec(`
		UPDATE groups
		SET sort_order = ?, updated_at = ?
		WHERE id = ?
	`, leftSortOrder, now, id)
	if err != nil {
		return fmt.Errorf("更新当前分组排序号失败: %w", err)
	}

	// 更新左边分组的排序号
	_, err = tx.Exec(`
		UPDATE groups
		SET sort_order = ?, updated_at = ?
		WHERE id = ?
	`, currentSortOrder, now, leftGroupID)
	if err != nil {
		return fmt.Errorf("更新左边分组排序号失败: %w", err)
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

/**
 * MoveGroupRight 将分组向右移动一位
 * @param id 分组ID
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增分组右移功能
 */
func (gs *GroupService) MoveGroupRight(id string) error {
	if !gs.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	if id == "" {
		return fmt.Errorf("分组ID不能为空")
	}

	db := gs.dbManager.GetDB()

	// 获取当前分组的排序号
	var currentSortOrder int
	err := db.QueryRow(`
		SELECT sort_order
		FROM groups
		WHERE id = ?
	`, id).Scan(&currentSortOrder)
	if err != nil {
		return fmt.Errorf("获取分组排序号失败: %w", err)
	}

	// 查找右边的分组（排序号大于当前分组的最小排序号）
	var rightGroupID string
	var rightSortOrder int
	err = db.QueryRow(`
		SELECT id, sort_order
		FROM groups
		WHERE sort_order > ?
		ORDER BY sort_order ASC
		LIMIT 1
	`, currentSortOrder).Scan(&rightGroupID, &rightSortOrder)
	if err != nil {
		// 没有找到右边的分组，说明已经是最右边了
		return fmt.Errorf("分组已经在最右边，无法继续右移")
	}

	// 交换两个分组的排序号
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()

	// 更新当前分组的排序号
	_, err = tx.Exec(`
		UPDATE groups
		SET sort_order = ?, updated_at = ?
		WHERE id = ?
	`, rightSortOrder, now, id)
	if err != nil {
		return fmt.Errorf("更新当前分组排序号失败: %w", err)
	}

	// 更新右边分组的排序号
	_, err = tx.Exec(`
		UPDATE groups
		SET sort_order = ?, updated_at = ?
		WHERE id = ?
	`, currentSortOrder, now, rightGroupID)
	if err != nil {
		return fmt.Errorf("更新右边分组排序号失败: %w", err)
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

/**
 * UpdateGroupSortOrder 更新分组排序
 * @param groupID 分组ID
 * @param newSortOrder 新的排序号
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增分组排序更新方法，用于拖拽排序
 */
func (gs *GroupService) UpdateGroupSortOrder(groupID string, newSortOrder int) error {
	if !gs.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	if groupID == "" {
		return fmt.Errorf("分组ID不能为空")
	}

	db := gs.dbManager.GetDB()
	now := time.Now()

	// 20251002 陈凤庆 不使用:=赋值，先声明变量类型
	var updateErr error
	_, updateErr = db.Exec(`
		UPDATE groups
		SET sort_order = ?, updated_at = ?
		WHERE id = ?
	`, newSortOrder, now, groupID)

	if updateErr != nil {
		return fmt.Errorf("更新分组排序失败: %w", updateErr)
	}

	return nil
}
