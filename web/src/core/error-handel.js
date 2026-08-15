// 全局未捕获错误处理：后端 sysError 接口已移除，这里改为本地控制台日志
function reportError(errorInfo) {
  const info = `${errorInfo.message}\nStack: ${errorInfo.stack}${errorInfo.component ? `\nComponent: ${errorInfo.component.name || 'Unknown'}` : ''}${errorInfo.vueInfo ? `\nVue Info: ${errorInfo.vueInfo}` : ''}${errorInfo.source ? `\nSource: ${errorInfo.source}:${errorInfo.lineno}:${errorInfo.colno}` : ''}`
  console.error(`[frontend error] type=${errorInfo.type}`, info)
}

window.addEventListener('unhandledrejection', (event) => {
  reportError({
    type: '前端',
    message: `错误信息: ${event.reason}`,
    stack: `调用栈: ${event.reason?.stack || '没有调用栈信息'}`
  })
})
