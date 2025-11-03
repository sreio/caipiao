<template>
  <el-container class="layout-container">
    <el-header class="layout-header">
      <div class="header-content">
        <h1 class="title">
          <el-icon><TrophyBase /></el-icon>
          彩票数据管理系统
        </h1>
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          @select="handleMenuSelect"
          class="header-menu"
        >
          <el-sub-menu index="shuangseqiu">
            <template #title>
              <el-icon><CirclePlus /></el-icon>
              <span>双色球</span>
            </template>
            <el-menu-item index="/shuangseqiu">
              <el-icon><Document /></el-icon>
              开奖记录
            </el-menu-item>
            <el-menu-item index="/shuangseqiu/trend">
              <el-icon><TrendCharts /></el-icon>
              走势图
            </el-menu-item>
            <el-menu-item index="/shuangseqiu/recommend">
              <el-icon><MagicStick /></el-icon>
              号码推荐
            </el-menu-item>
          </el-sub-menu>
          
          <el-sub-menu index="daletou">
            <template #title>
              <el-icon><Star /></el-icon>
              <span>大乐透</span>
            </template>
            <el-menu-item index="/daletou">
              <el-icon><Document /></el-icon>
              开奖记录
            </el-menu-item>
            <el-menu-item index="/daletou/trend">
              <el-icon><TrendCharts /></el-icon>
              走势图
            </el-menu-item>
            <el-menu-item index="/daletou/recommend">
              <el-icon><MagicStick /></el-icon>
              号码推荐
            </el-menu-item>
          </el-sub-menu>
        </el-menu>
      </div>
    </el-header>
    
    <el-main class="layout-main">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </el-main>
  </el-container>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const activeMenu = computed(() => route.path)

const handleMenuSelect = (index) => {
  router.push(index)
}
</script>

<style scoped>
.layout-container {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.layout-header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 0;
  height: 64px;
  display: flex;
  align-items: center;
}

.header-content {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: #409eff;
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
}

.header-menu {
  border-bottom: none;
  background-color: transparent;
  flex: 1;
  display: flex;
  justify-content: flex-end;
}

.header-menu .el-menu-item {
  font-size: 15px;
}

.layout-main {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  width: 100%;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

