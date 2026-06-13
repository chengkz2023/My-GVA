import packageInfo from '../../package.json'

const greenText = (text) => `\x1b[32m${text}\x1b[0m`

export const config = {
  appName: 'BoyKing Admin',
  showViteLogo: true,
  keepAliveTabs: false,
  useStaticMenu: false,
  useMockLogin: false,
  logs: []
}

export const viteLogo = (env) => {
  if (config.showViteLogo) {
    console.log(greenText(`> ${config.appName} v${packageInfo.version}`))
    console.log(greenText(`> 本地开发: http://127.0.0.1:${env.VITE_CLI_PORT || '5173'}`))
    console.log('')
  }
}

export default config
