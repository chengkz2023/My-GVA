import service from '@/utils/request'
export const uploadFile = (data) => service({ url: '/file/upload', method: 'post', data, headers: {'Content-Type': 'multipart/form-data'} })
export const getFileList = (data) => service({ url: '/file/list', method: 'get', params: data })
export const deleteFile = (data) => service({ url: `/file/${data.id||data.ID}`, method: 'delete' })
export const editFileName = (data) => service({ url: `/file/${data.id||data.ID}`, method: 'put', data })
