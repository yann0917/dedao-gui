import {defineStore} from "pinia";

export const settingStore = defineStore("setting", {
    state: () => {
        return {
            setting: {
                downloadDir: "",
                theme: "",
                color: ""
            },
        }
    },
    getters: {
        getDownloadDir: (state) => {
            return state.setting.downloadDir
        },
        getColor:(state)=>{
            return state.setting.color
        }
    },
    persist: true,
});
