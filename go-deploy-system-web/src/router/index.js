import Vue from 'vue';
import Router from 'vue-router';

Vue.use(Router);
const originalPush = Router.prototype.push
Router.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err)
}

export default new Router({
    mode: "history", // 去掉URL中的#号
    routes: [
        {
            path: '/',
            redirect: '/user'
        },
        {
            path: '/',
            component: () => import(/* webpackChunkName: "home" */ '../components/common/Home.vue'),
            meta: { title: '自述文件' },
            children: [
                {
                    path: '/releaseCode',
                    component: () => import(/* webpackChunkName: "releaseCode" */ '../components/page/releaseCode.vue'),
                    meta: { title: '发布代码' }
                },
                {
                    path: '/deploymentLogs',
                    component: () => import(/* webpackChunkName: "deploymentLogs" */ '../components/page/deploymentLogs.vue'),
                    meta: { title: '发布日志' }
                },
                {
                    path: '/projectConfig',
                    component: () => import(/* webpackChunkName: "customerList" */ '../components/page/projectConfig.vue'),
                    meta: { title: '项目配置' }
                },
                {
                    path: '/department',
                    component: () => import(/* webpackChunkName: "customerList" */ '../components/page/department.vue'),
                    meta: { title: '部门管理' }
                },
                {
                    path: '/service',
                    component: () => import(/* webpackChunkName: "customerList" */ '../components/page/service.vue'),
                    meta: { title: '服务器管理' }
                },
                {
                    path: '/engineroom',
                    component: () => import(/* webpackChunkName: "customerList" */ '../components/page/engineroom.vue'),
                    meta: { title: '机房管理' }
                },
                {
                    path: '/user',
                    component: () => import(/* webpackChunkName: "customerList" */ '../components/page/user.vue'),
                    meta: { title: '用户管理' }
                },
                {
                    path: '/404',
                    component: () => import(/* webpackChunkName: "404" */ '../components/page/404.vue'),
                    meta: { title: '404' }
                },
                {
                    path: '/403',
                    component: () => import(/* webpackChunkName: "403" */ '../components/page/403.vue'),
                    meta: {title: '403' }
                }
            ]
        },
        {
            path: '/login',
            component: () => import(/* webpackChunkName: "login" */ '../components/page/Login.vue'),
            meta: { title: '文件发布系统' }
        },
        {
            path: '*',
            redirect: '/404'
        }
    ]
});
