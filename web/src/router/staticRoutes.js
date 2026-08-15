/**
 * 前端静态菜单（脚手架默认）
 * 当 config.useStaticMenu 为 true 时使用，无需后端 getMenu 接口
 */
import { asyncRouterHandle } from '@/utils/asyncRouter'

const staticMenus = [
  {
    path: 'admin',
    name: 'superAdmin',
    meta: { title: '超级管理员', icon: 'user' },
    component: 'view/superAdmin/index.vue',
    children: [
      { path: 'authority', name: 'authority', meta: { title: '角色管理', icon: 'avatar' }, component: 'view/superAdmin/authority/authority.vue' },
      { path: 'menu', name: 'menu', meta: { title: '菜单管理', icon: 'tickets', keepAlive: true }, component: 'view/superAdmin/menu/menu.vue' },
      { path: 'api', name: 'api', meta: { title: 'API管理', icon: 'platform', keepAlive: true }, component: 'view/superAdmin/api/api.vue' },
      { path: 'user', name: 'user', meta: { title: '用户管理', icon: 'coordinate' }, component: 'view/superAdmin/user/user.vue' },
      { path: 'operation', name: 'operation', meta: { title: '操作历史', icon: 'pie-chart' }, component: 'view/superAdmin/operation/sysOperationRecord.vue' }
    ]
  }
]

export function getStaticRouter() {
  const baseRouter = [
    {
      path: '/layout',
      name: 'layout',
      component: 'view/layout/index.vue',
      meta: { title: '底层layout' },
      children: [...staticMenus]
    }
  ]
  asyncRouterHandle(baseRouter)
  return baseRouter
}

export { staticMenus }
