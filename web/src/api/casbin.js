import service from '@/utils/request'
export const updateCasbin = (data) => service({ url: '/system/api/policies', method: 'put', data: {authorityId: data.authorityId, policies: data.casbinInfos} })
export const getPolicyPathByAuthorityId = (data) => service({ url: `/system/api/policies/${data.authorityId}`, method: 'get' })
