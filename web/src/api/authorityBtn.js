// No dedicated endpoints — authority buttons are managed through menu module
export const getAuthorityBtnApi = () => Promise.resolve({ code: 0, data: {} })
export const setAuthorityBtnApi = () => Promise.resolve({ code: 0 })
export const canRemoveAuthorityBtnApi = () => Promise.resolve({ code: 0, data: true })
