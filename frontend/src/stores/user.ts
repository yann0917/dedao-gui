import { ref, computed, reactive } from "vue";
import { defineStore } from "pinia";
import { services } from '../../wailsjs/go/models'

export const userStore = defineStore("userStore",  {
    state:() =>{
        return {
            userList:[] as services.User[],
            user:null as services.User | null
        }
    },
    persist: true,
});
