// No V2 endpoints — attachment category not migrated
import service from '@/utils/request'
export const getCategoryList = () => Promise.resolve({ code: 0, data: { list: [] } })
export const addCategory = () => Promise.resolve({ code: 0 })
export const deleteCategory = () => Promise.resolve({ code: 0 })
