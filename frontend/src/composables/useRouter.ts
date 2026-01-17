import { useRouter as useVueRouter, useRoute as useVueRoute } from 'vue-router';
import type { RouteLocationRaw } from 'vue-router';
import { ROUTES, ROUTE_NAMES, buildRoute } from '../router/routes';
import { BrowserOpenURL } from '../../wailsjs/runtime';

// 类型安全的路由 hook
export function useAppRouter() {
  const router = useVueRouter();
  const route = useVueRoute();

  const openExternalUrl = (url: string) => {
    const normalized = String(url ?? '').trim();
    if (!normalized) return;

    try {
      BrowserOpenURL(normalized);
    } catch {
      window.open(normalized, '_blank');
    }
  };

  return {
    // 原始 router 和 route
    router,
    route,

    // 类型安全的导航方法
    pushCourseList() {
      return router.push({ name: ROUTE_NAMES.COURSE });
    },

    pushCourseDetail(classId: string | number, query?: Record<string, any>) {
      return router.push(buildRoute.courseDetail(classId, query));
    },

    pushArticleDetail(enid: string, from: string, query?: Record<string, any>) {
      return router.push(buildRoute.articleDetail(enid, from, query));
    },

    pushOdobList() {
      return router.push({ name: ROUTE_NAMES.ODOB });
    },

    pushOdobDetail(id: string) {
      return router.push(buildRoute.odobDetail(id));
    },

    pushEbookList() {
      return router.push({ name: ROUTE_NAMES.EBOOK });
    },

    pushEbookComment(query?: Record<string, any>) {
      return router.push(buildRoute.ebookComment(query));
    },

    pushVideo(query?: Record<string, any>) {
      return router.push({
        name: ROUTE_NAMES.VIDEO,
        query,
      });
    },

    openExternalUrl,

    openDedaoArticle(enid: string) {
      const id = String(enid ?? '').trim();
      if (!id) return;
      openExternalUrl(`https://www.dedao.cn/course/article?id=${encodeURIComponent(id)}`);
    },

    // 通用路径跳转
    push(path: string) {
      return router.push(path);
    },

    // 通用 name 跳转（推荐）
    pushByName(name: string, params?: Record<string, any>, query?: Record<string, any>) {
      return router.push({ name, params, query });
    },

    // 个人中心相关
    pushLogin() {
      return router.push({ name: ROUTE_NAMES.LOGIN });
    },

    pushSetting() {
      return router.push({ name: ROUTE_NAMES.SETTING });
    },

    // 路由返回
    back() {
      return router.back();
    },

    // 替换当前路由（不留下历史记录）
    replace(location: RouteLocationRaw) {
      return router.replace(location);
    },
  };
}

// 快速导入所有路由常量
export { ROUTES, ROUTE_NAMES, buildRoute };
