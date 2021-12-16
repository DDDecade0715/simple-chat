// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'

import LemonIMUI from 'lemon-imui'
import 'lemon-imui/dist/index.css'


import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
 
//自己写的样式
import './style/theme.css'
import './style/character.css'

import router from './router'
import store from './store'

import VueParticles from 'vue-particles'

import 'video.js/dist/video-js.css'

// 注册element-ui
Vue.use(ElementUI)

Vue.use(LemonIMUI)

Vue.config.productionTip = false

Vue.use(VueParticles)


/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>'
})
