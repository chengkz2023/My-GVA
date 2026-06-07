import service from '@/utils/request'

// GET /v2/system/api/list
export const getApiList = (params) => {
  return service({
    url: '/v2/system/api/list',
    method: 'get',
    params
  })
}

// POST /v2/system/api
export const createApi = (data) => {
  return service({
    url: '/v2/system/api',
    method: 'post',
    data
  })
}

// GET /v2/system/api/:id
export const getApiById = (id) => {
  return service({
    url: `/v2/system/api/${id}`,
    method: 'get'
  })
}

// PUT /v2/system/api/:id
export const updateApi = (id, data) => {
  return service({
    url: `/v2/system/api/${id}`,
    method: 'put',
    data
  })
}

// DELETE /v2/system/api/:id
export const deleteApi = (id) => {
  return service({
    url: `/v2/system/api/${id}`,
    method: 'delete'
  })
}

// POST /v2/system/api/batch-delete
export const deleteApisByIds = (data) => {
  return service({
    url: '/v2/system/api/batch-delete',
    method: 'post',
    data
  })
}

// GET /v2/system/api/all
export const getAllApis = (params) => {
  return service({
    url: '/v2/system/api/all',
    method: 'get',
    params
  })
}

// GET /v2/system/api/groups
export const getApiGroups = () => {
  return service({
    url: '/v2/system/api/groups',
    method: 'get'
  })
}

export const freshCasbin = () => {
  return service({
    url: '/v2/system/api/fresh-casbin',
    method: 'post'
  })
}

export const syncApi = () => {
  return service({
    url: '/v2/system/api/sync',
    method: 'get'
  })
}

export const ignoreApi = (data) => {
  return service({
    url: '/v2/system/api/ignore',
    method: 'post',
    data
  })
}

export const enterSyncApi = (data) => {
  return service({
    url: '/v2/system/api/batch-sync',
    method: 'post',
    data
  })
}
