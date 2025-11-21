import request from './request'
import api, { isDesktopApp } from './adapter'

// 双色球API
export const shuangseqiuAPI = {
  // 获取列表
  async getList(params) {
    const res = await api.getShuangseqiuList(
      params.page || 1,
      params.page_size || 20,
      params.issue || ''
    )
    return res.data
  },

  // 获取最新数据
  async fetch(issue = '') {
    const res = await api.fetchShuangseqiu(issue)
    return res.data
  },

  // 批量获取历史数据
  async fetchHistory(count = 100, async = false) {
    console.log('[fetchHistory] Desktop mode:', isDesktopApp, 'Count:', count, 'Async:', async)
    if (isDesktopApp) {
      // 桌面模式：直接调用Go方法（同步）
      console.log('[fetchHistory] Using desktop mode batch fetch')
      try {
        const res = await api.fetchShuangseqiuHistory(count)
        console.log('[fetchHistory] Desktop result:', res)
        return res.data
      } catch (error) {
        console.error('[fetchHistory] Desktop error:', error)
        throw error
      }
    }
    // Web模式：支持异步
    console.log('[fetchHistory] Using web mode batch fetch')
    return request.post('/api/shuangseqiu/fetch-history', null, {
      params: { count, async }
    })
  },

  // 获取统计数据
  async getStatistics(type = 'red') {
    const res = await api.getShuangseqiuStatistics(type)
    return res.data
  },

  // 获取走势数据
  async getTrend(limit = 50) {
    const res = await api.getShuangseqiuTrend(limit)
    return res.data
  }
}

// 大乐透API
export const daletouAPI = {
  // 获取列表
  async getList(params) {
    const res = await api.getDaletouList(
      params.page || 1,
      params.page_size || 20,
      params.issue || ''
    )
    return res.data
  },

  // 获取最新数据
  async fetch(issue = '') {
    const res = await api.fetchDaletou(issue)
    return res.data
  },

  // 批量获取历史数据
  async fetchHistory(count = 100, async = false) {
    if (isDesktopApp) {
      // 桌面模式：直接调用Go方法（同步）
      const res = await api.fetchDaletouHistory(count)
      return res.data
    }
    // Web模式：支持异步
    return request.post('/api/daletou/fetch-history', null, {
      params: { count, async }
    })
  },

  // 获取统计数据
  async getStatistics(type = 'front') {
    const res = await api.getDaletouStatistics(type)
    return res.data
  },

  // 获取走势数据
  async getTrend(limit = 50) {
    const res = await api.getDaletouTrend(limit)
    return res.data
  }
}

// 任务API (仅Web模式支持)
export const taskAPI = {
  // 获取任务状态
  getTask(taskId) {
    if (isDesktopApp) {
      // 桌面模式不需要任务API
      console.warn('Desktop app does not support task API')
      return Promise.resolve({ status: 'completed' })
    }
    return request.get(`/api/task/${taskId}`)
  }
}
