// Config and reload endpoints — /system/config/info is the V2 equivalent for read
import service from '@/utils/request'
export const getSystemConfig = () => service({ url: '/system/config/info', method: 'get' })
export const setSystemConfig = () => Promise.resolve({ code: 0 })
export const reloadSystem = () => Promise.resolve({ code: 0 })
