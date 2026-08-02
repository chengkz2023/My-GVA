// No dedicated endpoints — sys_dictionaries table is AutoMigrate-only
export const createSysDictionary = () => Promise.resolve({ code: 0 })
export const deleteSysDictionary = () => Promise.resolve({ code: 0 })
export const updateSysDictionary = () => Promise.resolve({ code: 0 })
export const findSysDictionary = () => Promise.resolve({ code: 0 })
export const getSysDictionaryList = () => Promise.resolve({ code: 0, data: { list: [], total: 0 } })
export const exportSysDictionary = () => Promise.resolve({ code: 0 })
export const importSysDictionary = () => Promise.resolve({ code: 0 })
