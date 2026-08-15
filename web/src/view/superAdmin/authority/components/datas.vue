<template>
  <div>
    <warning-bar
      title="该功能用于维护角色之间的数据权限关系，实际业务请结合自己的表结构和权限规则实现。"
      href="javascript:void(0)"
    />
    <div class="sticky top-0.5 z-10 my-4">
      <el-button class="float-left" type="primary" @click="all">全选</el-button>
      <el-button class="float-left" type="primary" @click="self">当前角色</el-button>
      <el-button class="float-left" type="primary" @click="selfAndChildren">
        当前角色及子角色
      </el-button>
      <el-button class="float-right" type="primary" @click="authDataEnter">
        确定
      </el-button>
    </div>
    <div class="clear-both pt-4">
      <el-checkbox-group v-model="dataAuthorityId" @change="selectAuthority">
        <el-checkbox
          v-for="(item, key) in authoritys"
          :key="key"
          :label="item"
        >
          {{ item.authorityName }}
        </el-checkbox>
      </el-checkbox-group>
    </div>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import { setDataAuthority, getAuthorityInfo } from '@/api/authority'
  import WarningBar from '@/components/warningBar/warningBar.vue'

  defineOptions({
    name: 'Datas'
  })

  const props = defineProps({
    row: {
      type: Object,
      default: () => ({})
    },
    authority: {
      type: Array,
      default: () => []
    }
  })

  const emit = defineEmits(['changeRow'])

  const authoritys = ref([])
  const dataAuthorityId = ref([])
  const needConfirm = ref(false)

  const roundAuthority = (authoritysData) => {
    authoritysData?.forEach((item) => {
      authoritys.value.push({
        authorityId: item.authorityId,
        authorityName: item.authorityName
      })
      if (item.children?.length) {
        roundAuthority(item.children)
      }
    })
  }

  const init = async () => {
    roundAuthority(props.authority)
    // 回显：GET /system/role/:id 返回 dataAuthorityIds（[]uint）
    const res = await getAuthorityInfo({ authorityId: props.row.authorityId })
    const authIds = res.code === 0 ? res.data.dataAuthorityIds || [] : []
    authIds.forEach((authId) => {
      const target = authoritys.value.find(
        (authority) => authority.authorityId === authId
      )
      if (target) {
        dataAuthorityId.value.push(target)
      }
    })
  }

  init()

  const enterAndNext = () => {
    authDataEnter()
  }

  const all = () => {
    dataAuthorityId.value = [...authoritys.value]
    emit('changeRow', 'dataAuthorityId', dataAuthorityId.value)
    needConfirm.value = true
  }

  const self = () => {
    dataAuthorityId.value = authoritys.value.filter(
      (item) => item.authorityId === props.row.authorityId
    )
    emit('changeRow', 'dataAuthorityId', dataAuthorityId.value)
    needConfirm.value = true
  }

  const getChildrenId = (row, ids) => {
    ids.push(row.authorityId)
    row.children?.forEach((item) => {
      getChildrenId(item, ids)
    })
  }

  const selfAndChildren = () => {
    const ids = []
    getChildrenId(props.row, ids)
    dataAuthorityId.value = authoritys.value.filter((item) =>
      ids.includes(item.authorityId)
    )
    emit('changeRow', 'dataAuthorityId', dataAuthorityId.value)
    needConfirm.value = true
  }

  const authDataEnter = async () => {
    // 后端期望 { authorityId, dataAuthorityIds: []uint }
    const res = await setDataAuthority({
      authorityId: props.row.authorityId,
      dataAuthorityIds: dataAuthorityId.value
        .filter(Boolean)
        .map((item) => item.authorityId)
    })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '数据权限设置成功' })
      needConfirm.value = false
    }
  }

  const selectAuthority = () => {
    dataAuthorityId.value = dataAuthorityId.value.filter(Boolean)
    emit('changeRow', 'dataAuthorityId', dataAuthorityId.value)
    needConfirm.value = true
  }

  defineExpose({
    enterAndNext,
    needConfirm
  })
</script>
