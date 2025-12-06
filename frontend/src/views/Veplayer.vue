<template>
    <div id="veplayer"></div>
</template>

<script lang="ts"  setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElTable, ElMessage } from 'element-plus'
import { GetVolcPlayAuthToken,GetVolcPlayInfo } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
let securityToken = ref()
let mediaId = ref()
let mediaVolc= reactive(new services.MediaVolc)
let playInfo = reactive(new services.VodPlayInfoResp)
let playAuthToken = ref()
let playInfoQuery = ref()

onMounted(() => {
  watch(() => {
    securityToken.value = route.query.security_token
    mediaId.value = route.query.media_id
    // console.log( route.query.media_id)
    // console.log( route.query.security_token)
  },
    () => getVolcPlayAuthToken(),
    { immediate: true },
  )

})


const getVolcPlayAuthToken = async () => {
  await GetVolcPlayAuthToken(mediaId.value, securityToken.value).then((volc) => {
    console.log(volc)
    Object.assign(mediaVolc, volc)
    playInfoQuery.value = mediaVolc?.tracks[0].formats[0].volc_key_token
    playAuthToken.value =  mediaVolc?.tracks[0].formats[0].volc_play_auth_token

  }).catch((error) => {
    ElMessage({
      message: error,
      type: 'warning'
    })
  })
}

const getVolcPlayInfo = async () => {
  console.log(playInfoQuery.value)
  await GetVolcPlayInfo(playInfoQuery.value).then((info) => {
    console.log(info)
    Object.assign(playInfo, info)
  }).catch((error) => {
    ElMessage({
      message: error,
      type: 'warning'
    })
  })
}

// const playerSdk = new VePlayer({
//   id: 'veplayer',
//   width: 800,
//   height: 500,
//   getVideoByToken: {
//     playAuthToken: playAuthToken.value, 
//     definitionMap: {
//       'original': {
//         definition: 'ori',
//         definitionTextKey: 'ORI'
//       },
//       '360p': {
//         definition: 'ld',
//         definitionTextKey: 'LD'
//       },
//       '480p': {
//         definition: 'hd',
//         definitionTextKey: 'HD'
//       },
//       '720p': {
//         definition: 'uhd',
//         definitionTextKey: 'UHD'
//       }
//     }
//   },
//   languages: {
//     'zh': {
//       ORI: '原画',
//       LD: '标清',
//       HD: '超清',
//       UHD: '蓝光'
//     }
//   }
// });

// getVolcPlayInfo();
</script>