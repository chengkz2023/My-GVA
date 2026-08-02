// JWT blacklist handled via middleware (global.BlackCache), no dedicated HTTP endpoint needed
import service from '@/utils/request'
export const jsonInBlacklist = () => Promise.resolve({ code: 0 })
