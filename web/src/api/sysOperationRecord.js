import service from '@/utils/request'
export const getSysOperationRecordList = (data) => service({ url: '/system/operation-record/list', method: 'get', params: data })
export const deleteSysOperationRecord = (data) => service({ url: `/system/operation-record/${data.id||data.ID}`, method: 'delete' })
export const deleteSysOperationRecordByIds = (data) => service({ url: '/system/operation-record/batch-delete', method: 'post', data })
export const findSysOperationRecord = (data) => service({ url: `/system/operation-record/${data.id||data.ID}`, method: 'get' })
