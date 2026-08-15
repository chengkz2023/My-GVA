// 独立日期格式化函数（不污染 Date.prototype）
// 月(M)、日(d)、小时(h)、分(m)、秒(s)、季度(q) 可用 1-2 个占位符，
// 年(y)可用 1-4 个占位符，毫秒(S)只能用 1 个占位符(1-3 位数字)
// formatDate(new Date(), "yyyy-MM-dd hh:mm:ss.S") => 2006-07-02 08:09:04.423
function formatDate(date, fmt) {
  if (!date || isNaN(date.getTime())) {
    return ''
  }
  const o = {
    'M+': date.getMonth() + 1, // 月份
    'd+': date.getDate(), // 日
    'h+': date.getHours(), // 小时
    'm+': date.getMinutes(), // 分
    's+': date.getSeconds(), // 秒
    'q+': Math.floor((date.getMonth() + 3) / 3), // 季度
    S: date.getMilliseconds() // 毫秒
  }
  const reg = /(y+)/
  if (reg.test(fmt)) {
    const t = reg.exec(fmt)[1]
    fmt = fmt.replace(t, (date.getFullYear() + '').substring(4 - t.length))
  }
  for (const k in o) {
    const regx = new RegExp('(' + k + ')')
    if (regx.test(fmt)) {
      const t = regx.exec(fmt)[1]
      fmt = fmt.replace(t, t.length === 1 ? o[k] : ('00' + o[k]).substring(('' + o[k]).length))
    }
  }
  return fmt
}

export function formatTimeToStr(times, pattern) {
  const d = new Date(times)
  if (isNaN(d.getTime())) {
    return ''
  }
  if (pattern) {
    return formatDate(d, pattern)
  }
  return formatDate(d, 'yyyy-MM-dd hh:mm:ss')
}
