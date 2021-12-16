import Vue from 'vue'
import Router from 'vue-router'
import Im from '@/components/im/Im'
import Login from '@/components/login/Login'

Vue.use(Router)

// 解决ElementUI导航栏中的vue-router在3.0版本以上重复点菜单报错问题
const originalPush = Router.prototype.push

Router.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err)
}

const baseRouters = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'login',
    component: Login
  },
  {
    path: '/im',
    name: 'im',
    component: Im
  },
]

// 需要通过后台数据来生成的组件

const createRouter = () => new Router({
  routes: baseRouters
})

const router = createRouter()

export default router
