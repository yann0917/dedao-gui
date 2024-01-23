import {settingStore} from "../stores/setting";
import {OpenDirectoryDialog} from "../../wailsjs/go/backend/App";
const store = settingStore()
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
  // const nprogress = document.querySelector("#nprogress .bar");
  // console.log(nprogress)
  el.style.setProperty("--el-color-primary", color);
  body?.style.setProperty("--van-primary-color", color);
  // nprogress.style.setProperty("background", color);
  // 此行判断是否为白天/暗夜模式，可根据自身业务调整代码
  let mixColor = store.setting.theme === "light" ? "#ffffff" : "#000000";
  // 此行判断是否为白天/暗夜模式，可根据自身业务调整代码
  for (let i = 1; i < 10; i++) {
    el.style.setProperty(`--el-color-primary-light-${i}`, colourBlend(color, mixColor, i / 10));
    el.style.setProperty(`--el-color-primary-dark-${i}`, colourBlend(color, mixColor, i / 10));
  }
  el.style.setProperty(`--el-color-primary-dark-2`, colourBlend(color, mixColor, 0.2));
  if (mixColor == "#ffffff") {
    store.setting.theme="light"
  } else {
    store.setting.theme="dark"
  }
  store.setting.color = color
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