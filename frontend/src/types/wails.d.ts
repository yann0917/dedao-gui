// Wails runtime types
declare module '../../wailsjs/go/backend/App' {
  export function CourseList(arg1: string, arg2: string, arg3: number, arg4: number): Promise<any>
  export function CourseCategory(): Promise<any>
  export function SetDir(arg: string[]): Promise<void>
  export function GetHomeInitialState(): Promise<any>
  export function SearchHot(): Promise<any>
  export function SunflowerLabelList(arg: number): Promise<any>
  export function SunflowerLabelContent(enid: string, type: number, page: number, pageSize: number): Promise<any>
  export function SunflowerResourceList(): Promise<any>
  export function UserInfo(): Promise<any>
  export function EbookUserInfo(): Promise<any>
  export function OdobUserInfo(): Promise<any>
  export function Logout(): Promise<any>
}

declare module '../../wailsjs/go/models' {
  export namespace services {
    export class CourseList {}
    export class CourseInfo {}
    export class ClassInfo {}
    export class ArticleIntro {}
    export class HomeInitState {}
    export class HomeData {}
    export class SunflowerLabelList {}
    export class SunflowerContent {}
    export class SunflowerResourceList {}
    export class User {}
    export class EbookVIPInfo {}
    export class OdobVip {}
    export class OdobUser {}
  }
}

declare module '../../wailsjs/runtime' {
  export function BrowserOpenURL(url: string): void
  export function WindowReloadApp(): void
}
