import service from '@/utils/request'

// GET /api/dictionary/list
export const getDictionaryList = (params) => {
  return service({
    url: '/dictionary/list',
    method: 'get',
    params
  })
}

// GET /api/dictionary/types — 业务引用接口
export const getDictionaryTypes = () => {
  return service({
    url: '/dictionary/types',
    method: 'get'
  })
}

// POST /api/dictionary
export const createDictionary = (data) => {
  return service({
    url: '/dictionary',
    method: 'post',
    data
  })
}

// PUT /api/dictionary/:id
export const updateDictionary = (id, data) => {
  return service({
    url: `/dictionary/${id}`,
    method: 'put',
    data
  })
}

// DELETE /api/dictionary/:id
export const deleteDictionary = (id) => {
  return service({
    url: `/dictionary/${id}`,
    method: 'delete'
  })
}

// GET /api/dictionary/details?dictionaryId=
export const getDictionaryDetails = (params) => {
  return service({
    url: '/dictionary/details',
    method: 'get',
    params
  })
}

// POST /api/dictionary/details
export const createDictionaryDetail = (data) => {
  return service({
    url: '/dictionary/details',
    method: 'post',
    data
  })
}

// PUT /api/dictionary/details/:id
export const updateDictionaryDetail = (id, data) => {
  return service({
    url: `/dictionary/details/${id}`,
    method: 'put',
    data
  })
}

// DELETE /api/dictionary/details/:id
export const deleteDictionaryDetail = (id) => {
  return service({
    url: `/dictionary/details/${id}`,
    method: 'delete'
  })
}
