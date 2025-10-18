<template>
  <div class="groups-bar">
    <!-- 正常分组 -->
    <div class="normal-groups">
      <GroupTag
        v-for="group in normalGroups"
        :key="group.id"
        :group="group"
        :active="isGroupActive(group)"
        @click="handleSelectGroup"
        @contextmenu="handleGroupContextMenu"
        @drag-reorder="handleGroupDragReorder"
      />
      <div class="group-tag add" @click="handleCreateGroup">
        <el-icon><Plus /></el-icon>
      </div>
    </div>

    <!-- 搜索结果分组（右侧隔离显示） -->
    <div v-if="searchGroups.length > 0" class="search-groups">
      <!-- <div class="search-separator"></div> -->
      <GroupTag
        v-for="group in searchGroups"
        :key="group.id"
        :group="group"
        :active="isGroupActive(group)"
        :class="{ 'search-result': group.isSearchResult }"
        @click="handleSelectGroup"
      />
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import GroupTag from './GroupTag.vue'

/**
 * 分组栏组件
 * @author 陈凤庆
 * @date 20251001
 * @description 显示分组标签，处理分组选择和创建
 */
export default {
  name: 'GroupsBar',
  components: {
    GroupTag
  },
  props: {
    groups: {
      type: Array,
      default: () => []
    },
    currentGroupId: {
      type: Number,
      default: null
    }
  },
  emits: ['select-group', 'create-group', 'show-group-menu', 'drag-reorder'],
  setup(props, { emit }) {
    /**
     * 分离正常分组和搜索分组
     * @author 陈凤庆
     * @date 2025-10-03
     * @description 将搜索结果分组与正常分组分开显示
     */
    const normalGroups = computed(() => {
      return props.groups.filter(group => !group.isSearchResult)
    })

    const searchGroups = computed(() => {
      return props.groups.filter(group => group.isSearchResult)
    })

    /**
     * 判断分组是否处于选中状态
     * @modify 20251002 陈凤庆 使用字符串比较，避免大整数精度问题
     */
    const isGroupActive = (group) => {
      // 直接使用字符串比较，后端已返回字符串ID
      const currentIdStr = String(props.currentGroupId || '')
      const groupIdStr = String(group.id || '')
      const active = currentIdStr === groupIdStr

      console.log(`[GroupsBar] 检查分组 "${group.name}" (ID: ${group.id}) 是否选中:`, {
        currentGroupId: props.currentGroupId,
        groupId: group.id,
        active: active,
        currentIdStr: currentIdStr,
        groupIdStr: groupIdStr,
        stringComparison: `"${currentIdStr}" === "${groupIdStr}"`
      })
      return active
    }

    /**
     * 处理选择分组
     */
    const handleSelectGroup = (groupId) => {
      console.log(`[GroupsBar] 选择分组: ${groupId}`)
      emit('select-group', groupId)
    }

    /**
     * 处理创建分组
     */
    const handleCreateGroup = () => {
      emit('create-group')
    }

    /**
     * 处理分组右键菜单
     */
    const handleGroupContextMenu = (event, group) => {
      console.log('[GroupsBar] 分组右键菜单事件:', group.name, event)
      emit('show-group-menu', event, group)
    }

    /**
     * 处理分组拖拽排序
     * @param {string} sourceGroupId 源分组ID
     * @param {string} targetGroupId 目标分组ID
     * @author 20251002 陈凤庆 处理分组拖拽排序
     */
    const handleGroupDragReorder = (sourceGroupId, targetGroupId) => {
      console.log('[GroupsBar] 分组拖拽排序:', sourceGroupId, '->', targetGroupId)
      emit('drag-reorder', sourceGroupId, targetGroupId)
    }

    return {
      Plus,
      normalGroups,
      searchGroups,
      isGroupActive,
      handleSelectGroup,
      handleCreateGroup,
      handleGroupContextMenu,
      handleGroupDragReorder
    }
  }
}
</script>

<style scoped>
.groups-bar {
  height: 36px;
  border-bottom: 1px solid #dee2e6;
  padding: 0 4px;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  overflow-x: auto;
  background-color: #f8f9fa;
  /* 20251003 陈凤庆 隐藏左侧滚动条 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE and Edge */
}

/* 20251003 陈凤庆 WebKit浏览器隐藏滚动条 */
.groups-bar::-webkit-scrollbar {
  display: none;
}

.normal-groups {
  display: flex;
  align-items: flex-end;
  gap: 2px;
  flex: 1;
}

.search-groups {
  display: flex;
  align-items: flex-end;
  gap: 2px;
  margin-left: auto;
}

.search-separator {
  width: 1px;
  height: 20px;
  background-color: #dee2e6;
  margin: 0 8px;
  align-self: center;
}

.group-tag.add {
  padding: 6px 10px;
  font-size: 16px;
  background-color: #e9ecef;
  border: 1px solid #dee2e6;
  border-bottom-color: transparent;
  border-radius: 4px 4px 0 0;
  cursor: pointer;
  white-space: nowrap;
  height: 35px;
  display: flex;
  align-items: center;
  gap: 6px;
  color: #495057;
}

.group-tag.add:hover {
  background-color: #dee2e6;
  color: #212529;
}

/* 搜索结果分组特殊样式 */
.group-tag.search-result {
  background-color: #e3f2fd;
  border-color: #90caf9;
  color: #1976d2;
}

.group-tag.search-result.active {
  background-color: #ffffff;
  color: #1976d2;
  border-color: #90caf9 #90caf9 #ffffff;
  border-top: 2px solid #1976d2;
}

.group-tag.search-result:hover:not(.active) {
  background-color: #bbdefb;
  color: #0d47a1;
}
</style>

