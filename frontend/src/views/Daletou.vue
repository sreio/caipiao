<template>
  <div class="daletou-container">
    <el-card class="search-card">
      <div class="search-bar">
        <el-input
          v-model="searchIssue"
          placeholder="请输入期号搜索"
          clearable
          style="width: 300px"
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="handleSearch">
          <el-icon><Search /></el-icon>
          搜索
        </el-button>
        <el-button type="success" @click="handleFetch('')">
          <el-icon><Download /></el-icon>
          获取最新数据
        </el-button>
        <el-button type="warning" @click="showFetchDialog = true">
          <el-icon><Calendar /></el-icon>
          指定期数获取
        </el-button>
        <el-button @click="showStatistics = true">
          <el-icon><DataAnalysis /></el-icon>
          数据统计
        </el-button>
        <el-button type="info" @click="showHistoryDialog = true">
          <el-icon><Refresh /></el-icon>
          批量获取历史数据
        </el-button>
      </div>
    </el-card>

    <el-card class="table-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">大乐透开奖记录</span>
          <span class="total-count">共 {{ total }} 条记录</span>
        </div>
      </template>
      
      <el-table
        v-loading="loading"
        :data="tableData"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="issue" label="期号" min-width="120" fixed />
        <el-table-column label="前区" min-width="380">
          <template #default="{ row }">
            <div class="balls-container">
              <span class="ball front-ball">{{ row.front_ball_1 }}</span>
              <span class="ball front-ball">{{ row.front_ball_2 }}</span>
              <span class="ball front-ball">{{ row.front_ball_3 }}</span>
              <span class="ball front-ball">{{ row.front_ball_4 }}</span>
              <span class="ball front-ball">{{ row.front_ball_5 }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="后区" min-width="180">
          <template #default="{ row }">
            <div class="balls-container">
              <span class="ball back-ball">{{ row.back_ball_1 }}</span>
              <span class="ball back-ball">{{ row.back_ball_2 }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="开奖日期" min-width="180">
          <template #default="{ row }">
            {{ formatDate(row.draw_date) }}
          </template>
        </el-table-column>
        <el-table-column label="创建时间" min-width="200">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 指定期数获取对话框 -->
    <el-dialog
      v-model="showFetchDialog"
      title="指定期数获取数据"
      width="400px"
    >
      <el-form>
        <el-form-item label="期号">
          <el-input v-model="fetchIssue" placeholder="请输入期号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFetchDialog = false">取消</el-button>
        <el-button type="primary" @click="handleFetchByIssue">确定</el-button>
      </template>
    </el-dialog>

    <!-- 批量获取历史数据对话框 -->
    <el-dialog
      v-model="showHistoryDialog"
      title="批量获取历史数据"
      width="500px"
    >
      <el-form label-width="100px">
        <el-form-item label="获取期数">
          <el-input-number
            v-model="historyCount"
            :min="10"
            :max="500"
            :step="10"
            placeholder="请输入期数"
          />
          <div style="color: #909399; font-size: 12px; margin-top: 5px;">
            默认100期，最多500期
          </div>
        </el-form-item>
      </el-form>
      
      <div v-if="historyResult" class="history-result">
        <el-descriptions title="获取结果" :column="2" border>
          <el-descriptions-item label="总计">{{ historyResult.total }}</el-descriptions-item>
          <el-descriptions-item label="成功">
            <el-tag type="success">{{ historyResult.success }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="跳过">
            <el-tag type="warning">{{ historyResult.skipped }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="失败">
            <el-tag type="danger">{{ historyResult.failed }}</el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </div>
      
      <template #footer>
        <el-button @click="showHistoryDialog = false">取消</el-button>
        <el-button 
          type="primary" 
          :loading="historyLoading"
          @click="handleFetchHistory"
        >
          {{ historyLoading ? '获取中...' : '开始获取' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 统计数据对话框 -->
    <el-dialog
      v-model="showStatistics"
      title="大乐透统计数据"
      width="900px"
    >
      <el-tabs v-model="statsType" @tab-change="loadStatistics">
        <el-tab-pane label="前区统计" name="front">
          <div class="statistics-container">
            <div
              v-for="stat in statistics"
              :key="stat.number"
              class="stat-item"
            >
              <div class="stat-number front-ball">{{ stat.number }}</div>
              <div class="stat-count">出现 {{ stat.count }} 次</div>
              <el-progress
                :percentage="getPercentage(stat.count)"
                :color="'#f56c6c'"
              />
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="后区统计" name="back">
          <div class="statistics-container">
            <div
              v-for="stat in statistics"
              :key="stat.number"
              class="stat-item"
            >
              <div class="stat-number back-ball">{{ stat.number }}</div>
              <div class="stat-count">出现 {{ stat.count }} 次</div>
              <el-progress
                :percentage="getPercentage(stat.count)"
                :color="'#409eff'"
              />
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { daletouAPI, taskAPI } from '@/api/lottery'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const searchIssue = ref('')
const showFetchDialog = ref(false)
const fetchIssue = ref('')
const showStatistics = ref(false)
const statsType = ref('front')
const statistics = ref([])
const showHistoryDialog = ref(false)
const historyCount = ref(100)
const historyLoading = ref(false)
const historyResult = ref(null)

const pagination = reactive({
  page: 1,
  page_size: 20
})

// 加载列表数据
const loadData = async () => {
  loading.value = true
  try {
    const res = await daletouAPI.getList({
      page: pagination.page,
      page_size: pagination.page_size,
      issue: searchIssue.value
    })
    tableData.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadData()
}

// 获取最新数据
const handleFetch = async (issue) => {
  loading.value = true
  try {
    await daletouAPI.fetch(issue)
    // 成功或已存在的消息会由拦截器自动显示
    loadData()
  } catch (error) {
    console.error('获取数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 指定期数获取
const handleFetchByIssue = () => {
  if (!fetchIssue.value) {
    ElMessage.warning('请输入期号')
    return
  }
  handleFetch(fetchIssue.value)
  showFetchDialog.value = false
  fetchIssue.value = ''
}

// 批量获取历史数据
const handleFetchHistory = async () => {
  historyLoading.value = true
  historyResult.value = null
  
  try {
    // 判断是否需要异步处理（超过100期）
    const useAsync = historyCount.value > 100
    const result = await daletouAPI.fetchHistory(historyCount.value, useAsync)
    
    if (result.task_id) {
      // 异步任务，开始轮询
      ElMessage.info('任务已创建，正在后台获取数据...')
      await pollTask(result.task_id)
    } else {
      // 同步完成
      historyResult.value = result
      
      if (result.success > 0) {
        ElMessage.success(`成功获取 ${result.success} 期数据`)
        loadData()
      } else if (result.skipped === result.total) {
        ElMessage.info('所有数据都已存在')
      } else {
        ElMessage.warning('批量获取完成，部分数据失败')
      }
    }
  } catch (error) {
    console.error('批量获取失败:', error)
  } finally {
    historyLoading.value = false
  }
}

// 轮询任务状态
const pollTask = async (taskId) => {
  let pollCount = 0
  const maxPolls = 300 // 最多轮询5分钟（每2秒一次）
  
  const pollInterval = setInterval(async () => {
    pollCount++
    
    try {
      const task = await taskAPI.getTask(taskId)
      
      // 更新进度
      historyResult.value = {
        total: task.total || 0,
        success: task.success || 0,
        skipped: task.skipped || 0,
        failed: task.failed || 0
      }
      
      if (task.status === 'completed') {
        clearInterval(pollInterval)
        historyLoading.value = false
        ElMessage.success(`批量获取完成！成功 ${task.success} 期`)
        loadData()
      } else if (task.status === 'failed') {
        clearInterval(pollInterval)
        historyLoading.value = false
        ElMessage.error('批量获取失败：' + task.message)
      } else if (pollCount >= maxPolls) {
        clearInterval(pollInterval)
        historyLoading.value = false
        ElMessage.warning('任务执行超时，请稍后查看列表')
      }
    } catch (error) {
      console.error('查询任务失败:', error)
      if (pollCount >= 3) {
        clearInterval(pollInterval)
        historyLoading.value = false
        ElMessage.error('无法获取任务状态')
      }
    }
  }, 2000) // 每2秒查询一次
}

// 加载统计数据
const loadStatistics = async () => {
  try {
    const res = await daletouAPI.getStatistics(statsType.value)
    statistics.value = res || []
    // 按出现次数排序
    statistics.value.sort((a, b) => b.count - a.count)
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 计算百分比
const getPercentage = (count) => {
  if (!statistics.value.length) return 0
  const maxCount = Math.max(...statistics.value.map(s => s.count))
  return Math.round((count / maxCount) * 100)
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString('zh-CN')
}

// 格式化日期时间
const formatDateTime = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleString('zh-CN')
}

// 分页
const handleSizeChange = () => {
  loadData()
}

const handleCurrentChange = () => {
  loadData()
}

// 监听统计对话框打开，自动加载数据
watch(showStatistics, (newVal) => {
  if (newVal) {
    loadStatistics()
  }
})

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.daletou-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.search-bar {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.table-card {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.total-count {
  font-size: 14px;
  color: #909399;
}

.balls-container {
  display: flex;
  gap: 8px;
}

.ball {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  font-size: 16px;
  font-weight: 600;
  color: #fff;
}

.front-ball {
  background: linear-gradient(135deg, #f56c6c, #e74c3c);
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.3);
}

.back-ball {
  background: linear-gradient(135deg, #409eff, #3498db);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.statistics-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
  max-height: 600px;
  overflow-y: auto;
  padding: 10px;
}

.stat-item {
  padding: 16px;
  background-color: #f5f7fa;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.stat-number {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 8px;
}

.stat-count {
  font-size: 14px;
  color: #606266;
  margin-bottom: 8px;
}

.history-result {
  margin-top: 20px;
}
</style>

