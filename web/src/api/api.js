import service from '@/utils/request'

// GET /api/system/api/list
export const getApiList = (params) => {
  return service({
    url: '/system/api/list',
    method: 'get',
    params
  })
}

// POST /api/system/api
export const createApi = (data) => {
  return service({
    url: '/system/api',
    method: 'post',
    data
  })
}

// GET /api/system/api/:id
export const getApiById = (id) => {
  return service({
    url: `/system/api/${id}`,
    method: 'get'
  })
}

// PUT /api/system/api/:id
export const updateApi = (id, data) => {
  return service({
    url: `/system/api/${id}`,
    method: 'put',
    data
  })
}

// DELETE /api/system/api/:id
export const deleteApi = (id) => {
  return service({
    url: `/system/api/${id}`,
    method: 'delete'
  })
}

// POST /api/system/api/batch-delete
export const deleteApisByIds = (data) => {
  return service({
    url: '/system/api/batch-delete',
    method: 'post',
    data
  })
}

// GET /api/system/api/all
export const getAllApis = (params) => {
  return service({
    url: '/system/api/all',
    method: 'get',
    params
  })
}

// GET /api/system/api/groups
export const getApiGroups = () => {
  return service({
    url: '/system/api/groups',
    method: 'get'
  })
}

export const freshCasbin = () => {
  return service({
    url: '/system/api/fresh-casbin',
    method: 'post'
  })
}

export const syncApi = () => {
  return service({
    url: '/system/api/sync',
    method: 'get'
  })
}

export const ignoreApi = (data) => {
  return service({
    url: '/system/api/ignore',
    method: 'post',
    data
  })
}

export const enterSyncApi = (data) => {
  return service({
    url: '/system/api/batch-sync',
    method: 'post',
    data
  })
}
