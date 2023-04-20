import {defineStore} from "pinia";

export const settingStore = defineStore("setting", {
    state: () => {
        return {
            setting: {
                downloadDir: "",
                theme: "",
                color: "",
                ffmpegDir: "",
                wkhtmltopdfDir: "",
            },
        }
    },
    getters: {
        getDownloadDir: (state) => {
            return state.setting.downloadDir
        },
        getFfmpegDirDir: (state) => {
            return state.setting.ffmpegDir
        },
        getWkDir: (state) => {
            return state.setting.wkhtmltopdfDir
        },
        getColor:(state)=>{
            return state.setting.color
        }
    },
    persist: true,
});
