import { ref, computed, reactive } from "vue";
import { defineStore } from "pinia";
import { services } from '../../wailsjs/go/models'



export const themeStore = defineStore("themeStore",  {
    state:() =>{
        return {
            settings:{theme:"",color:""},
        }
    },
    persist: true,
});
