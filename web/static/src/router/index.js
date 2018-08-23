import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/views/layout/Layout'

/** note: submenu only apppear when children.length>=1
 *  detail see  https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 **/

/**
* hidden: true                   if `hidden:true` will not show in the sidebar(default is false)
* alwaysShow: true               if set true, will always show the root menu, whatever its child routes length
*                                if not set alwaysShow, only more than one route under the children
*                                it will becomes nested mode, otherwise not show the root menu
* redirect: noredirect           if `redirect:noredirect` will no redirct in the breadcrumb
* name:'router-name'             the name is used by <keep-alive> (must set!!!)
* meta : {
    roles: ['admin','editor']     will control the page roles (you can set multiple roles)
    title: 'title'               the name show in submenu and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar,
    noCache: true                if true ,the page will no be cached(default is false)
  }
**/
export const constantRouterMap = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/authredirect',
    component: () => import('@/views/login/authredirect'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/errorPage/404'),
    hidden: true
  },
  {
    path: '/401',
    component: () => import('@/views/errorPage/401'),
    hidden: true
  },
  {
    path: '',
    component: Layout,
    redirect: 'dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/dashboard/index'),
        name: 'dashboard',
        meta: { title: 'dashboard', icon: 'dashboard', noCache: true }
      }
    ]
  }
]

export default new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRouterMap
})

export const asyncRouterMap = [
  {
    path: '/taskManage',
    component: Layout,
    redirect: '/taskManage/list',
    name: 'taskManage',
    meta: {
      title: 'taskManage',
      icon: 'form'
    },
    children: [
      {
        path: 'list',
        component: () => import('@/views/task/list'),
        name: 'taskList',
        meta: { title: 'taskList', icon: 'list' }
      },
      {
        path: 'edit/:id(\\d+)',
        component: () => import('@/views/task/create'),
        name: 'editTask',
        meta: { title: 'editTask', noCache: true },
        hidden: true
      },
      {
        path: 'create',
        component: () => import('@/views/task/create'),
        name: 'createTask',
        meta: { title: 'createTask', icon: 'edit' }
      }
    ]
  },

  {
    path: '/dbManage',
    component: Layout,
    redirect: '/dbManage/list',
    name: 'dbManage',
    meta: {
      title: 'dbManage',
      icon: 'nested'
    },
    children: [
      {
        path: 'list',
        component: () => import('@/views/exportdb/list'),
        name: 'expDbList',
        meta: { title: 'expDbList', icon: 'list' }
      },
      {
        path: 'add',
        component: () => import('@/views/exportdb/add'),
        name: 'addExpDb',
        meta: { title: 'addExpDb', noCache: true },
        hidden: true
      }
    ]
  },

  { path: '*', redirect: '/404', hidden: true }
]
