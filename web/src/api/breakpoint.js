// No dedicated endpoints — breakpoint upload not migrated
export const findFile = () => Promise.resolve({ code: 0, data: {} })
export const breakpointContinue = () => Promise.resolve({ code: 0, data: {} })
export const breakpointContinueFinish = () => Promise.resolve({ code: 0, data: {} })
export const removeChunk = () => Promise.resolve({ code: 0 })
