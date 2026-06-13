// No V2 endpoints — SysVersion is AutoMigrate-only, version served via /system/version/info
export const deleteSysVersion = () => Promise.resolve({ code: 0 })
export const deleteSysVersionByIds = () => Promise.resolve({ code: 0 })
export const findSysVersion = () => Promise.resolve({ code: 0, data: {} })
export const getSysVersionList = () => Promise.resolve({ code: 0, data: { list: [], total: 0 } })
export const exportVersion = () => Promise.resolve({ code: 0 })
export const downloadVersionJson = () => Promise.resolve({ code: 0 })
export const importVersion = () => Promise.resolve({ code: 0 })
