<template>
  <div class="dictionary">
    <warning-bar title="字典管理：集中维护枚举型键值，业务模块通过 GET /dictionary/types 引用" />
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-input
          v-model="searchInfo.type"
          placeholder="字典类型"
          clearable
          style="width: 200px"
          @keyup.enter="getTableData"
          @clear="getTableData"
        />
        <el-button type="primary" icon="search" @click="getTableData">搜索</el-button>
        <el-button type="primary" icon="plus" @click="openDictForm()">新增字典</el-button>
      </div>
      <el-table :data="tableData" row-key="id" style="width: 100%">
        <el-table-column label="ID" min-width="80" prop="id" />
        <el-table-column label="字典类型" min-width="140" prop="type" />
        <el-table-column label="字典名称" min-width="140" prop="name" />
        <el-table-column label="排序" min-width="80" prop="sort" />
        <el-table-column label="状态" min-width="90">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="260">
          <template #default="scope">
            <el-button icon="tickets" type="primary" link @click="openDetailDrawer(scope.row)">字典项</el-button>
            <el-button icon="edit" type="primary" link @click="openDictForm(scope.row)">编辑</el-button>
            <el-button icon="delete" type="primary" link @click="removeDict(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="getTableData"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 字典表单抽屉 -->
    <el-drawer v-model="dictFormVisible" :size="appStore.drawerSize" :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dictFormTitle }}</span>
          <div>
            <el-button @click="dictFormVisible = false">取 消</el-button>
            <el-button type="primary" @click="submitDictForm">确 定</el-button>
          </div>
        </div>
      </template>
      <el-form ref="dictFormRef" :model="dictForm" :rules="dictRules" label-width="90px">
        <el-form-item label="字典类型" prop="type">
          <el-input v-model="dictForm.type" :disabled="dictFormType === 'edit'" autocomplete="off" placeholder="如 gender" />
        </el-form-item>
        <el-form-item label="字典名称" prop="name">
          <el-input v-model="dictForm.name" autocomplete="off" placeholder="如 性别" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="dictForm.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="dictForm.status" :active-value="1" :inactive-value="2" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- 字典项抽屉 -->
    <el-drawer v-model="detailDrawerVisible" :size="appStore.drawerSize" :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ activeDict.name }}（{{ activeDict.type }}）字典项</span>
          <div>
            <el-button @click="detailDrawerVisible = false">关 闭</el-button>
            <el-button type="primary" icon="plus" @click="openDetailForm()">新增字典项</el-button>
          </div>
        </div>
      </template>
      <el-table :data="detailData" row-key="id" style="width: 100%">
        <el-table-column label="显示值" min-width="120" prop="label" />
        <el-table-column label="存储值" min-width="120" prop="value" />
        <el-table-column label="排序" min-width="80" prop="sort" />
        <el-table-column label="状态" min-width="90">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="140">
          <template #default="scope">
            <el-button icon="edit" type="primary" link @click="openDetailForm(scope.row)">编辑</el-button>
            <el-button icon="delete" type="primary" link @click="removeDetail(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>

    <!-- 字典项表单弹窗 -->
    <el-dialog v-model="detailFormVisible" :title="detailFormTitle" width="480px">
      <el-form ref="detailFormRef" :model="detailForm" :rules="detailRules" label-width="90px">
        <el-form-item label="显示值" prop="label">
          <el-input v-model="detailForm.label" autocomplete="off" placeholder="如 男" />
        </el-form-item>
        <el-form-item label="存储值" prop="value">
          <el-input v-model="detailForm.value" autocomplete="off" placeholder="如 male" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="detailForm.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="detailForm.status" :active-value="1" :inactive-value="2" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="detailFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="submitDetailForm">确 定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import {
    getDictionaryList,
    createDictionary,
    updateDictionary,
    deleteDictionary,
    getDictionaryDetails,
    createDictionaryDetail,
    updateDictionaryDetail,
    deleteDictionaryDetail
  } from '@/api/dictionary'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { useAppStore } from '@/pinia'
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'

  defineOptions({
    name: 'Dictionary'
  })

  const appStore = useAppStore()

  const tableData = ref([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const searchInfo = ref({ type: '' })

  const getTableData = async () => {
    const res = await getDictionaryList({ page: page.value, pageSize: pageSize.value, type: searchInfo.value.type })
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    }
  }
  getTableData()

  const handleSizeChange = () => {
    page.value = 1
    getTableData()
  }

  // ---- 字典表单 ----
  const dictFormVisible = ref(false)
  const dictFormType = ref('add')
  const dictFormTitle = ref('新增字典')
  const dictFormRef = ref(null)
  const dictForm = ref({ id: 0, type: '', name: '', sort: 0, status: 1 })
  const dictRules = {
    type: [{ required: true, message: '请输入字典类型', trigger: 'blur' }],
    name: [{ required: true, message: '请输入字典名称', trigger: 'blur' }]
  }

  const openDictForm = (row) => {
    dictFormType.value = row ? 'edit' : 'add'
    dictFormTitle.value = row ? '编辑字典' : '新增字典'
    dictForm.value = row
      ? { id: row.id, type: row.type, name: row.name, sort: row.sort, status: row.status }
      : { id: 0, type: '', name: '', sort: 0, status: 1 }
    dictFormVisible.value = true
  }

  const submitDictForm = () => {
    dictFormRef.value.validate(async (valid) => {
      if (!valid) {
        return
      }
      if (dictFormType.value === 'add') {
        const res = await createDictionary(dictForm.value)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '添加成功' })
          dictFormVisible.value = false
          getTableData()
        }
      } else {
        const res = await updateDictionary(dictForm.value.id, dictForm.value)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '更新成功' })
          dictFormVisible.value = false
          getTableData()
        }
      }
    })
  }

  const removeDict = (row) => {
    ElMessageBox.confirm(`删除字典「${row.name}」将同时删除其全部字典项，是否继续？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
      .then(async () => {
        const res = await deleteDictionary(row.id)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '删除成功' })
          getTableData()
        }
      })
      .catch(() => {
        ElMessage({ type: 'info', message: '已取消删除' })
      })
  }

  // ---- 字典项 ----
  const detailDrawerVisible = ref(false)
  const activeDict = ref({})
  const detailData = ref([])

  const getDetailData = async () => {
    const res = await getDictionaryDetails({ dictionaryId: activeDict.value.id })
    if (res.code === 0) {
      detailData.value = res.data.list || []
    }
  }

  const openDetailDrawer = (row) => {
    activeDict.value = row
    detailDrawerVisible.value = true
    getDetailData()
  }

  const detailFormVisible = ref(false)
  const detailFormType = ref('add')
  const detailFormTitle = ref('新增字典项')
  const detailFormRef = ref(null)
  const detailForm = ref({ id: 0, dictionaryId: 0, label: '', value: '', sort: 0, status: 1 })
  const detailRules = {
    label: [{ required: true, message: '请输入显示值', trigger: 'blur' }],
    value: [{ required: true, message: '请输入存储值', trigger: 'blur' }]
  }

  const openDetailForm = (row) => {
    detailFormType.value = row ? 'edit' : 'add'
    detailFormTitle.value = row ? '编辑字典项' : '新增字典项'
    detailForm.value = row
      ? { id: row.id, dictionaryId: activeDict.value.id, label: row.label, value: row.value, sort: row.sort, status: row.status }
      : { id: 0, dictionaryId: activeDict.value.id, label: '', value: '', sort: 0, status: 1 }
    detailFormVisible.value = true
  }

  const submitDetailForm = () => {
    detailFormRef.value.validate(async (valid) => {
      if (!valid) {
        return
      }
      if (detailFormType.value === 'add') {
        const res = await createDictionaryDetail(detailForm.value)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '添加成功' })
          detailFormVisible.value = false
          getDetailData()
        }
      } else {
        const res = await updateDictionaryDetail(detailForm.value.id, detailForm.value)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '更新成功' })
          detailFormVisible.value = false
          getDetailData()
        }
      }
    })
  }

  const removeDetail = (row) => {
    ElMessageBox.confirm(`确定删除字典项「${row.label}」？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
      .then(async () => {
        const res = await deleteDictionaryDetail(row.id)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '删除成功' })
          getDetailData()
        }
      })
      .catch(() => {
        ElMessage({ type: 'info', message: '已取消删除' })
      })
  }
</script>

<style lang="scss" scoped>
  .gva-pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: 16px;
  }
</style>
