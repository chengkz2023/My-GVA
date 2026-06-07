<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="路径">
          <el-input v-model="searchInfo.path" placeholder="请输入路径" clearable />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="searchInfo.description" placeholder="请输入描述" clearable />
        </el-form-item>
        <el-form-item label="分组">
          <el-select v-model="searchInfo.apiGroup" placeholder="请选择分组" clearable>
            <el-option
              v-for="item in apiGroupOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="请求方式">
          <el-select v-model="searchInfo.method" placeholder="请选择请求方式" clearable>
            <el-option
              v-for="item in methodOptions"
              :key="item.value"
              :label="`${item.label} (${item.value})`"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog('addApi')">新增</el-button>
        <el-button icon="delete" :disabled="!apis.length" @click="onDelete">删除</el-button>
        <el-button icon="refresh" @click="onFresh">刷新缓存</el-button>
        <el-button icon="compass" @click="onSync">同步 API</el-button>
      </div>

      <el-table
        :data="tableData"
        @sort-change="sortChange"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" min-width="80" prop="ID" sortable="custom" />
        <el-table-column align="left" label="路径" min-width="180" prop="path" sortable="custom" />
        <el-table-column align="left" label="分组" min-width="140" prop="apiGroup" sortable="custom" />
        <el-table-column align="left" label="描述" min-width="180" prop="description" sortable="custom" />
        <el-table-column align="left" label="请求方式" min-width="120" prop="method" sortable="custom">
          <template #default="{ row }">
            <span>{{ row.method }} / {{ methodFilter(row.method) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          fixed="right"
          label="操作"
          :min-width="appStore.operateMinWith"
        >
          <template #default="{ row }">
            <el-button icon="edit" type="primary" link @click="editApiFunc(row)">编辑</el-button>
            <el-button icon="delete" type="primary" link @click="deleteApiFunc(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-drawer
      v-model="syncApiFlag"
      :size="appStore.drawerSize"
      :show-close="false"
      :before-close="closeSyncDialog"
    >
      <template #header>
        <div class="drawer-header">
          <span class="text-lg">同步 API</span>
          <div>
            <el-button @click="closeSyncDialog">取消</el-button>
            <el-button type="primary" :loading="syncing" @click="enterSyncDialog">确认</el-button>
          </div>
        </div>
      </template>

      <warning-bar title="同步结果会把当前路由与 API 表进行对齐，缺失项会新增，失效项会移除。" />

      <h4 class="section-title">待新增</h4>
      <el-table v-loading="syncing" :data="syncApiData.newApis">
        <el-table-column align="left" label="路径" min-width="180" prop="path" />
        <el-table-column align="left" label="分组" min-width="160" prop="apiGroup">
          <template #default="{ row }">
            <el-select
              v-model="row.apiGroup"
              placeholder="请选择或新建分组"
              allow-create
              filterable
              default-first-option
            >
              <el-option
                v-for="item in apiGroupOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column align="left" label="描述" min-width="180" prop="description">
          <template #default="{ row }">
            <el-input v-model="row.description" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="请求方式" min-width="120" prop="method">
          <template #default="{ row }">
            <span>{{ row.method }} / {{ methodFilter(row.method) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="150" fixed="right">
          <template #default="{ row }">
            <el-button icon="plus" type="primary" link @click="addApiFunc(row)">单条新增</el-button>
            <el-button icon="minus" type="primary" link @click="ignoreApiFunc(row, true)">忽略</el-button>
          </template>
        </el-table-column>
      </el-table>

      <h4 class="section-title">待删除</h4>
      <el-table :data="syncApiData.deleteApis">
        <el-table-column align="left" label="路径" min-width="180" prop="path" />
        <el-table-column align="left" label="分组" min-width="160" prop="apiGroup" />
        <el-table-column align="left" label="描述" min-width="180" prop="description" />
        <el-table-column align="left" label="请求方式" min-width="120" prop="method">
          <template #default="{ row }">
            <span>{{ row.method }} / {{ methodFilter(row.method) }}</span>
          </template>
        </el-table-column>
      </el-table>

      <h4 class="section-title">已忽略</h4>
      <el-table :data="syncApiData.ignoreApis">
        <el-table-column align="left" label="路径" min-width="180" prop="path" />
        <el-table-column align="left" label="分组" min-width="160" prop="apiGroup" />
        <el-table-column align="left" label="描述" min-width="180" prop="description" />
        <el-table-column align="left" label="请求方式" min-width="120" prop="method">
          <template #default="{ row }">
            <span>{{ row.method }} / {{ methodFilter(row.method) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="120" fixed="right">
          <template #default="{ row }">
            <el-button icon="refresh-left" type="primary" link @click="ignoreApiFunc(row, false)">
              取消忽略
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>

    <el-drawer
      v-model="dialogFormVisible"
      :size="appStore.drawerSize"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="drawer-header">
          <span class="text-lg">{{ dialogTitle }}</span>
          <div>
            <el-button @click="closeDialog">取消</el-button>
            <el-button type="primary" @click="enterDialog">确认</el-button>
          </div>
        </div>
      </template>

      <warning-bar title="新增 API 后，请到角色管理里分配对应权限。" />
      <el-form ref="apiForm" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="路径" prop="path">
          <el-input v-model="form.path" autocomplete="off" />
        </el-form-item>
        <el-form-item label="请求方式" prop="method">
          <el-select v-model="form.method" placeholder="请选择" style="width: 100%">
            <el-option
              v-for="item in methodOptions"
              :key="item.value"
              :label="`${item.label} (${item.value})`"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="分组" prop="apiGroup">
          <el-select
            v-model="form.apiGroup"
            placeholder="请选择或新建分组"
            allow-create
            filterable
            default-first-option
            style="width: 100%"
          >
            <el-option
              v-for="item in apiGroupOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" autocomplete="off" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createApi,
  deleteApi,
  deleteApisByIds,
  enterSyncApi,
  freshCasbin,
  getApiById,
  getApiGroups,
  getApiList,
  ignoreApi,
  syncApi,
  updateApi
} from '@/api/api'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { useAppStore } from '@/pinia'
import { toSQLLine } from '@/utils/stringFun'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref } from 'vue'

defineOptions({
  name: 'Api'
})

const appStore = useAppStore()

const apis = ref([])
const apiForm = ref()
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const dialogFormVisible = ref(false)
const dialogTitle = ref('新增 API')
const type = ref('addApi')
const syncing = ref(false)
const syncApiFlag = ref(false)
const searchInfo = ref({})
const apiGroupOptions = ref([])
const apiGroupMap = ref({})
const syncApiData = ref({
  newApis: [],
  deleteApis: [],
  ignoreApis: []
})
const form = ref({
  path: '',
  apiGroup: '',
  method: '',
  description: ''
})

const methodOptions = ref([
  { value: 'POST', label: '新增' },
  { value: 'GET', label: '查询' },
  { value: 'PUT', label: '更新' },
  { value: 'DELETE', label: '删除' }
])

const rules = {
  path: [{ required: true, message: '请输入 API 路径', trigger: 'blur' }],
  apiGroup: [{ required: true, message: '请输入 API 分组', trigger: 'blur' }],
  method: [{ required: true, message: '请选择请求方式', trigger: 'change' }],
  description: [{ required: true, message: '请输入 API 描述', trigger: 'blur' }]
}

const methodFilter = (value) => methodOptions.value.find((item) => item.value === value)?.label || value

const resetForm = () => {
  form.value = {
    path: '',
    apiGroup: '',
    method: '',
    description: ''
  }
}

const getGroup = async () => {
  const res = await getApiGroups()
  if (res.code === 0) {
    const groups = res.data.groups || []
    apiGroupOptions.value = groups.map((item) => ({ label: item, value: item }))
    apiGroupMap.value = res.data.apiGroupMap || res.data.groups.reduce((acc, g) => { acc[g] = g; return acc }, {})
  }
}

const getTableData = async () => {
  const res = await getApiList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total || 0
    page.value = res.data.page || 1
    pageSize.value = res.data.pageSize || 10
  }
}

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {}
  page.value = 1
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const sortChange = ({ prop, order }) => {
  if (prop) {
    const field = prop === 'ID' ? 'id' : prop
    searchInfo.value.orderKey = toSQLLine(field)
    searchInfo.value.desc = order === 'descending'
  } else {
    delete searchInfo.value.orderKey
    delete searchInfo.value.desc
  }
  getTableData()
}

const handleSelectionChange = (val) => {
  apis.value = val
}

const openDialog = (key) => {
  type.value = key
  dialogTitle.value = key === 'edit' ? '编辑 API' : '新增 API'
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  apiForm.value?.clearValidate()
  resetForm()
}

const editApiFunc = async (row) => {
  const res = await getApiById(row.ID)
  if (res.code === 0) {
    form.value = { ...res.data }
    openDialog('edit')
  }
}

const enterDialog = () => {
  apiForm.value.validate(async (valid) => {
    if (!valid) return
    const request = { ...form.value }
    const res = type.value === 'edit' ? await updateApi(request.id, request) : await createApi(request)
    if (res.code === 0) {
      ElMessage.success(type.value === 'edit' ? '更新成功' : '新增成功')
      closeDialog()
      await getTableData()
      await getGroup()
    }
  })
}

const deleteApiFunc = async (row) => {
  ElMessageBox.confirm('删除后会同步移除相关角色的 API 权限，是否继续？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteApi(row.ID)
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === 1 && page.value > 1) {
        page.value -= 1
      }
      await getTableData()
      await getGroup()
    }
  })
}

