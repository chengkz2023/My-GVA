export const toLowerCase = (str) => {
  if (!str) {
    return ''
  }
  return str.charAt(0).toLowerCase() + str.slice(1)
}

// 驼峰转换下划线
export const toSQLLine = (str) => {
  if (!str) {
    return ''
  }
  if (str === 'ID') return 'ID'
  return str.replace(/([A-Z])/g, '_$1').toLowerCase()
}
