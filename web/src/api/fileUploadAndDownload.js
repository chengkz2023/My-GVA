import service from '@/utils/request'
export const uploadFile = (data) => service({ url: '/v2/file/upload', method: 'post', data, headers: {'Content-Type': 'multipart/form-data'} })
export const getFileList = (data) => service({ url: '/v2/file/list', method: 'get', params: data })
export const deleteFile = (data) => service({ url: `/v2/file/${data.id||data.ID}`, method: 'delete' })
export const editFileName = (data) => service({ url: `/v2/file/${data.id||data.ID}`, method: 'put', data })
