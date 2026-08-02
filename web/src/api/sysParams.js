// No dedicated endpoints — sys_params table is AutoMigrate-only
export const createSysParams = () => Promise.resolve({ code: 0 })
export const deleteSysParams = () => Promise.resolve({ code: 0 })
export const deleteSysParamsByIds = () => Promise.resolve({ code: 0 })
export const updateSysParams = () => Promise.resolve({ code: 0 })
export const findSysParams = () => Promise.resolve({ code: 0 })
export const getSysParamsList = () => Promise.resolve({ code: 0, data: { list: [], total: 0 } })
export const getSysParam = () => Promise.resolve({ code: 0, data: {} })
