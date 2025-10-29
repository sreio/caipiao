import axios from 'axios'
import { ElMessage } from 'element-plus'

const request = axios.create({
  baseURL: '',
  timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    const res = response.data
    
    // code = 1 表示数据已存在（警告）
    if (res.code === 1) {
      ElMessage.warning(res.msg || '数据已存在')
      // 返回数据，不中断流程
      return res.data
    }
    
    // code 为其他非0值表示错误
    if (res.code !== 0) {
      ElMessage.error(res.msg || '请求失败')
      return Promise.reject(new Error(res.msg || '请求失败'))
    }
    
    return res.data
  },
  error => {
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

export default request

