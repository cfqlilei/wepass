package services

import (
	"fmt"
	"log"
	"time"

	"wepassword/internal/database"
	"wepassword/internal/models"
	"wepassword/internal/utils"
)

/**
 * 类型服务
 * @author 陈凤庆
 * @description 管理类型的增删改查操作
 * @modify 20251002 陈凤庆 tab_service.go改名为type_service.go，对应types表
 */

/**
 * TypeService 类型服务
 */
type TypeService struct {
	dbManager *database.DatabaseManager
}

/**
 * NewTypeService 创建新的类型服务
 * @param dbManager 数据库管理器
 * @return *TypeService 类型服务实例
 */
func NewTypeService(dbManager *database.DatabaseManager) *TypeService {
	return &TypeService{
		dbManager: dbManager,
	}
}

/**
 * GetTypesByGroup 根据分组ID获取类型列表
 * @param groupID 分组ID
 * @return []models.Type 类型列表
 * @return error 错误信息
 * @modify 20251002 陈凤庆 GetTabsByGroup改名为GetTypesByGroup，增强日志记录
 */
func (ts *TypeService) GetTypesByGroup(groupID string) ([]models.Type, error) {
	// 20251002 陈凤庆 添加详细日志记录
	log.Printf("[TypeService] GetTypesByGroup 开始执行，分组ID: %s", groupID)

	// 20251002 陈凤庆 参数验证
	if groupID == "" {
		log.Printf("[TypeService] GetTypesByGroup 失败：分组ID为空")
		return nil, fmt.Errorf("分组ID不能为空")
	}

	if !ts.dbManager.IsOpened() {
		log.Printf("[TypeService] GetTypesByGroup 失败：数据库未打开")
		return nil, fmt.Errorf("数据库未打开")
	}

	db := ts.dbManager.GetDB()
	log.Printf("[TypeService] 开始查询types表，分组ID: %s", groupID)

	// 20251002 陈凤庆 查询types表
	rows, err := db.Query(`
		SELECT id, name, icon, filter, group_id, sort_order, created_at, updated_at
		FROM types
		WHERE group_id = ?
		ORDER BY sort_order, name
	`, groupID)
	if err != nil {
		log.Printf("[TypeService] 查询types表失败，分组ID: %s, 错误: %v", groupID, err)
		return nil, fmt.Errorf("查询类型失败: %w", err)
	}
	defer rows.Close()

	// 20251002 陈凤庆 初始化为非nil的空切片，确保前端接收到[]而不是null
	types := make([]models.Type, 0)
	var count int
	for rows.Next() {
		var typeItem models.Type
		err := rows.Scan(
			&typeItem.ID, &typeItem.Name, &typeItem.Icon, &typeItem.Filter, &typeItem.GroupID, &typeItem.SortOrder,
			&typeItem.CreatedAt, &typeItem.UpdatedAt,
		)
		if err != nil {
			log.Printf("[TypeService] 扫描类型数据失败，分组ID: %s, 错误: %v", groupID, err)
			return nil, fmt.Errorf("扫描类型数据失败: %w", err)
		}
		types = append(types, typeItem)
		count++
		log.Printf("[TypeService] 扫描到类型 %d: ID=%s, Name='%s', Icon='%s'", count, typeItem.ID, typeItem.Name, typeItem.Icon)
	}

	log.Printf("[TypeService] GetTypesByGroup 完成，分组ID: %s, 返回 %d 个类型", groupID, len(types))
	return types, nil
}

/**
 * GetAllTypes 获取所有类型
 * @return []models.Type 所有类型列表
 * @return error 错误信息
 * @author 20251004 陈凤庆 新增获取所有类型的方法，用于导出功能
 */
