import Vue from 'vue'
import Vuex from 'vuex'

//挂载Vuex
Vue.use(Vuex)

//创建VueX对象
const store = new Vuex.Store({
    state: {
        contact: [],
    },
    mutations: {
        changeContact(state, payload) {
            state.contact = payload;
        },
    }
})

export default store