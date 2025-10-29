import request from './request'

// 双色球API
export const shuangseqiuAPI = {
  // 获取列表
  getList(params) {
    return request.get('/api/shuangseqiu/list', { params })
  },
  // 获取最新数据
  fetch(issue = '') {
    return request.post('/api/shuangseqiu/fetch', null, {
      params: { issue }
    })
  },
  // 批量获取历史数据
  fetchHistory(count = 100, async = false) {
    return request.post('/api/shuangseqiu/fetch-history', null, {
      params: { count, async }
    })
  },
  // 获取统计数据
  getStatistics(type = 'red') {
    return request.get('/api/shuangseqiu/statistics', {
      params: { type }
    })
  },
  // 获取走势数据
  getTrend(limit = 50) {
    return request.get('/api/shuangseqiu/trend', {
      params: { limit }
    })
  }
}

// 大乐透API
export const daletouAPI = {
  // 获取列表
  getList(params) {
    return request.get('/api/daletou/list', { params })
  },
  // 获取最新数据
  fetch(issue = '') {
    return request.post('/api/daletou/fetch', null, {
      params: { issue }
    })
  },
  // 批量获取历史数据
  fetchHistory(count = 100, async = false) {
    return request.post('/api/daletou/fetch-history', null, {
      params: { count, async }
    })
  },
  // 获取统计数据
  getStatistics(type = 'front') {
    return request.get('/api/daletou/statistics', {
      params: { type }
    })
  },
  // 获取走势数据
  getTrend(limit = 50) {
    return request.get('/api/daletou/trend', {
      params: { limit }
    })
  }
}

// 任务API
export const taskAPI = {
  // 获取任务状态
  getTask(taskId) {
    return request.get(`/api/task/${taskId}`)
  }
}

