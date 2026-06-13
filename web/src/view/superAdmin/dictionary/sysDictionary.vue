<template>
  <div class="dictionary-page">
    <div class="dictionary-sidebar">
      <div class="panel-card">
        <warning-bar title="字典变更后，前端依赖字典的下拉与展示内容会随之更新。" />
        <div class="toolbar">
          <el-input
            v-model.trim="searchName"
            placeholder="搜索字典名称"
            clearable
            @keyup.enter="getTableData"
            @clear="getTableData"
          >
            <template #append>
              <el-button @click="getTableData">查询</el-button>
            </template>
          </el-input>
          <div class="toolbar-actions">
            <el-button type="success" :icon="Upload" @click="openImportDialog">导入</el-button>
            <el-button type="primary" :icon="Plus" @click="openDrawer">新增</el-button>
          </div>
        </div>

        <el-scrollbar class="dict-list">
          <div
            v-for="dictionary in dictionaryData"
            :key="dictionary.ID"
            class="dict-item"
            :class="{ active: selectID === dictionary.ID, child: !!dictionary.parentID }"
            @click="toDetail(dictionary)"
          >
            <div class="dict-main">
              <div class="dict-name">{{ dictionary.name }}</div>
              <div class="dict-type">{{ dictionary.type }}</div>
            </div>
            <div class="dict-actions">
              <el-button link type="success" @click.stop="exportDictionary(dictionary)">导出</el-button>
              <el-button link type="primary" @click.stop="updateSysDictionaryFunc(dictionary)">编辑</el-button>
              <el-button link type="danger" @click.stop="deleteSysDictionaryFunc(dictionary)">删除</el-button>
            </div>
          </div>
          <el-empty v-if="!dictionaryData.length" description="暂无字典数据" />
        </el-scrollbar>
      </div>
    </div>

    <div class="dictionary-detail">
      <div class="panel-card detail-panel">
        <sysDictionaryDetail :sys-dictionary-i-d="selectID" />
      </div>
    </div>

    <el-drawer
      v-model="drawerFormVisible"
      :size="appStore.drawerSize"
      :show-close="false"
      :before-close="closeDrawer"
    >
      <template #header>
        <div class="drawer-header">
          <span class="text-lg">{{ type === 'update' ? '编辑字典' : '新增字典' }}</span>
          <div>
            <el-button @click="closeDrawer">取消</el-button>
            <el-button type="primary" @click="enterDrawer">确认</el-button>
          </div>
        </div>
      </template>

      <el-form ref="drawerForm" :model="formData" :rules="rules" label-width="110px">
        <el-form-item label="父级字典" prop="parentID">
          <el-select v-model="formData.parentID" placeholder="可选" clearable filterable style="width: 100%">
            <el-option
              v-for="dict in availableParentDictionaries"
              :key="dict.ID"
              :label="`${dict.name} (${dict.type})`"
              :value="dict.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="字典名称" prop="name">
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item label="字典类型" prop="type">
          <el-input v-model="formData.type" />
        </el-form-item>
        <el-form-item label="启用状态" prop="status">
          <el-switch v-model="formData.status" active-text="启用" inactive-text="停用" />
        </el-form-item>
        <el-form-item label="描述" prop="desc">
          <el-input v-model="formData.desc" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer
      v-model="importDrawerVisible"
      :size="appStore.drawerSize"
      :show-close="false"
      :before-close="closeImportDrawer"
    >
      <template #header>
        <div class="drawer-header">
          <span class="text-lg">导入字典 JSON</span>
          <div>
            <el-button @click="closeImportDrawer">取消</el-button>
            <el-button type="primary" :loading="importing" @click="handleImport">确认导入</el-button>
          </div>
        </div>
      </template>

      <el-alert
        title="支持直接粘贴 JSON，也支持拖拽或点击上传 .json 文件。"
        type="info"
        :closable="false"
        show-icon
      />

      <div
        class="drag-upload-area"
        :class="{ 'is-dragging': isDragging }"
        @drop.prevent="handleDrop"
        @dragover.prevent="handleDragOver"
        @dragleave.prevent="handleDragLeave"
        @click="triggerFileInput"
      >
        <el-icon class="upload-icon"><Upload /></el-icon>
        <p>拖拽 JSON 文件到这里，或点击选择文件</p>
        <p class="upload-hint">文件内容会自动填充到下方编辑区</p>
        <input
          ref="fileInputRef"
          type="file"
          accept=".json,application/json"
          style="display: none"
          @change="handleFileSelect"
        />
      </div>

      <div class="json-editor-container">
        <el-input
          v-model="importJsonText"
          type="textarea"
          :rows="16"
          class="json-textarea"
          placeholder='请输入字典 JSON，例如 {"name":"状态","type":"status","status":true,"desc":"系统状态","sysDictionaryDetails":[]}'
        />
      </div>

      <el-alert
        v-if="jsonPreviewError"
        :title="jsonPreviewError"
        type="error"
        :closable="false"
        show-icon
        class="mt-4"
      />
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createSysDictionary,
  deleteSysDictionary,
  exportSysDictionary,
  findSysDictionary,
  getSysDictionaryList,
  importSysDictionary,
  updateSysDictionary
} from '@/api/sysDictionary'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { useAppStore } from '@/pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Upload } from '@element-plus/icons-vue'
import { ref, watch } from 'vue'
import sysDictionaryDetail from './sysDictionaryDetail.vue'