func (ts *TypeService) GetAllTypes() ([]models.Type, error) {
	log.Printf("[TypeService] GetAllTypes 开始执行")

	if !ts.dbManager.IsOpened() {
		log.Printf("[TypeService] GetAllTypes 失败：数据库未打开")
		return nil, fmt.Errorf("数据库未打开")
	}

	db := ts.dbManager.GetDB()
	log.Printf("[TypeService] 开始查询所有types")

	rows, err := db.Query(`
		SELECT id, name, icon, filter, group_id, sort_order, created_at, updated_at
		FROM types
		ORDER BY group_id, sort_order, name
	`)
	if err != nil {
		log.Printf("[TypeService] 查询所有types失败，错误: %v", err)
		return nil, fmt.Errorf("查询所有类型失败: %w", err)
	}
	defer rows.Close()

	types := make([]models.Type, 0)
	var count int
	for rows.Next() {
		var typeItem models.Type
		err := rows.Scan(
			&typeItem.ID, &typeItem.Name, &typeItem.Icon, &typeItem.Filter, &typeItem.GroupID, &typeItem.SortOrder,
			&typeItem.CreatedAt, &typeItem.UpdatedAt,
		)
		if err != nil {
			log.Printf("[TypeService] 扫描类型数据失败，错误: %v", err)
			return nil, fmt.Errorf("扫描类型数据失败: %w", err)
		}
		types = append(types, typeItem)
		count++
		log.Printf("[TypeService] 扫描到类型 %d: ID=%s, Name='%s', GroupID='%s'", count, typeItem.ID, typeItem.Name, typeItem.GroupID)
	}

	log.Printf("[TypeService] GetAllTypes 完成，返回 %d 个类型", len(types))
	return types, nil
}

/**
 * GetTypeByID 根据ID获取类型
 * @param id 类型ID
 * @return *models.Type 类型信息
 * @return error 错误信息
 * @modify 20251002 陈凤庆 GetTabByID改名为GetTypeByID
 */
func (ts *TypeService) GetTypeByID(id string) (*models.Type, error) {
	if !ts.dbManager.IsOpened() {
		return nil, fmt.Errorf("数据库未打开")
	}

	db := ts.dbManager.GetDB()
	var typeItem models.Type
	// 20251002 陈凤庆 查询types表
	err := db.QueryRow(`
		SELECT id, name, icon, filter, group_id, sort_order, created_at, updated_at
		FROM types
		WHERE id = ?
	`, id).Scan(
		&typeItem.ID, &typeItem.Name, &typeItem.Icon, &typeItem.Filter, &typeItem.GroupID, &typeItem.SortOrder,
		&typeItem.CreatedAt, &typeItem.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("查询类型失败: %w", err)
	}

	return &typeItem, nil
}

/**
 * CreateType 创建类型
 * @param name 类型名称
 * @param groupID 分组ID
 * @param icon 图标
 * @return models.Type 创建的类型
 * @return error 错误信息
 * @modify 20251002 陈凤庆 CreateTab改名为CreateType
 */