const onDelete = async () => {
  ElMessageBox.confirm('确定删除选中的 API 吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const ids = apis.value.map((item) => item.ID)
    const res = await deleteApisByIds({ ids })
    if (res.code === 0) {
      ElMessage.success(res.msg)
      if (tableData.value.length === ids.length && page.value > 1) {
        page.value -= 1
      }
      await getTableData()
    }
  })
}

const onFresh = async () => {
  ElMessageBox.confirm('确定刷新 Casbin 缓存吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await freshCasbin()
    if (res.code === 0) {
      ElMessage.success(res.msg)
    }
  })
}

const ignoreApiFunc = async (row, flag) => {
  const res = await ignoreApi({ path: row.path, method: row.method, flag })
  if (res.code === 0) {
    ElMessage.success(res.msg)
    if (flag) {
      syncApiData.value.newApis = syncApiData.value.newApis.filter(
        (item) => !(item.path === row.path && item.method === row.method)
      )
      syncApiData.value.ignoreApis.push(row)
      return
    }
    syncApiData.value.ignoreApis = syncApiData.value.ignoreApis.filter(
      (item) => !(item.path === row.path && item.method === row.method)
    )
    syncApiData.value.newApis.push(row)
  }
}

const addApiFunc = async (row) => {
  if (!row.apiGroup) {
    ElMessage.error('请先选择 API 分组')
    return
  }
  if (!row.description) {
    ElMessage.error('请先填写 API 描述')
    return
  }
  const res = await createApi(row)
  if (res.code === 0) {
    ElMessage.success('新增成功，请到角色管理中分配权限')
    syncApiData.value.newApis = syncApiData.value.newApis.filter(
      (item) => !(item.path === row.path && item.method === row.method)
    )
    await getTableData()
    await getGroup()
  }
}

const onSync = async () => {
  const res = await syncApi()
  if (res.code === 0) {
    ;(res.data.newApis || []).forEach((item) => {
      item.apiGroup = apiGroupMap.value[item.path.split('/')[1]] || item.apiGroup
    })
    syncApiData.value = {
      newApis: res.data.newApis || [],
      deleteApis: res.data.deleteApis || [],
      ignoreApis: res.data.ignoreApis || []
    }
    syncApiFlag.value = true
  }
}

const closeSyncDialog = () => {
  syncApiFlag.value = false
}

const enterSyncDialog = async () => {
  const hasInvalidItem = syncApiData.value.newApis.some((item) => !item.apiGroup || !item.description)
  if (hasInvalidItem) {
    ElMessage.error('存在待同步 API 未填写分组或描述')
    return
  }
  syncing.value = true
  const res = await enterSyncApi(syncApiData.value)
  syncing.value = false
  if (res.code === 0) {
    ElMessage.success(res.msg)
    syncApiFlag.value = false
    await getTableData()
    await getGroup()
  }
}

getTableData()
getGroup()
</script>

<style scoped>
.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.section-title {
  margin: 20px 0 12px;
  font-size: 15px;
  font-weight: 600;
}
</style>
