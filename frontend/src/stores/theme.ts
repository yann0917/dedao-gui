import { ref, computed, reactive } from "vue";
import { defineStore } from "pinia";
import * as services from '@backend-bindings/services/models';
import { setThemeColor } from '../utils/utils'

export const themeStore = defineStore("themeStore",  {
    state:() =>{
        return {
            // 主题模式：'light' | 'dark'
            theme: 'light' as 'light' | 'dark',
            // 主题颜色(朱砂红,与 dedao-dl webui 一致)
            color: '#d32e22',
        }
    },
    getters: {
        // 是否为暗色主题
        isDark: (state) => state.theme === 'dark',
        // 当前主题类名
        themeClass: (state) => `theme-${state.theme}`,
    },
    actions: {
        // 切换主题
        toggleTheme() {
            console.log('切换主题前:', this.theme);
            this.theme = this.theme === 'light' ? 'dark' : 'light';
            console.log('切换主题后:', this.theme);
            this.applyTheme();
        },
        // 设置主题
        setTheme(theme: 'light' | 'dark') {
            console.log('设置主题:', theme);
            this.theme = theme;
            this.applyTheme();
        },
        // 设置主题色
        setThemeColor(color: string) {
            console.log('设置主题色:', color);
            this.color = color;
            setThemeColor(color);
        },
        // 应用主题到 DOM
        applyTheme() {
            const html = document.documentElement;
            console.log('应用主题到 DOM:', this.theme, this.themeClass);
            // 移除之前的主题类
            html.classList.remove('theme-light', 'theme-dark');
            // 添加当前主题类
            html.classList.add(this.themeClass);
            // 设置 CSS 变量
            html.setAttribute('data-theme', this.theme);
            // 应用主题色
            setThemeColor(this.color);
            console.log('HTML 类名:', html.className);
        },
        // 初始化主题
        initTheme() {
            console.log('初始化主题，当前主题:', this.theme);
            // 旧默认主题色(橙色)一次性迁移到朱砂红,持久化存储会随之更新
            if (this.color === '#ff6b00' || !this.color) {
                this.color = '#d32e22';
            }
            this.applyTheme();
        }
    },
    persist: true,
});
