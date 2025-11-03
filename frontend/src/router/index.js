import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/components/Layout.vue'

const routes = [
  {
    path: '/',
    component: Layout,
    redirect: '/shuangseqiu',
    children: [
      {
        path: '/shuangseqiu',
        name: 'Shuangseqiu',
        component: () => import('@/views/Shuangseqiu.vue'),
        meta: { title: '双色球' }
      },
      {
        path: '/shuangseqiu/trend',
        name: 'ShuangseqiuTrend',
        component: () => import('@/views/ShuangseqiuTrend.vue'),
        meta: { title: '双色球走势' }
      },
      {
        path: '/shuangseqiu/recommend',
        name: 'ShuangseqiuRecommend',
        component: () => import('@/views/ShuangseqiuRecommend.vue'),
        meta: { title: '双色球号码推荐' }
      },
      {
        path: '/daletou',
        name: 'Daletou',
        component: () => import('@/views/Daletou.vue'),
        meta: { title: '大乐透' }
      },
      {
        path: '/daletou/trend',
        name: 'DaletouTrend',
        component: () => import('@/views/DaletouTrend.vue'),
        meta: { title: '大乐透走势' }
      },
      {
        path: '/daletou/recommend',
        name: 'DaletouRecommend',
        component: () => import('@/views/DaletouRecommend.vue'),
        meta: { title: '大乐透号码推荐' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

