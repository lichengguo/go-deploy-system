import Vue from 'vue'
import App from './App'
import router from './router'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css'; // 默认主题
import  axios from 'axios';
// import  qs from 'qs'
// import md5 from 'js-md5';

// Vue.prototype.$qs = qs;
// Vue.prototype.$md5 = md5;
Vue.prototype.$axios = axios;
Vue.config.productionTip = false;

Vue.use(ElementUI, {
    size: 'small'
});

// axios拦截器 发送请求的时候，携带token
axios.interceptors.request.use(config => {
    config.headers.Authorization = `Bearer ${localStorage.getItem('token')}`
    return config
})

// 响应拦截
axios.interceptors.response.use(response=> {
      if (response.data.code == 200) {
        return response.data.result;
      } else if (response.data.code == 3001) {
        localStorage.removeItem('token');
        router.push('/login')
      }
      return response
    }, err=> {
      console.log(err)
      return Promise.reject(err);
  })
  

// 使用钩子函数对路由进行权限跳转
router.beforeEach((to, from, next) => {
  document.title = `${to.meta.title}`;
  let role = localStorage.getItem('token');
  let path = to.path

  // // 简单的判断IE10及以下不进入富文本编辑器，该组件不兼容
  // if (navigator.userAgent.indexOf('MSIE') > -1 && to.path === '/editor') {
  //     Vue.prototype.$alert('vue-quill-editor组件不兼容IE10及以下浏览器，请使用更高版本的浏览器查看', '浏览器不兼容通知', {
  //         confirmButtonText: '确定'
  //     });
  // }
  if(!role && path !== '/login'){
      next("/login");
  }else if(role && path === '/login'){
      next("/");
  }else{
      next();
  }
});


/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
