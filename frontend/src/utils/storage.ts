/**
 * window.localStorage 浏览器永久缓存
 * @method set 设置永久缓存
 * @method get 获取永久缓存
 * @method remove 移除永久缓存
 * @method clear 移除全部永久缓存
 */

export const Local = {
    // 存储
    set(key: string, value: any) {
        localStorage.setItem(key, JSON.stringify(value))
    },
    // 取出
    get(key: string) {
        let value = localStorage.getItem(key)
        if (value && value !== 'undefined' && value !== 'null') {
            return JSON.parse(value)
        }
        return null
    },
    // 删除
    remove(key: string) {
        localStorage.removeItem(key)
    },
    clear() {
        localStorage.clear()
    }
}

/**
 * window.sessionStorage 浏览器临时缓存
 * @method set 设置临时缓存
 * @method get 获取临时缓存
 * @method remove 移除临时缓存
 * @method clear 移除全部临时缓存
 */
export const Session = {
    // 存储
    set(key: string, value: any) {
        sessionStorage.setItem(key, JSON.stringify(value))
    },
    // 取出
    get(key: string) {
        let value = sessionStorage.getItem(key)
        if (value && value !== 'undefined' && value !== 'null') {
            return JSON.parse(value)
        }
        return null
    },
    // 删除
    remove(key: string) {
        sessionStorage.removeItem(key)
    },
    clear() {
        sessionStorage.clear()
    }
}