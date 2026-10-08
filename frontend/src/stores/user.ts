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
                Window.Reload()
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
