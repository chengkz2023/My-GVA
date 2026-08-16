import { createI18n } from 'vue-i18n'

// i18n 接缝（见 docs/i18n.md 与 CONTEXT.md「多语言接缝」）：
// 默认中文。新增 UI 文案必须走 t()（示例：view/superAdmin/dictionary/dictionary.vue）。
// 有海外交付需求的项目在此追加语言包与 locale 切换，并把存量硬编码文案抽取到对应语言包。
const messages = {
  'zh-CN': {
    dictionary: {
      type: '字典类型',
      name: '字典名称',
      search: '搜索',
      add: '新增字典',
      addDetail: '新增字典项',
      edit: '编辑',
      editDict: '编辑字典',
      editDetail: '编辑字典项',
      delete: '删除',
      details: '字典项',
      sort: '排序',
      status: '状态',
      enabled: '启用',
      disabled: '禁用',
      label: '显示值',
      value: '存储值',
      typePlaceholder: '如 gender',
      namePlaceholder: '如 性别',
      labelPlaceholder: '如 男',
      valuePlaceholder: '如 male',
      cancel: '取 消',
      confirm: '确 定',
      close: '关 闭',
      addSuccess: '添加成功',
      updateSuccess: '更新成功',
      deleteSuccess: '删除成功',
      deleteCanceled: '已取消删除',
      confirmDeleteDict: '删除字典「{name}」将同时删除其全部字典项，是否继续？',
      confirmDeleteDetail: '确定删除字典项「{label}」？',
      tip: '提示'
    }
  }
}

const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages
})

export default i18n
