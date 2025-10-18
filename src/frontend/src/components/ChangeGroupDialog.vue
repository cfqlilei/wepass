<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="$emit('update:visible', $event)"
    :title="$t('dialog.changeGroupTitle')"
    width="400px"
    :before-close="handleClose"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
  >
    <div class="change-group-form">
      <el-form :model="form" label-width="80px" label-position="left">
        <el-form-item :label="$t('dialog.accountTitleLabel')">
          <el-input v-model="form.title" readonly />
        </el-form-item>
        
        <el-form-item :label="$t('dialog.selectGroupLabel')" required>
          <el-select
            v-model="form.groupId"
            :placeholder="$t('dialog.selectGroupPlaceholder')"
            style="width: 100%"
            @change="handleGroupChange"
          >
            <el-option
              v-for="group in groups"
              :key="group.id"
              :label="group.name"
              :value="group.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item :label="$t('dialog.selectTypeLabel')" required>
          <el-select
            v-model="form.typeId"
            :placeholder="$t('dialog.selectTypePlaceholder')"
            style="width: 100%"
            :loading="typesLoading"
            @change="handleTypeChange"
          >
            <el-option
              v-for="type in types"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
    </div>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleConfirm" :disabled="!canConfirm">{{ $t('common.confirm') }}</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { apiService } from '@/services/api'

/**
 * 更改分组对话框组件
 * @author 陈凤庆
 * @date 20251005
 * @description 用于更改账号的分组和分类
 */
export default {
  name: 'ChangeGroupDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    account: {
      type: Object,
      default: null
    },
    groups: {
      type: Array,
      default: () => []
    }
  },
  emits: ['update:visible', 'confirm'],
  setup(props, { emit }) {
    // 表单数据
    const form = ref({
      title: '',
      groupId: '',
      typeId: ''
    })
    
    // 分类列表和加载状态
    const types = ref([])
    const typesLoading = ref(false)
    
    // 上一次选择的分组ID和分类ID（用于记忆功能）
    const lastSelectedGroupId = ref('')
    const lastSelectedTypeId = ref('')
    
    /**
     * 是否可以确认
     */
    const canConfirm = computed(() => {
      return form.value.groupId && form.value.typeId
    })
    
    /**
     * 监听账号变化，初始化表单
     */
    watch(() => props.account, (newAccount) => {
      if (newAccount) {
        form.value.title = newAccount.title

        // 默认显示上一次选择的分组，如果没有则使用第一个分组
        let defaultGroupId = lastSelectedGroupId.value
        if (!defaultGroupId && props.groups.length > 0) {
          defaultGroupId = props.groups[0].id
        }

        form.value.groupId = defaultGroupId
        form.value.typeId = ''

        // 如果有默认分组，加载对应的分类列表
        if (defaultGroupId) {
          loadTypes(defaultGroupId, true) // 传入 true 表示需要恢复上一次选择的分类
        }
      }
    }, { immediate: true })
    
    /**
     * 处理分组变化
     */
    const handleGroupChange = (groupId) => {
      console.log('[ChangeGroupDialog] 分组变化:', groupId)
      form.value.typeId = '' // 清空分类选择
      if (groupId) {
        loadTypes(groupId, false) // 分组变化时不恢复上一次的分类选择
        lastSelectedGroupId.value = groupId // 记住选择的分组
      } else {
        types.value = []
      }
    }

    /**
     * 处理分类变化
     */
    const handleTypeChange = (typeId) => {
      console.log('[ChangeGroupDialog] 分类变化:', typeId)
      if (typeId) {
        lastSelectedTypeId.value = typeId // 记住选择的分类
      }
    }
    
    /**
     * 加载分类列表
     * @param {string} groupId 分组ID
     * @param {boolean} restoreLastType 是否恢复上一次选择的分类
     */
    const loadTypes = async (groupId, restoreLastType = false) => {
      if (!groupId) return

      try {
        typesLoading.value = true
        console.log('[ChangeGroupDialog] 加载分类列表，分组ID:', groupId, '恢复上次分类:', restoreLastType)

        const typeList = await apiService.getTypesByGroupId(groupId)
        types.value = typeList || []

        console.log('[ChangeGroupDialog] 分类列表加载完成:', types.value.length, '个分类')

        // 如果需要恢复上一次选择的分类，并且上一次的分类在当前分组中存在
        if (restoreLastType && lastSelectedTypeId.value) {
          const typeExists = types.value.some(type => type.id === lastSelectedTypeId.value)
          if (typeExists) {
            form.value.typeId = lastSelectedTypeId.value
            console.log('[ChangeGroupDialog] 恢复上一次选择的分类:', lastSelectedTypeId.value)
          } else {
            console.log('[ChangeGroupDialog] 上一次选择的分类在当前分组中不存在，无法恢复')
          }
        }
      } catch (error) {
        console.error('[ChangeGroupDialog] 加载分类列表失败:', error)
        ElMessage.error(this.$t('error.loadTypesFailed'))
        types.value = []
      } finally {
        typesLoading.value = false
      }
    }
    
    /**
     * 处理关闭
     */
    const handleClose = () => {
      emit('update:visible', false)
    }
    
    /**
     * 处理确认
     */
    const handleConfirm = () => {
      if (!canConfirm.value) {
        ElMessage.warning(this.$t('warning.selectGroupAndType'))
        return
      }
      
      const result = {
        accountId: props.account.id,
        groupId: form.value.groupId,
        typeId: form.value.typeId
      }
      
      console.log('[ChangeGroupDialog] 确认更改分组:', result)
      emit('confirm', result)
      handleClose()
    }
    
    return {
      form,
      types,
      typesLoading,
      canConfirm,
      handleGroupChange,
      handleTypeChange,
      handleClose,
      handleConfirm
    }
  }
}
</script>

<style scoped>
.change-group-form {
  padding: 10px 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.el-form-item {
  margin-bottom: 20px;
}

.el-input[readonly] {
  background-color: #f5f7fa;
  cursor: not-allowed;
}
</style>
