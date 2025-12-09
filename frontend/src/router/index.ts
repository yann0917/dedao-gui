import { createRouter, createWebHashHistory } from "vue-router";
import { Store } from "pinia";
import {userStore} from "../stores/user"
// import { ref, reactive, onMounted } from 'vue'
// import { CourseCategory } from '../../wailsjs/go/backend/App'


// const category = () => {
//     CourseCategory().then(list => {
//         console.log(list)
//         return list
//     })
// }
/**
  * @params menuType
  * -1 一般为首页 / -> /home 只显示第一个子项
  * -2 为无子菜单的菜单项 /config -> /config/person 无上下级，使用一级title
  * -3 正常菜单 /system -> /system/1 | /system/2 有上下级
  */
const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: "/",
            meta: { name: "首页", icon:"House", menuType: 2 },
            redirect: '/home',
            children: [
                {
                    path: 'home',
                    name: "home",
                    component: () => import("../views/Home.vue"),
                    meta: {
                        name: "首页", requiresAuth:false
                    },
                },{
                    path: 'category',
                    name: "category",
                    component: () => import("../views/Algo.vue"),
                    meta: {
                        name: "分类", hideMenu:true, requiresAuth:false
                    },
                },
            ],
        },
        {
            path: "/bought",
            meta: { name: "我的学习", icon:"Collection", menuType: 3 },
            redirect: '/bought/course',
            children: [
                {
                    path: 'course',
                    name: "course",
                    component: () => import("../views/Course.vue"),
                    meta: {
                        name: "课程", requiresAuth:true
                    },
                },
                {
                    path: 'course/:id',
                    name: "articleList",
                    component: () => import("../views/ArticleList.vue"),
                    meta: {
                        name: "章节列表", hideMenu:true, requiresAuth:true
                    },
                },
                {
                    path: 'article/:id',
                    name: "article",
                    component: () => import("../views/Article.vue"),
                    meta: {
                        name: "文稿", hideMenu:true, requiresAuth:true
                    },
                },
                {
                    path: 'video',
                    name: "video",
                    component: () => import("../views/Veplayer.vue"),
                    meta: {
                        name: "视频", hideMenu:true, requiresAuth:true
                    },
                },
                {
                    path: 'odob',
                    name: "odob",
                    component: () => import("../views/Odob.vue"),
                    meta: {
                        name: "听书书架", requiresAuth:true
                    },
                },
                {
                    path: 'odob/:id',
                    name: "odobArticle",
                    component: () => import("../views/Article.vue"),
                    meta: {
                        name: "文稿",hideMenu:true, requiresAuth:true
                    },
                },
                {
                    path: 'ebook',
                    name: "ebook",
                    component: () => import("../views/Ebook.vue"),
                    meta: {
                        name: "电子书架", requiresAuth:true
                    },
                },
                {
                    path: 'ebook/comment',
                    name: "ebookComment",
                    component: () => import("../views/EbookComment.vue"),
                    meta: {
                        name: "书评", hideMenu:true,requiresAuth:false
                    },
                }
            ],
        },
        {
            path: "/knowledge",
            name: "knowledge",
            meta: { name: "知识城邦", icon:"Notebook", menuType: 2 },
            redirect: "/knowledge/home",
            children: [
                {
                    path: 'home',
                    name: "knowledgeHome",
                    component: () => import("../views/Knowledge.vue"),
                    meta: {
                        name: "知识城邦", requiresAuth:false
                    },
                }
            ]
        },
        {
            path: "/",
            // name: "compass",
            meta: { name: "锦囊", icon:"Compass", menuType: 2 },
            redirect: "/compass",
            children: [
                {
                    path: 'compass',
                    name: "compass",
                    component: () => import("../views/Compass.vue"),
                    meta: {
                        name: "锦囊", requiresAuth:true
                    },
                }
            ]
        },
        {
            path: "/channel",
            meta: { name: "学习圈", icon:"CollectionTag", menuType: 3 },
            redirect: '/channel/ai',
            children: [
                {
                    path: 'ai',
                    name: "aiChannel",
                    component: () => import("../views/AIChannel.vue"),
                    meta: {
                        name: "AI学习圈", requiresAuth:true
                    },
                    props: { channelId: 1000 } // 设置默认channelId
                },
            ],
        },
        // {
        //     path: "/",
        //     // name: "about",
        //     meta: { name: "关于", icon:"InfoFilled", menuType: 2 },
        //     redirect: "/about",
        //     children: [
        //         {
        //             path: 'about',
        //             name: "about",
        //             component: () => import("../views/About.vue"),
        //             meta: {
        //                 name: "关于", requiresAuth:false
        //             },
        //         }
        //     ],
        // },
        {
            path: "/user",
            // name: "user",
            meta: { name: "个人中心", icon:"Avatar", menuType: 3 },
            redirect: "/user/profile",
            children: [
                {
                    path: 'login',
                    name: "login",
                    component: () => import("../views/Login.vue"),
                    meta: {
                        name: "登录",icon:"promotion", hideMenu:true
                    },
                },
                {
                    path: 'profile',
                    name: "profile",
                    component: () => import("../views/UserCenter.vue"),
                    meta: {
                        name: "个人简介",icon:"User", requiresAuth:true
                    },
                },
                {
                    path: 'switch',
                    name: "switch",
                    component: () => import("../views/Home.vue"),
                    meta: {
                        name: "切换账号",icon:"switchButton"
                    },
                },
                {
                    path: 'setting',
                    name: "setting",
                    component: () => import("../views/Setting.vue"),
                    meta: {
                        name: "设置",icon:"Setting", requiresAuth:false
                    },
                },
            ],
        }
    ],
});


// router.beforeEach(async (to, from, next) => {
//     // ✅ 这样做是可行的，因为路由器在安装完之后就会开始导航。
//     // Pinia 也将被安装。
//     const store = userStore()
//     if (to.name !=="login" && to.meta.requiresAuth && !store.user?.nickname) {
//         next( {path:'/user/login'})
//     } else {
//         next()
//     }
//   })

  // 全局解析守卫 它在每次导航时都会触发，但是确保在导航被确认之前，同时在所有组件内守卫和异步路由组件被解析之后，解析守卫就被正确调用。
  router.beforeResolve(async to => {
    // if (to.meta.requiresCamera) {
    //   try {
    //     await askForCameraPermission()
    //   } catch (error) {
    //     if (error instanceof NotAllowedError) {
    //       // ... 处理错误，然后取消导航
    //       return false
    //     } else {
    //       // 意料之外的错误，取消导航并把错误传给全局处理器
    //       throw error
    //     }
    //   }
    // }
  })
export default router;