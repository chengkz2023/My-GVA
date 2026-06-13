import service from '@/utils/request'
export const getAuthorityList = (data) => service({ url: '/system/role/tree', method: 'get' })
export const createAuthority = (data) => service({ url: '/system/role', method: 'post', data })
export const deleteAuthority = (data) => service({ url: `/system/role/${data.authorityId}`, method: 'delete' })
export const updateAuthority = (data) => service({ url: `/system/role/${data.authorityId}`, method: 'put', data })
export const copyAuthority = (data) => service({ url: '/system/role/copy', method: 'post', data })
export const setDataAuthority = (data) => service({ url: '/system/role/data-authority', method: 'post', data })
export const getAuthorityInfo = (data) => service({ url: `/system/role/${data.authorityId}`, method: 'get' })
