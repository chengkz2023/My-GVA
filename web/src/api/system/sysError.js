import service from '@/utils/request'

// SysError endpoints were removed (sys_errors table has no business logic).
// These functions are kept as no-ops to prevent 404 loops from error-handel.js.
export const createSysError = () => Promise.resolve({ code: 0 })
export const deleteSysError = () => Promise.resolve({ code: 0 })
export const deleteSysErrorByIds = () => Promise.resolve({ code: 0 })
export const updateSysError = () => Promise.resolve({ code: 0 })
export const findSysError = () => Promise.resolve({ code: 0 })
export const getSysErrorList = () => Promise.resolve({ code: 0 })
export const getSysErrorPublic = () => Promise.resolve({ code: 0 })
export const getSysErrorSolution = () => Promise.resolve({ code: 0 })
