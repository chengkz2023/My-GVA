// JWT blacklist handled via middleware (global.BlackCache), no V2 HTTP endpoint needed
import service from '@/utils/request'
export const jsonInBlacklist = () => Promise.resolve({ code: 0 })
