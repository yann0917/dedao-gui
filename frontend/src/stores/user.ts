import { ref, computed, reactive } from "vue";
import { defineStore } from "pinia";
import * as services from '@backend-bindings/services/models';
import { Logout } from '@backend-bindings/app'
import { Window } from '@wailsio/runtime'
import { Local } from '../utils/storage'

export const userStore = defineStore("userStore",  {
    state:() =>{
        return {
            userList:[] as services.User[],
            user:null as services.User | null,
            // 登录态：Local.get 是普通函数调用不具备响应式，模板里直接读它
            // 会在登录后不刷新（首页仍显示“立即登录”），统一改读这个状态
            loggedIn: !!Local.get("cookies")
        }
    },
    actions: {
        async logout() {
            try {
                await Logout()
                this.user = null
                this.loggedIn = false
                Local.remove("cookies")
                Local.remove("userStore")
                Window.Reload()
            } catch (error) {
                console.error('Logout failed:', error)
            }
        },
        loginSuccess(cookie: string, user: services.User) {
            this.loggedIn = true
            Local.set("cookies", cookie)
            this.user = user
        },
        isLoggedIn(): boolean {
            return this.loggedIn || !!this.user?.nickname
        }
    },
    persist: true,
});