defineOptions({
  name: 'SysDictionary'
})

const appStore = useAppStore()

const selectID = ref(0)
const searchName = ref('')
const type = ref('create')
const dictionaryData = ref([])
const availableParentDictionaries = ref([])
const drawerFormVisible = ref(false)
const drawerForm = ref()
const importDrawerVisible = ref(false)
const importJsonText = ref('')
const importing = ref(false)
const jsonPreviewError = ref('')
const jsonPreview = ref(null)
const isDragging = ref(false)
const fileInputRef = ref()
const formData = ref({
  name: '',
  type: '',
  status: true,
  desc: '',
  parentID: null
})

const rules = {
  name: [{ required: true, message: '请输入字典名称', trigger: 'blur' }],
  type: [{ required: true, message: '请输入字典类型', trigger: 'blur' }],
  desc: [{ required: true, message: '请输入描述', trigger: 'blur' }]
}

watch(importJsonText, (value) => {
  if (!value.trim()) {
    jsonPreview.value = null
    jsonPreviewError.value = ''
    return
  }
  try {
    jsonPreview.value = JSON.parse(value)
    jsonPreviewError.value = ''
  } catch (error) {
    jsonPreview.value = null
    jsonPreviewError.value = `JSON 格式错误: ${error.message}`
  }
})

const resetFormData = () => {
  formData.value = {
    name: '',
    type: '',
    status: true,
    desc: '',
    parentID: null
  }
}

const isChildDictionary = (dictId, parentId) => {
  const dict = dictionaryData.value.find((item) => item.ID === dictId)
  if (!dict || !dict.parentID) return false
  if (dict.parentID === parentId) return true
  return isChildDictionary(dict.parentID, parentId)
}

const updateAvailableParentDictionaries = () => {
  if (type.value === 'update' && formData.value.ID) {
    availableParentDictionaries.value = dictionaryData.value.filter((dict) => {
      return dict.ID !== formData.value.ID && !isChildDictionary(dict.ID, formData.value.ID)
    })
    return
  }
  availableParentDictionaries.value = [...dictionaryData.value]
}

const getTableData = async () => {
  const res = await getSysDictionaryList({ name: searchName.value.trim() })
  if (res.code === 0) {
    dictionaryData.value = res.data || []
    if (dictionaryData.value.length) {
      if (!dictionaryData.value.some((item) => item.ID === selectID.value)) {
        selectID.value = dictionaryData.value[0].ID
      }
    } else {
      selectID.value = 0
    }
    updateAvailableParentDictionaries()
  }
}

const toDetail = (row) => {
  selectID.value = row.ID
}

const openDrawer = () => {
  type.value = 'create'
  resetFormData()
  drawerForm.value?.clearValidate()
  updateAvailableParentDictionaries()
  drawerFormVisible.value = true
}

const closeDrawer = () => {
  drawerFormVisible.value = false
  drawerForm.value?.clearValidate()
  resetFormData()
}

const updateSysDictionaryFunc = async (row) => {
  const res = await findSysDictionary({ ID: row.ID })
  if (res.code === 0) {
    type.value = 'update'
    formData.value = {
      ...(res.data.resysDictionary || {}),
      status: res.data.resysDictionary.status ?? true
    }
    updateAvailableParentDictionaries()
    drawerFormVisible.value = true
  }
}

const enterDrawer = () => {
  drawerForm.value.validate(async (valid) => {
    if (!valid) return
    const req = { ...formData.value }
    const res = type.value === 'update'
      ? await updateSysDictionary(req)
      : await createSysDictionary(req)
    if (res.code === 0) {
      ElMessage.success(type.value === 'update' ? '更新成功' : '创建成功')
      closeDrawer()
      await getTableData()
    }
  })
}

