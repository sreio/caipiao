// API适配层 - 自动检测并使用正确的API调用方式
import axios from 'axios'

// 检测是否在Wails环境中
const isWailsApp = typeof window !== 'undefined' && window.runtime && typeof window.runtime.EventsOn === 'function'

// Wails API导入（懒加载）
let WailsAPI = null
let wailsAPIPromise = null

// 懒加载Wails API
function loadWailsAPI() {
    if (!wailsAPIPromise) {
        wailsAPIPromise = import('../../wailsjs/go/main/App.js')
            .then(module => {
                WailsAPI = module
                return module
            })
            .catch(error => {
                console.warn('Failed to load Wails API:', error)
                return null
            })
    }
    return wailsAPIPromise
}

// 获取Wails API（确保已加载）
async function getWailsAPI() {
    if (isWailsApp) {
        if (!WailsAPI) {
            await loadWailsAPI()
        }
        return WailsAPI
    }
    return null
}

// API适配器
export const api = {
    // 双色球列表
    async getShuangseqiuList(page, pageSize, issue = '') {
        const wails = await getWailsAPI()
        if (wails) {
            const result = await wails.GetShuangseqiuList(page, pageSize, issue)
            return { data: result }
        } return axios.get('/api/shuangseqiu/list', {
            params: { page, page_size: pageSize, issue }
        })
    },

    // 获取双色球数据
    async fetchShuangseqiu(issue = '') {
        const wails = await getWailsAPI()
        if (wails) {
            const result = await wails.FetchShuangseqiu(issue)
            return { data: result }
        }
        return axios.post('/api/shuangseqiu/fetch', null, {
            params: { issue }
        })
    },

    // 双色球统计
    async getShuangseqiuStatistics(ballType) {
        const wails = await getWailsAPI()
        if (wails) {
            const result = await wails.GetShuangseqiuStatistics(ballType)
            return { data: result }
        }
        return axios.get('/api/shuangseqiu/statistics', {
            params: { ball_type: ballType }
        })
    },

    // 双色球走势
    async getShuangseqiuTrend(limit = 30) {
        const wails = await getWailsAPI()
        if (wails) {
            const result = await wails.GetShuangseqiuTrend(limit)
            return { data: result }
        }
        return axios.get('/api/shuangseqiu/trend', {
            params: { limit }
        })
    },

    // 双色球推荐
    async getShuangseqiuRecommendation(count = 5) {
        const wails = await getWailsAPI()
        if (wails) {
            const result = await wails.GetShuangseqiuRecommendation(count)
            return { data: result }
        }
        return axios.get('/api/shuangseqiu/recommend', {
            params: { count }
        })
    },

    // 大乐透列表
    async getDaletouList(page, pageSize, issue = '') {
        const wails = await getWailsAPI()
        if (wails) {
            const result = await wails.GetDaletouList(page, pageSize, issue)
            return { data: result }
        }
        return axios.get('/api/daletou/list', {
            params: { page, page_size: pageSize, issue }
        })
    },

    // 获取大乐透数据
    async fetchDaletou(issue = '') {
        const wails = await getWailsAPI()
        if (wails) {
            const result = await wails.FetchDaletou(issue)
            return { data: result }
        }
        return axios.post('/api/daletou/fetch', null, {
            params: { issue }
        })
    },

    // 大乐透统计
    async getDaletouStatistics(ballType) {
        const wails = await getWailsAPI()
        if (wails) {
            const result = await wails.GetDaletouStatistics(ballType)
            return { data: result }
        }
        return axios.get('/api/daletou/statistics', {
            params: { ball_type: ballType }
        })
    },

    // 大乐透走势
    async getDaletouTrend(limit = 30) {
        const wails = await getWailsAPI()
        if (wails) {
            const result = await wails.GetDaletouTrend(limit)
            return { data: result }
        }
        return axios.get('/api/daletou/trend', {
            params: { limit }
        })
    },

    // 大乐透推荐
    async getDaletouRecommendation(count = 5) {
        const wails = await getWailsAPI()
        if (wails) {
            const result = await wails.GetDaletouRecommendation(count)
            return { data: result }
        }
        return axios.get('/api/daletou/recommend', {
            params: { count }
        })
    },

    // 双色球批量获取历史数据
    async fetchShuangseqiuHistory(count = 100) {
        const wails = await getWailsAPI()
        if (wails) {
            const result = await wails.FetchShuangseqiuHistory(count)
            return { data: result }
        }
        // Web模式：使用HTTP请求（同步模式）
        return axios.post('/api/shuangseqiu/fetch-history', null, {
            params: { count, async: false }
        })
    },

    // 大乐透批量获取历史数据
    async fetchDaletouHistory(count = 100) {
        const wails = await getWailsAPI()
        if (wails) {
            const result = await wails.FetchDaletouHistory(count)
            return { data: result }
        }
        // Web模式：使用HTTP请求（同步模式）
        return axios.post('/api/daletou/fetch-history', null, {
            params: { count, async: false }
        })
    }
}

// 导出环境检测
export const isDesktopApp = isWailsApp

export default api
