import { ref, computed, reactive } from "vue";
import { defineStore } from "pinia";
import { services } from '../../wailsjs/go/models'
import { Logout } from '../../wailsjs/go/backend/App'
import { WindowReloadApp } from '../../wailsjs/runtime'
import { Local } from '../utils/storage'

export const userStore = defineStore("userStore",  {
    state:() =>{
        return {
            userList:[] as services.User[],
            user:null as services.User | null
        }
    },
    actions: {
        async logout() {
            try {
                await Logout()
                this.user = null
                Local.remove("cookies")
                Local.remove("userStore")
                WindowReloadApp()
            } catch (error) {
                console.error('Logout failed:', error)
            }
        },
        isLoggedIn(): boolean {
            return !!this.user?.nickname
        }
    },
    persist: true,
});
