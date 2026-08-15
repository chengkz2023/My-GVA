export const fmtTitle = (title, now) => {
  if (!title || typeof title !== 'string') {
    return ''
  }
  const reg = /\$\{(.+?)\}/
  const reg_g = /\$\{(.+?)\}/g
  const result = title.match(reg_g)
  if (result) {
    result.forEach((item) => {
      const key = item.match(reg)[1]
      const value = now?.params?.[key] ?? now?.query?.[key] ?? ''
      title = title.replace(item, value)
    })
  }
  return title
}