const deleteSysDictionaryFunc = async (row) => {
  ElMessageBox.confirm('确定删除这个字典吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteSysDictionary({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      await getTableData()
    }
  })
}

const exportDictionary = async (row) => {
  const res = await exportSysDictionary({ ID: row.ID })
  if (res.code !== 0) return
  const jsonStr = JSON.stringify(res.data, null, 2)
  const blob = new Blob([jsonStr], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `${row.type || row.name}-dictionary.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
  ElMessage.success('导出成功')
}

const openImportDialog = () => {
  importDrawerVisible.value = true
  importJsonText.value = ''
  jsonPreview.value = null
  jsonPreviewError.value = ''
  isDragging.value = false
}

const closeImportDrawer = () => {
  importDrawerVisible.value = false
  importJsonText.value = ''
  jsonPreview.value = null
  jsonPreviewError.value = ''
  isDragging.value = false
}

const handleDragOver = () => {
  isDragging.value = true
}

const handleDragLeave = () => {
  isDragging.value = false
}

const handleDrop = (event) => {
  isDragging.value = false
  const file = event.dataTransfer.files?.[0]
  if (file) {
    readJsonFile(file)
  }
}

const triggerFileInput = () => {
  fileInputRef.value?.click()
}

const handleFileSelect = (event) => {
  const file = event.target.files?.[0]
  if (file) {
    readJsonFile(file)
  }
  event.target.value = ''
}

const readJsonFile = (file) => {
  if (!file.name.endsWith('.json')) {
    ElMessage.warning('请选择 JSON 文件')
    return
  }
  const reader = new FileReader()
  reader.onload = (event) => {
    try {
      const content = event.target?.result || ''
      JSON.parse(content)
      importJsonText.value = content
      ElMessage.success('文件读取成功')
    } catch {
      ElMessage.error('文件内容不是有效的 JSON')
    }
  }
  reader.onerror = () => {
    ElMessage.error('文件读取失败')
  }
  reader.readAsText(file)
}

const handleImport = async () => {
  if (!importJsonText.value.trim()) {
    ElMessage.warning('请输入或上传 JSON 内容')
    return
  }
  if (jsonPreviewError.value) {
    ElMessage.error('请先修正 JSON 格式错误')
    return
  }
  importing.value = true
  try {
    const res = await importSysDictionary({ json: importJsonText.value })
    if (res.code === 0) {
      ElMessage.success('导入成功')
      closeImportDrawer()
      await getTableData()
    }
  } finally {
    importing.value = false
  }
}

getTableData()
</script>

<style scoped>
.dictionary-page {
  display: flex;
  gap: 16px;
  height: 100%;
}

.dictionary-sidebar {
  width: 360px;
  min-width: 320px;
}

.dictionary-detail {
  flex: 1;
  min-width: 0;
}

.panel-card {
  height: 100%;
  padding: 16px;
  background: #fff;
  border-radius: 8px;
}

.detail-panel {
  overflow: hidden;
}

.toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.toolbar-actions {
  display: flex;
  gap: 8px;
}

.dict-list {
  height: calc(100vh - 280px);
  margin-top: 16px;
}

.dict-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 12px;
  margin-top: 8px;
  border-radius: 8px;
  background: #f8fafc;
  cursor: pointer;
}

.dict-item.child {
  margin-left: 16px;
  border-left: 2px solid #bfdbfe;
}

.dict-item.active {
  outline: 1px solid var(--el-color-primary);
}

.dict-main {
  min-width: 0;
  flex: 1;
}

.dict-name {
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dict-type {
  margin-top: 4px;
  font-size: 12px;
  color: #64748b;
}

.dict-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.drag-upload-area {
  margin-top: 16px;
  padding: 36px 20px;
  text-align: center;
  border: 2px dashed #dcdfe6;
  border-radius: 8px;
  background: #fafafa;
  cursor: pointer;
  transition: all 0.3s ease;
}

.drag-upload-area:hover,
.drag-upload-area.is-dragging {
  border-color: #409eff;
  background: #ecf5ff;
}

.upload-icon {
  margin-bottom: 12px;
  font-size: 44px;
  color: #8c939d;
}

.upload-hint {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}

.json-editor-container {
  margin-top: 16px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

.json-textarea :deep(.el-textarea__inner) {
  font-family: Consolas, 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
}
</style>
