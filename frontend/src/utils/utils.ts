import {settingStore} from "../stores/setting";
import {OpenDirectoryDialog} from "../../wailsjs/go/backend/App";
export const secondToHour = function (msd:number) {
    let second = 0; // 秒
    let minute = 0; // 分
    let hour = 0; // 小时
    if (msd > 0) {
      minute = Number(Math.floor(msd / 60).toFixed(0));
      second = msd % 60;
      if (minute > 60) {
        hour = Number(Math.floor(minute / 60).toFixed(0));
        minute = minute % 60;
      }
    }
    var result = ""
    if (second > 0) {
      result += second >=10 ? second + "秒": '0'+second + "秒";
    }
    if (minute > 0) {
      result = "" + minute + "分" + result;
    }
    if (hour > 0) {
      result = "" + hour + "小时" + result;
    }
    return result;
  }

  export const  timestampToTime = (timestamp:number)=> {
    var date = new Date(timestamp * 1000);//时间戳为10位需*1000，时间戳为13位的话不需乘1000
    var Y = date.getFullYear() + '-';
    var M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1) + '-';
    var D = date.getDate() + ' ';
    var h = date.getHours() + ':';
    var m = date.getMinutes() + ':';
    var s = date.getSeconds();
    return Y+M+D+h+m+s;
}

export const  timestampToDate = (timestamp:number)=> {
  var date = new Date(timestamp * 1000);//时间戳为10位需*1000，时间戳为13位的话不需乘1000
  var Y = date.getFullYear() + '-';
  var M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1) + '-';
  var D = date.getDate() + ' ';
  return Y+M+D;
}


export function setThemeColor(color:any) {
  const el = document.documentElement;
  const body = document.querySelector("body");
  
  // 设置 Element Plus 主题色
  el.style.setProperty("--el-color-primary", color);
  body?.style.setProperty("--van-primary-color", color);
  
  // 设置自定义主题色变量
  el.style.setProperty("--accent-color", color);
  
  // 生成主题色的悬停效果色
  const hoverColor = adjustBrightness(color, 20);
  el.style.setProperty("--accent-hover", hoverColor);
  
  // 判断当前主题模式
  const isDark = el.classList.contains('theme-dark');
  let mixColor = isDark ? "#000000" : "#ffffff";
  
  // 生成主题色的各种变体
  for (let i = 1; i < 10; i++) {
    el.style.setProperty(`--el-color-primary-light-${i}`, colourBlend(color, mixColor, i / 10));
    el.style.setProperty(`--el-color-primary-dark-${i}`, colourBlend(color, mixColor, i / 10));
  }
  el.style.setProperty(`--el-color-primary-dark-2`, colourBlend(color, mixColor, 0.2));
  
  // 更新设置存储 - 在函数内部获取 store 实例
  try {
    const store = settingStore();
    store.setting.color = color;
    store.setting.theme = isDark ? "dark" : "light";
  } catch (error) {
    console.warn('无法更新设置存储，可能 Pinia 尚未初始化:', error);
  }
  
  // 生成暗色模式下的主题色变体
  if (isDark) {
    generateDarkThemeColors(color);
  }
}

export function setFontFamily(font: string) {
  const el = document.documentElement;
  let fontFamilyValue = "";
  
  switch (font) {
    case 'jetbrains':
      fontFamilyValue = "'JetBrainsMono', 'PingFang SC', sans-serif";
      break;
    case 'wenkai':
      fontFamilyValue = "'LXGW WenKai Lite', 'PingFang SC', sans-serif";
      break;
    default:
      el.style.removeProperty('--font-family-base');
      el.style.removeProperty('--el-font-family');
      return;
  }
  
  el.style.setProperty('--font-family-base', fontFamilyValue);
  el.style.setProperty('--el-font-family', fontFamilyValue);
}

// 调整颜色亮度的辅助函数
function adjustBrightness(color: string, percent: number): string {
  const num = parseInt(color.replace("#", ""), 16);
  const amt = Math.round(2.55 * percent);
  const R = (num >> 16) + amt;
  const G = (num >> 8 & 0x00FF) + amt;
  const B = (num & 0x0000FF) + amt;
  return "#" + (0x1000000 + (R < 255 ? R < 1 ? 0 : R : 255) * 0x10000 +
    (G < 255 ? G < 1 ? 0 : G : 255) * 0x100 +
    (B < 255 ? B < 1 ? 0 : B : 255)).toString(16).slice(1);
}

// 生成暗色模式下的主题色变体
function generateDarkThemeColors(primaryColor: string) {
  const el = document.documentElement;
  
  // 生成暗色模式下的主题色变体
  const darkAccent = adjustBrightness(primaryColor, -10);
  const darkAccentHover = adjustBrightness(primaryColor, 10);
  
  el.style.setProperty("--accent-color-dark", darkAccent);
  el.style.setProperty("--accent-hover-dark", darkAccentHover);
  
  // 生成暗色模式下的背景色变体
  const darkBg = adjustBrightness(primaryColor, -95);
  const darkCardBg = adjustBrightness(primaryColor, -90);
  const darkHoverBg = adjustBrightness(primaryColor, -85);
  
  el.style.setProperty("--dark-bg-color", darkBg);
  el.style.setProperty("--dark-card-bg", darkCardBg);
  el.style.setProperty("--dark-hover-bg", darkHoverBg);
}

export function colourBlend(c1:any, c2:any, ratio:any) {
  ratio = Math.max(Math.min(Number(ratio), 1), 0)
  let r1 = parseInt(c1.substring(1, 3), 16)
  let g1 = parseInt(c1.substring(3, 5), 16)
  let b1 = parseInt(c1.substring(5, 7), 16)
  let r2 = parseInt(c2.substring(1, 3), 16)
  let g2 = parseInt(c2.substring(3, 5), 16)
  let b2 = parseInt(c2.substring(5, 7), 16)

  let r = Math.round(r1 * (1 - ratio) + r2 * ratio)
  let g = Math.round(g1 * (1 - ratio) + g2 * ratio)
  let b = Math.round(b1 * (1 - ratio) + b2 * ratio)

  let rs = ('0' + (r || 0).toString(16)).slice(-2)
  let gs = ('0' + (g || 0).toString(16)).slice(-2)
  let bs = ('0' + (b || 0).toString(16)).slice(-2)

  return '#' + rs + gs + bs
}

export function loadJs(src:string) {
  return new Promise((resolve,reject)=>{
    let script = document.createElement('script');
    script.type = "text/javascript";
    script.src= src;
    document.body.appendChild(script);
      
    script.onload = ()=>{
      resolve(src);
    }
    script.onerror = ()=>{
      reject();
    }
  })
}

export function unloadJs(src:string) {
  return new Promise((resolve,reject)=>{
    let script = document.createElement('script');
    script.type = "text/javascript";
    script.src= src;
    document.body.removeChild(script);
    script.onload = ()=>{
      resolve(src);
    }
    script.onerror = ()=>{
      reject();
    }
  })
}

export const openDialogDir = async (title: string) => {
  await OpenDirectoryDialog(title).then((result)=>{
    console.log(result)
  }).catch((error)=>{
    console.log(error)
  })
}