func (ts *TypeService) CreateType(name string, groupID string, icon string) (models.Type, error) {
	if !ts.dbManager.IsOpened() {
		return models.Type{}, fmt.Errorf("数据库未打开")
	}

	if name == "" {
		return models.Type{}, fmt.Errorf("类型名称不能为空")
	}

	if groupID == "" {
		return models.Type{}, fmt.Errorf("分组ID不能为空")
	}

	db := ts.dbManager.GetDB()
	now := time.Now()

	// 获取最大排序号
	var maxSortOrder int
	err := db.QueryRow(`
		SELECT COALESCE(MAX(sort_order), 0)
		FROM types
		WHERE group_id = ?
	`, groupID).Scan(&maxSortOrder)
	if err != nil {
		return models.Type{}, fmt.Errorf("获取排序号失败: %w", err)
	}

	// 生成新的GUID
	newID := utils.GenerateGUID()

	// 插入新类型
	// 20251002 陈凤庆 不使用:=赋值，先声明变量类型
	var insertErr error
	_, insertErr = db.Exec(`
		INSERT INTO types (id, name, icon, group_id, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, newID, name, icon, groupID, maxSortOrder+1, now, now)
	if insertErr != nil {
		return models.Type{}, fmt.Errorf("创建类型失败: %w", insertErr)
	}

	// 返回创建的类型
	typeItem := models.Type{
		ID:        newID,
		Name:      name,
		Icon:      icon,
		GroupID:   groupID,
		SortOrder: maxSortOrder + 1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return typeItem, nil
}

/**
 * UpdateType 更新类型
 * @param typeItem 类型信息
 * @return error 错误信息
 * @modify 20251002 陈凤庆 UpdateTab改名为UpdateType
 */
func (ts *TypeService) UpdateType(typeItem models.Type) error {
	if !ts.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	db := ts.dbManager.GetDB()
	now := time.Now()

	// 20251002 陈凤庆 不使用:=赋值，先声明变量类型
	var updateErr error
	_, updateErr = db.Exec(`
		UPDATE types 
		SET name = ?, icon = ?, filter = ?, group_id = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`, typeItem.Name, typeItem.Icon, typeItem.Filter, typeItem.GroupID, typeItem.SortOrder, now, typeItem.ID)

	if updateErr != nil {
		return fmt.Errorf("更新类型失败: %w", updateErr)
	}

	return nil
}

/**
 * DeleteType 删除类型
 * @param id 类型ID
 * @return error 错误信息
 * @modify 20251002 陈凤庆 DeleteTab改名为DeleteType，增加账号检查
 */
func (ts *TypeService) DeleteType(id string) error {
	if !ts.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	db := ts.dbManager.GetDB()

	// 20251002 陈凤庆 删除前检查accounts表中是否有对应typeid
	var accountCount int
	err := db.QueryRow("SELECT COUNT(*) FROM accounts WHERE typeid = ?", id).Scan(&accountCount)
	if err != nil {
		return fmt.Errorf("检查账号失败: %w", err)
	}

	if accountCount > 0 {
		return fmt.Errorf("该标签下还有 %d 个账号，请先删除相关账号", accountCount)
	}

	// 删除类型
	// 20251002 陈凤庆 不使用:=赋值，先声明变量类型
	var deleteErr error
	_, deleteErr = db.Exec("DELETE FROM types WHERE id = ?", id)
	if deleteErr != nil {
		return fmt.Errorf("删除类型失败: %w", deleteErr)
	}

	return nil
}

/**
 * SearchTypes 搜索类型
 * @param keyword 搜索关键词
 * @return []models.Type 搜索结果
 * @return error 错误信息
 * @modify 20251002 陈凤庆 SearchTabs改名为SearchTypes
 */
func (ts *TypeService) SearchTypes(keyword string) ([]models.Type, error) {
	if !ts.dbManager.IsOpened() {
		return nil, fmt.Errorf("数据库未打开")
	}

	db := ts.dbManager.GetDB()
	// 20251002 陈凤庆 查询types表
	rows, err := db.Query(`
		SELECT id, name, icon, filter, group_id, sort_order, created_at, updated_at
		FROM types
		WHERE name LIKE ?
		ORDER BY sort_order, name
	`, "%"+keyword+"%")
	if err != nil {
		return nil, fmt.Errorf("搜索类型失败: %w", err)
	}
	defer rows.Close()

	// 20251002 陈凤庆 初始化为非nil的空切片，确保前端接收到[]而不是null
	types := make([]models.Type, 0)
	for rows.Next() {
		var typeItem models.Type
		err := rows.Scan(
			&typeItem.ID, &typeItem.Name, &typeItem.Icon, &typeItem.Filter, &typeItem.GroupID, &typeItem.SortOrder,
			&typeItem.CreatedAt, &typeItem.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描类型数据失败: %w", err)
		}
		types = append(types, typeItem)
	}

	return types, nil
}

/**
 * MoveTypeUp 将类型向上移动一位
 * @param id 类型ID
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增类型上移功能
 */
func (ts *TypeService) MoveTypeUp(id string) error {
	if !ts.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	if id == "" {
		return fmt.Errorf("类型ID不能为空")
	}

	db := ts.dbManager.GetDB()

	// 获取当前类型的排序号和分组ID
	var currentSortOrder int
	var groupID string
	err := db.QueryRow(`
		SELECT sort_order, group_id
		FROM types
		WHERE id = ?
	`, id).Scan(&currentSortOrder, &groupID)
	if err != nil {
		return fmt.Errorf("获取类型排序号失败: %w", err)
	}

	// 查找上面的类型（同分组内排序号小于当前类型的最大排序号）
	var upperTypeID string
	var upperSortOrder int
	err = db.QueryRow(`
		SELECT id, sort_order
		FROM types
		WHERE group_id = ? AND sort_order < ?
		ORDER BY sort_order DESC
		LIMIT 1
	`, groupID, currentSortOrder).Scan(&upperTypeID, &upperSortOrder)
	if err != nil {
		// 没有找到上面的类型，说明已经是最上面了
		return fmt.Errorf("类型已经在最上面，无法继续上移")
	}

	// 交换两个类型的排序号
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()

	// 更新当前类型的排序号
	_, err = tx.Exec(`
		UPDATE types
		SET sort_order = ?, updated_at = ?
		WHERE id = ?
	`, upperSortOrder, now, id)
	if err != nil {
		return fmt.Errorf("更新当前类型排序号失败: %w", err)
	}

	// 更新上面类型的排序号
	_, err = tx.Exec(`
		UPDATE types
		SET sort_order = ?, updated_at = ?
		WHERE id = ?
	`, currentSortOrder, now, upperTypeID)
	if err != nil {
		return fmt.Errorf("更新上面类型排序号失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

/**
 * MoveTypeDown 将类型向下移动一位
 * @param id 类型ID
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增类型下移功能
 */
func (ts *TypeService) MoveTypeDown(id string) error {
	if !ts.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	if id == "" {
		return fmt.Errorf("类型ID不能为空")
	}

	db := ts.dbManager.GetDB()

	// 获取当前类型的排序号和分组ID
	var currentSortOrder int
	var groupID string
	err := db.QueryRow(`
		SELECT sort_order, group_id
		FROM types
		WHERE id = ?
	`, id).Scan(&currentSortOrder, &groupID)
	if err != nil {
		return fmt.Errorf("获取类型排序号失败: %w", err)
	}

	// 查找下面的类型（同分组内排序号大于当前类型的最小排序号）
	var lowerTypeID string
	var lowerSortOrder int
	err = db.QueryRow(`
		SELECT id, sort_order
		FROM types
		WHERE group_id = ? AND sort_order > ?
		ORDER BY sort_order ASC
		LIMIT 1
	`, groupID, currentSortOrder).Scan(&lowerTypeID, &lowerSortOrder)
	if err != nil {
		// 没有找到下面的类型，说明已经是最下面了
		return fmt.Errorf("类型已经在最下面，无法继续下移")
	}

	// 交换两个类型的排序号
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()

	// 更新当前类型的排序号
	_, err = tx.Exec(`
		UPDATE types
		SET sort_order = ?, updated_at = ?
		WHERE id = ?
	`, lowerSortOrder, now, id)
	if err != nil {
		return fmt.Errorf("更新当前类型排序号失败: %w", err)
	}

	// 更新下面类型的排序号
	_, err = tx.Exec(`
		UPDATE types
		SET sort_order = ?, updated_at = ?
		WHERE id = ?
	`, currentSortOrder, now, lowerTypeID)
	if err != nil {
		return fmt.Errorf("更新下面类型排序号失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

/**
 * UpdateTypeSortOrder 更新类型排序
 * @param typeID 类型ID
 * @param newSortOrder 新的排序号
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增类型排序更新功能，用于拖拽排序
 */
func (ts *TypeService) UpdateTypeSortOrder(typeID string, newSortOrder int) error {
	if !ts.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	if typeID == "" {
		return fmt.Errorf("类型ID不能为空")
	}

	db := ts.dbManager.GetDB()
	now := time.Now()

	// 更新类型的排序号
	// 20251002 陈凤庆 不使用:=赋值，先声明变量类型
	var updateErr error
	_, updateErr = db.Exec(`
		UPDATE types
		SET sort_order = ?, updated_at = ?
		WHERE id = ?
	`, newSortOrder, now, typeID)

	if updateErr != nil {
		return fmt.Errorf("更新类型排序失败: %w", updateErr)
	}

	return nil
}

/**
 * InsertTypeAfter 在指定标签后插入新标签
 * @param name 新标签名称
 * @param groupID 分组ID
 * @param icon 图标
 * @param afterTypeID 在此标签后插入，如果为空则插入到最后
 * @return models.Type 创建的类型
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增在指定标签后插入新标签功能
 */
func (ts *TypeService) InsertTypeAfter(name string, groupID string, icon string, afterTypeID string) (models.Type, error) {
	if !ts.dbManager.IsOpened() {
		return models.Type{}, fmt.Errorf("数据库未打开")
	}

	if name == "" {
		return models.Type{}, fmt.Errorf("类型名称不能为空")
	}

	if groupID == "" {
		return models.Type{}, fmt.Errorf("分组ID不能为空")
	}

	db := ts.dbManager.GetDB()
	now := time.Now()

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return models.Type{}, fmt.Errorf("开始事务失败: %w", err)
	}
	defer tx.Rollback()

	var insertSortOrder int

	if afterTypeID == "" {
		// 插入到最后，获取最大排序号
		err = tx.QueryRow(`
			SELECT COALESCE(MAX(sort_order), 0)
			FROM types
			WHERE group_id = ?
		`, groupID).Scan(&insertSortOrder)
		if err != nil {
			return models.Type{}, fmt.Errorf("获取最大排序号失败: %w", err)
		}
		insertSortOrder++
	} else {
		// 获取指定标签的排序号
		var afterSortOrder int
		err = tx.QueryRow(`
			SELECT sort_order
			FROM types
			WHERE id = ? AND group_id = ?
		`, afterTypeID, groupID).Scan(&afterSortOrder)
		if err != nil {
			return models.Type{}, fmt.Errorf("获取指定标签排序号失败: %w", err)
		}

		// 新标签的排序号为指定标签的排序号+1
		insertSortOrder = afterSortOrder + 1

		// 将所有排序号大于等于insertSortOrder的标签排序号+1
		_, err = tx.Exec(`
			UPDATE types
			SET sort_order = sort_order + 1, updated_at = ?
			WHERE group_id = ? AND sort_order >= ?
		`, now, groupID, insertSortOrder)
		if err != nil {
			return models.Type{}, fmt.Errorf("调整其他标签排序号失败: %w", err)
		}
	}

	// 生成新的GUID
	newID := utils.GenerateGUID()

	// 插入新类型
	// 20251002 陈凤庆 不使用:=赋值，先声明变量类型
	var insertErr error
	_, insertErr = tx.Exec(`
		INSERT INTO types (id, name, icon, group_id, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, newID, name, icon, groupID, insertSortOrder, now, now)
	if insertErr != nil {
		return models.Type{}, fmt.Errorf("创建类型失败: %w", insertErr)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return models.Type{}, fmt.Errorf("提交事务失败: %w", err)
	}

	// 返回创建的类型
	typeItem := models.Type{
		ID:        newID,
		Name:      name,
		Icon:      icon,
		GroupID:   groupID,
		SortOrder: insertSortOrder,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return typeItem, nil
}

/**
 * CreateDefaultTypeForGroup 为分组创建默认类型
 * @param groupID 分组ID
 * @return error 错误信息
 * @modify 20251002 陈凤庆 CreateDefaultTabForGroup改名为CreateDefaultTypeForGroup
 */
func (ts *TypeService) CreateDefaultTypeForGroup(groupID string) error {
	if !ts.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	// 检查是否已经有类型
	types, err := ts.GetTypesByGroup(groupID)
	if err != nil {
		return fmt.Errorf("检查类型失败: %w", err)
	}

	if len(types) > 0 {
		return nil // 已经有类型，不需要创建默认类型
	}

	// 创建默认类型
	_, err = ts.CreateType("默认", groupID, "fa-folder")
	if err != nil {
		return fmt.Errorf("创建默认类型失败: %w", err)
	}

	return nil
}
