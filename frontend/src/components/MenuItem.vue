<template>
	<template v-if="!props.menu.meta?.hideMenu">
		<el-sub-menu v-if="props.menu.meta?.menuType === 3" :index="menuPath">
			<template #title>
				<el-icon>
					<component :is="props.menu.meta?.icon"></component>
				</el-icon>
				<span>{{ props.menu.meta?.name }}</span>
			</template>
			<template v-for="children in props.menu.children" :key="children.path">
				<menu-item v-if="!children.meta?.hideMenu" :menu="children" :path="`${menuPath}/${children.path}`" />
			</template>
		</el-sub-menu>
		<router-link v-else :to="menuPath">
			<el-menu-item :index="menuPath">
				<el-icon v-if="props.menu.meta?.icon">
					<component :is="props.menu.meta.icon"></component>
				</el-icon>
				<template #title>
					<span>{{
					props.menu.meta?.menuType === 1
					? props.menu.children![0].meta?.name
					: props.menu.meta?.name
					}}</span>
				</template>
			</el-menu-item>
		</router-link>
	</template>
</template>

<script lang="ts" setup>
import { computed, PropType } from 'vue'
import { RouteRecordRaw } from 'vue-router'
const props = defineProps({
	menu: {
		type: Object as PropType<RouteRecordRaw>,
		required: true
	},
	path: {
		type: String,
		default: ''
	}
})
const menuPath = computed(() => {
	if ([1, 2].includes(props.menu.meta?.menuType as number)) {
		return (
			(props.path === '/' ? props.path : props.path + '/') +
			props.menu.children![0].path
		)
	}
	return props.path
})
</script>
<style scoped>
	a{
		text-decoration: none;
	}
</style>