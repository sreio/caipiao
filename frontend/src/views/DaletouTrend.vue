<template>
  <div class="trend-container">
    <!-- 控制栏 -->
    <el-card class="control-card">
      <div class="control-bar">
        <el-select v-model="limit" @change="loadTrend" style="width: 150px">
          <el-option label="最近30期" :value="30" />
          <el-option label="最近50期" :value="50" />
          <el-option label="最近100期" :value="100" />
        </el-select>
        
        <el-button type="primary" @click="loadTrend" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新数据
        </el-button>
      </div>
    </el-card>

    <!-- 走势图 -->
    <el-card class="trend-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span class="card-title">大乐透号码走势图</span>
          <span class="period-info">共 {{ tableData.length }} 期</span>
        </div>
      </template>
      
      <div class="trend-table-wrapper">
        <table class="trend-table">
          <thead>
            <tr>
              <th class="period-col">期号</th>
              <th class="date-col">开奖日期</th>
              <th colspan="35" class="number-group-header front-header">前区</th>
              <th colspan="12" class="number-group-header back-header">后区</th>
            </tr>
            <tr>
              <th class="period-col"></th>
              <th class="date-col"></th>
              <!-- 前区号码 1-35 -->
              <th v-for="num in 35" :key="`front-${num}`" class="number-header front-bg">
                {{ String(num).padStart(2, '0') }}
              </th>
              <!-- 后区号码 1-12 -->
              <th v-for="num in 12" :key="`back-${num}`" class="number-header back-bg">
                {{ String(num).padStart(2, '0') }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, rowIndex) in tableData" :key="row.issue" class="trend-row">
              <td class="period-cell">{{ row.issue }}</td>
              <td class="date-cell">{{ formatDate(row.draw_date) }}</td>
              
              <!-- 前区走势 -->
              <td v-for="num in 35" :key="`front-${row.issue}-${num}`" class="trend-cell">
                <div class="cell-content">
                  <span v-if="row.front_balls.includes(num)" class="ball-mark front-mark">
                    {{ String(num).padStart(2, '0') }}
                  </span>
                  <span v-else class="empty-mark">{{ num }}</span>
                  <!-- 连接线 -->
                  <svg v-if="rowIndex < tableData.length - 1 && hasConnection(rowIndex, num, 'front')" 
                       class="connect-line" viewBox="0 0 2 40">
                    <line x1="1" y1="0" x2="1" y2="40" stroke="#f56c6c" stroke-width="2"/>
                  </svg>
                </div>
              </td>
              
              <!-- 后区走势 -->
              <td v-for="num in 12" :key="`back-${row.issue}-${num}`" class="trend-cell">
                <div class="cell-content">
                  <span v-if="row.back_balls.includes(num)" class="ball-mark back-mark">
                    {{ String(num).padStart(2, '0') }}
                  </span>
                  <span v-else class="empty-mark">{{ num }}</span>
                  <!-- 连接线 -->
                  <svg v-if="rowIndex < tableData.length - 1 && hasConnection(rowIndex, num, 'back')" 
                       class="connect-line" viewBox="0 0 2 40">
                    <line x1="1" y1="0" x2="1" y2="40" stroke="#409eff" stroke-width="2"/>
                  </svg>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 前区频率统计 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">前区频率统计（Top 10）</span>
            </div>
          </template>
          <div ref="frontFreqChartRef" style="width: 100%; height: 400px"></div>
          <el-table :data="topFrontFreq" style="margin-top: 20px">
            <el-table-column prop="number" label="号码" min-width="100">
              <template #default="{ row }">
                <span class="ball front-ball">{{ row.number }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="count" label="出现次数" min-width="120" />
            <el-table-column label="出现率" min-width="120">
              <template #default="{ row }">
                {{ ((row.count / (limit * 5)) * 100).toFixed(1) }}%
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- 后区频率统计 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">后区频率统计（Top 10）</span>
            </div>
          </template>
          <div ref="backFreqChartRef" style="width: 100%; height: 400px"></div>
          <el-table :data="topBackFreq" style="margin-top: 20px">
            <el-table-column prop="number" label="号码" min-width="100">
              <template #default="{ row }">
                <span class="ball back-ball">{{ row.number }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="count" label="出现次数" min-width="120" />
            <el-table-column label="出现率" min-width="120">
              <template #default="{ row }">
                {{ ((row.count / (limit * 2)) * 100).toFixed(1) }}%
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 遗漏值分析 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 前区遗漏 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">前区遗漏值（Top 10）</span>
              <el-tag type="warning" size="small">提醒：遗漏值较大的号码可能即将出现</el-tag>
            </div>
          </template>
          <el-alert
            v-if="topFrontMissing.length > 0 && topFrontMissing[0].missing > 20"
            title="长期遗漏提醒"
            :description="`号码 ${topFrontMissing[0].number} 已连续 ${topFrontMissing[0].missing} 期未出现！`"
            type="warning"
            show-icon
            style="margin-bottom: 20px"
          />
          <el-table :data="topFrontMissing">
            <el-table-column prop="number" label="号码" min-width="100">
              <template #default="{ row }">
                <span class="ball front-ball">{{ row.number }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="missing" label="遗漏期数" min-width="140">
              <template #default="{ row }">
                <el-tag :type="getMissingTagType(row.missing)">
                  {{ row.missing }} 期
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" min-width="120">
              <template #default="{ row }">
                <span v-if="row.missing === 0" style="color: #67c23a">●刚出现</span>
                <span v-else-if="row.missing < 10" style="color: #409eff">●正常</span>
                <span v-else-if="row.missing < 20" style="color: #e6a23c">●较长</span>
                <span v-else style="color: #f56c6c">●长期未出</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- 后区遗漏 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">后区遗漏值（Top 10）</span>
              <el-tag type="warning" size="small">提醒：遗漏值较大的号码可能即将出现</el-tag>
            </div>
          </template>
          <el-alert
            v-if="topBackMissing.length > 0 && topBackMissing[0].missing > 10"
            title="长期遗漏提醒"
            :description="`号码 ${topBackMissing[0].number} 已连续 ${topBackMissing[0].missing} 期未出现！`"
            type="warning"
            show-icon
            style="margin-bottom: 20px"
          />
          <el-table :data="topBackMissing">
            <el-table-column prop="number" label="号码" min-width="100">
              <template #default="{ row }">
                <span class="ball back-ball">{{ row.number }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="missing" label="遗漏期数" min-width="140">
              <template #default="{ row }">
                <el-tag :type="getMissingTagType(row.missing)">
                  {{ row.missing }} 期
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" min-width="120">
              <template #default="{ row }">
                <span v-if="row.missing === 0" style="color: #67c23a">●刚出现</span>
                <span v-else-if="row.missing < 5" style="color: #409eff">●正常</span>
                <span v-else-if="row.missing < 10" style="color: #e6a23c">●较长</span>
                <span v-else style="color: #f56c6c">●长期未出</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { daletouAPI } from '@/api/lottery'
import * as echarts from 'echarts'

const loading = ref(false)
const limit = ref(50)
const tableData = ref([])
const trendData = ref(null)

// 图表实例
const frontFreqChartRef = ref()
const backFreqChartRef = ref()
let frontFreqChart = null
let backFreqChart = null

// 计算属性
const topFrontFreq = computed(() => {
  if (!trendData.value?.front_freq) return []
  return Object.entries(trendData.value.front_freq)
    .map(([number, count]) => ({ number: parseInt(number), count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 10)
})

const topBackFreq = computed(() => {
  if (!trendData.value?.back_freq) return []
  return Object.entries(trendData.value.back_freq)
    .map(([number, count]) => ({ number: parseInt(number), count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 10)
})

const topFrontMissing = computed(() => {
  if (!trendData.value?.front_missing) return []
  return Object.entries(trendData.value.front_missing)
    .map(([number, missing]) => ({ number: parseInt(number), missing }))
    .sort((a, b) => b.missing - a.missing)
    .slice(0, 10)
})

const topBackMissing = computed(() => {
  if (!trendData.value?.back_missing) return []
  return Object.entries(trendData.value.back_missing)
    .map(([number, missing]) => ({ number: parseInt(number), missing }))
    .sort((a, b) => b.missing - a.missing)
    .slice(0, 10)
})

// 获取遗漏标签类型
const getMissingTagType = (missing) => {
  if (missing === 0) return 'success'
  if (missing < 10) return 'info'
  if (missing < 20) return 'warning'
  return 'danger'
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}-${date.getDate()}`
}

// 判断是否有连接线
const hasConnection = (rowIndex, num, type) => {
  if (rowIndex >= tableData.value.length - 1) return false
  
  const currentRow = tableData.value[rowIndex]
  const nextRow = tableData.value[rowIndex + 1]
  
  if (type === 'front') {
    return currentRow.front_balls.includes(num) && nextRow.front_balls.includes(num)
  } else {
    return currentRow.back_balls.includes(num) && nextRow.back_balls.includes(num)
  }
}

// 加载走势数据
const loadTrend = async () => {
  loading.value = true
  try {
    const data = await daletouAPI.getTrend(limit.value)
    trendData.value = data
    
    // 处理表格数据
    tableData.value = data.issues.map((issue, index) => ({
      issue: issue,
      draw_date: data.draw_dates ? data.draw_dates[index] : '',
      front_balls: data.front_balls[index] || [],
      back_balls: data.back_balls[index] || []
    }))
    
    // 渲染图表
    renderFreqCharts()
  } catch (error) {
    console.error('加载走势数据失败:', error)
    ElMessage.error('加载走势数据失败')
  } finally {
    loading.value = false
  }
}

// 渲染频率图表
const renderFreqCharts = () => {
  // 前区频率图
  if (!frontFreqChart) {
    frontFreqChart = echarts.init(frontFreqChartRef.value)
  }

  const frontFreqData = topFrontFreq.value
  const frontFreqOption = {
    title: {
      text: '前区出现频率',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    xAxis: {
      type: 'category',
      data: frontFreqData.map(item => item.number)
    },
    yAxis: {
      type: 'value',
      name: '出现次数'
    },
    series: [
      {
        data: frontFreqData.map(item => ({
          value: item.count,
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#f56c6c' },
              { offset: 1, color: '#e74c3c' }
            ])
          }
        })),
        type: 'bar',
        barWidth: '60%'
      }
    ]
  }

  frontFreqChart.setOption(frontFreqOption)

  // 后区频率图
  if (!backFreqChart) {
    backFreqChart = echarts.init(backFreqChartRef.value)
  }

  const backFreqData = topBackFreq.value
  const backFreqOption = {
    title: {
      text: '后区出现频率',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    xAxis: {
      type: 'category',
      data: backFreqData.map(item => item.number)
    },
    yAxis: {
      type: 'value',
      name: '出现次数'
    },
    series: [
      {
        data: backFreqData.map(item => ({
          value: item.count,
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#409eff' },
              { offset: 1, color: '#3498db' }
            ])
          }
        })),
        type: 'bar',
        barWidth: '60%'
      }
    ]
  }

  backFreqChart.setOption(backFreqOption)
}

// 窗口大小变化时重新渲染
const handleResize = () => {
  frontFreqChart?.resize()
  backFreqChart?.resize()
}

onMounted(() => {
  loadTrend()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  frontFreqChart?.dispose()
  backFreqChart?.dispose()
})
</script>

<style scoped>
.trend-container {
  padding: 20px;
}

.control-card {
  margin-bottom: 20px;
}

.control-bar {
  display: flex;
  gap: 12px;
  align-items: center;
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

.period-info {
  font-size: 14px;
  color: #909399;
}

/* 走势表格 */
.trend-card {
  margin-bottom: 20px;
}

.trend-table-wrapper {
  overflow-x: auto;
  max-height: 800px;
  overflow-y: auto;
}

.trend-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
  min-width: 1400px;
}

.trend-table thead {
  position: sticky;
  top: 0;
  z-index: 10;
  background: #fff;
}

.number-group-header {
  background: #f5f7fa;
  font-weight: 600;
  padding: 8px 4px;
  text-align: center;
  border: 1px solid #e4e7ed;
  font-size: 14px;
}

.front-header {
  background: linear-gradient(to bottom, #fef0f0, #fde2e2);
  color: #f56c6c;
}

.back-header {
  background: linear-gradient(to bottom, #ecf5ff, #d9ecff);
  color: #409eff;
}

.period-col,
.date-col,
.number-header {
  padding: 6px 4px;
  text-align: center;
  border: 1px solid #e4e7ed;
  background: #fff;
  font-weight: 600;
}

.period-col {
  min-width: 80px;
  position: sticky;
  left: 0;
  z-index: 11;
  background: #fff;
}

.date-col {
  min-width: 60px;
  position: sticky;
  left: 80px;
  z-index: 11;
  background: #fff;
}

.number-header {
  width: 28px;
  min-width: 28px;
  font-size: 11px;
}

.front-bg {
  background: #fef0f0;
  color: #f56c6c;
}

.back-bg {
  background: #ecf5ff;
  color: #409eff;
}

.trend-row {
  border-bottom: 1px solid #e4e7ed;
}

.trend-row:hover {
  background: #f5f7fa;
}

.period-cell {
  padding: 4px;
  text-align: center;
  font-weight: 600;
  font-size: 12px;
  border: 1px solid #e4e7ed;
  position: sticky;
  left: 0;
  background: #fff;
  z-index: 1;
}

.trend-row:hover .period-cell {
  background: #f5f7fa;
}

.date-cell {
  padding: 4px;
  text-align: center;
  font-size: 11px;
  color: #909399;
  border: 1px solid #e4e7ed;
  position: sticky;
  left: 80px;
  background: #fff;
  z-index: 1;
}

.trend-row:hover .date-cell {
  background: #f5f7fa;
}

.trend-cell {
  padding: 0;
  text-align: center;
  border: 1px solid #e4e7ed;
  height: 40px;
  position: relative;
  width: 28px;
  min-width: 28px;
}

.cell-content {
  position: relative;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ball-mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  font-size: 11px;
  font-weight: 600;
  color: #fff;
  position: relative;
  z-index: 2;
}

.front-mark {
  background: linear-gradient(135deg, #f56c6c, #e74c3c);
  box-shadow: 0 1px 4px rgba(245, 108, 108, 0.4);
}

.back-mark {
  background: linear-gradient(135deg, #409eff, #3498db);
  box-shadow: 0 1px 4px rgba(64, 158, 255, 0.4);
}

.empty-mark {
  font-size: 10px;
  color: #dcdfe6;
  font-weight: normal;
}

.connect-line {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translateX(-50%);
  width: 2px;
  height: 40px;
  pointer-events: none;
  z-index: 1;
}

/* 频率统计表格样式 */
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
</style>

