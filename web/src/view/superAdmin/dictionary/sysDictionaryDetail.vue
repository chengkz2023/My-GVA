<template>
  <div class="detail-page">
    <div class="gva-table-box">
      <div class="gva-btn-list justify-between flex items-center">
        <span class="font-bold">字典详情</span>
        <div class="flex items-center gap-2">
          <el-button type="primary" icon="plus" :disabled="!props.sysDictionaryID" @click="openDrawer">
            新增字典项
          </el-button>
        </div>
      </div>

      <el-empty v-if="!props.sysDictionaryID" description="请先选择左侧字典" />

      <el-table
        v-else
        :data="treeData"
        row-key="ID"
        default-expand-all
        :tree-props="{ children: 'children' }"
      >
        <el-table-column align="left" label="显示值" prop="label" min-width="160" />
        <el-table-column align="left" label="字典值" prop="value" min-width="140" />
        <el-table-column align="left" label="扩展值" prop="extend" min-width="140" />
        <el-table-column align="left" label="层级" prop="level" width="80" />
        <el-table-column align="left" label="状态" prop="status" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status ? 'success' : 'info'">
              {{ row.status ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="排序" prop="sort" width="90" />
        <el-table-column align="left" label="操作" :min-width="appStore.operateMinWith" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link icon="plus" @click="addChildNode(row)">新增子项</el-button>
            <el-button type="primary" link icon="edit" @click="updateSysDictionaryDetailFunc(row)">编辑</el-button>
            <el-button type="danger" link icon="delete" @click="deleteSysDictionaryDetailFunc(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-drawer
      v-model="drawerFormVisible"
      :size="appStore.drawerSize"
      :show-close="false"
      :before-close="closeDrawer"
    >
      <template #header>
        <div class="drawer-header">
          <span class="text-lg">{{ type === 'create' ? '新增字典项' : '编辑字典项' }}</span>
          <div>
            <el-button @click="closeDrawer">取消</el-button>
            <el-button type="primary" @click="enterDrawer">确认</el-button>
          </div>
        </div>
      </template>

      <el-form ref="drawerForm" :model="formData" :rules="rules" label-width="110px">
        <el-form-item label="父级节点" prop="parentID">
          <el-cascader
            v-model="formData.parentID"
            :options="[rootOption, ...treeData]"
            :props="cascadeProps"
            clearable
            filterable
            placeholder="可选，不选则为根节点"
            style="width: 100%"
            @change="handleParentChange"
          />
        </el-form-item>
        <el-form-item label="显示值" prop="label">
          <el-input v-model="formData.label" />
        </el-form-item>
        <el-form-item label="字典值" prop="value">
          <el-input v-model="formData.value" />
        </el-form-item>
        <el-form-item label="扩展值" prop="extend">
          <el-input v-model="formData.extend" />
        </el-form-item>
        <el-form-item label="启用状态" prop="status">
          <el-switch v-model="formData.status" active-text="启用" inactive-text="停用" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createSysDictionaryDetail,
  deleteSysDictionaryDetail,
  findSysDictionaryDetail,
  getDictionaryTreeList,
  updateSysDictionaryDetail
} from '@/api/sysDictionaryDetail'
import { useAppStore } from '@/pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, watch } from 'vue'

defineOptions({
  name: 'SysDictionaryDetail'
})

const props = defineProps({
  sysDictionaryID: {
    type: Number,
    default: 0
  }
})

const appStore = useAppStore()
const drawerForm = ref()
const drawerFormVisible = ref(false)
const type = ref('create')
const treeData = ref([])
const formData = ref({
  label: '',
  value: '',
  extend: '',
  status: true,
  sort: 0,
  parentID: null,
  sysDictionaryID: 0
})

const rules = {
  label: [{ required: true, message: '请输入显示值', trigger: 'blur' }],
  value: [{ required: true, message: '请输入字典值', trigger: 'blur' }],
  sort: [{ required: true, message: '请输入排序值', trigger: 'blur' }]
}

const cascadeProps = {
  value: 'ID',
  label: 'label',
  children: 'children',
  checkStrictly: true,
  emitPath: false
}

const rootOption = {
  ID: null,
  label: '无父级（根节点）'
}

const resetForm = () => {
  formData.value = {
    label: '',
    value: '',
    extend: '',
    status: true,
    sort: 0,
    parentID: null,
    sysDictionaryID: props.sysDictionaryID
  }
}

const getTreeData = async () => {
  if (!props.sysDictionaryID) {
    treeData.value = []
    return
  }
  const res = await getDictionaryTreeList({
    sysDictionaryID: props.sysDictionaryID
  })
  if (res.code === 0) {
    treeData.value = res.data.list || []
  }
}

const openDrawer = () => {
  type.value = 'create'
  resetForm()
  drawerForm.value?.clearValidate()
  drawerFormVisible.value = true
}

const addChildNode = (parentNode) => {
  type.value = 'create'
  resetForm()
  formData.value.parentID = parentNode.ID
  drawerForm.value?.clearValidate()
  drawerFormVisible.value = true
}

const updateSysDictionaryDetailFunc = async (row) => {
  const res = await findSysDictionaryDetail({ ID: row.ID })
  if (res.code === 0) {
    type.value = 'update'
    formData.value = {
      ...(res.data.reSysDictionaryDetail || {}),
      status: res.data.reSysDictionaryDetail.status ?? true,
      sysDictionaryID: props.sysDictionaryID
    }
    drawerFormVisible.value = true
  }
}

const handleParentChange = (value) => {
  formData.value.parentID = value
}

const closeDrawer = () => {
  drawerFormVisible.value = false
  drawerForm.value?.clearValidate()
  resetForm()
}

const enterDrawer = () => {
  drawerForm.value.validate(async (valid) => {
    if (!valid) return
    const payload = {
      ...formData.value,
      sysDictionaryID: props.sysDictionaryID
    }
    const res = type.value === 'update'
      ? await updateSysDictionaryDetail(payload)
      : await createSysDictionaryDetail(payload)
    if (res.code === 0) {
      ElMessage.success(type.value === 'update' ? '更新成功' : '创建成功')
      closeDrawer()
      await getTreeData()
    }
  })
}

const deleteSysDictionaryDetailFunc = async (row) => {
  ElMessageBox.confirm('确定删除这个字典项吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteSysDictionaryDetail({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      await getTreeData()
    }
  })
}

watch(
  () => props.sysDictionaryID,
  () => {
    resetForm()
    getTreeData()
  },
  { immediate: true }
)
</script>

<style scoped>
.detail-page {
  height: 100%;
}

.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
