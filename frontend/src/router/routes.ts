// 路由路径常量定义
export const ROUTES = {
  // 首页
  HOME: '/home',
  CATEGORY: '/category',

  // 我的学习
  BOUGHT: {
    BASE: '/bought',
    COURSE: '/bought/course',
    COURSE_DETAIL: (id: string | number) => `/bought/course/${id}`,
    ARTICLE: (id: string) => `/bought/article/${id}`,
    VIDEO: '/bought/video',
    ODOB: '/bought/odob',
    ODOB_DETAIL: (id: string) => `/bought/odob/${id}`,
    EBOOK: '/bought/ebook',
    EBOOK_COMMENT: '/bought/ebook/comment',
  },

  // 知识城邦
  KNOWLEDGE: '/knowledge/home',

  // 锦囊
  COMPASS: '/compass',

  // 个人中心
  USER: {
    LOGIN: '/user/login',
    PROFILE: '/user/profile',
    SWITCH: '/user/switch',
    SETTING: '/user/setting',
  },

  // 其他
  ALGO: '/algo',
} as const;

// 路由名称常量（用于命名路由跳转）
export const ROUTE_NAMES = {
  // 首页
  HOME: 'home',
  CATEGORY: 'category',

  // 我的学习
  COURSE: 'course',
  ARTICLE_LIST: 'articleList',
  ARTICLE: 'article',
  VIDEO: 'video',
  ODOB: 'odob',
  ODOB_ARTICLE: 'odobArticle',
  EBOOK: 'ebook',
  EBOOK_COMMENT: 'ebookComment',

  // 知识城邦
  KNOWLEDGE: 'knowledgeHome',

  // 锦囊
  COMPASS: 'compass',

  // 个人中心
  LOGIN: 'login',
  PROFILE: 'profile',
  SWITCH: 'switch',
  SETTING: 'setting',

  // 其他
  ALGO: 'algo',
} as const;

// 路由查询参数类型
type RouteQuery = Record<string, string | number | undefined>;

// 路由构建工具 - 使用命名路由
export const buildRoute = {
  // 通过路径跳转
  path: (path: string, query?: RouteQuery) => ({ path, query }),

  // 通过名称跳转（推荐）
  name: (name: string, params?: Record<string, any>, query?: RouteQuery) => ({
    name,
    params,
    query,
  }),

  // 课程相关
  courseDetail: (classId: string | number, query?: RouteQuery) => ({
    name: ROUTE_NAMES.ARTICLE_LIST,
    params: { id: String(classId) },
    query,
  }),

  // 文章详情
  articleDetail: (enid: string, from: string, query?: RouteQuery) => ({
    name: ROUTE_NAMES.ARTICLE,
    params: { id: enid },
    query: { from, ...query },
  }),

  // 听书详情
  odobDetail: (id: string, query?: RouteQuery) => ({
    name: ROUTE_NAMES.ODOB_ARTICLE,
    params: { id },
    query: { from: 'odob', ...query },
  }),

  // 书评列表
  ebookComment: (query?: RouteQuery) => ({
    name: ROUTE_NAMES.EBOOK_COMMENT,
    query,
  }),
};
