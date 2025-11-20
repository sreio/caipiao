<template>
  <div class="lottery-item">
    <div class="issue-number">{{ data.issue }}</div>
    <div class="balls-container">
      <span 
        v-for="i in 6" 
        :key="'red-' + i" 
        class="ball red-ball"
      >
        {{ data[`red_ball_${i}`] }}
      </span>
      <span class="ball blue-ball">{{ data.blue_ball }}</span>
    </div>
    <div class="date-info">
      <span class="draw-date">{{ formatDate(data.draw_date) }}</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  data: {
    type: Object,
    required: true
  }
})

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.lottery-item {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 15px 20px;
  border-bottom: 1px solid #ebeef5;
  transition: background-color 0.2s;
}

.lottery-item:hover {
  background-color: #f5f7fa;
}

.issue-number {
  font-weight: 600;
  color: #303133;
  min-width: 100px;
}

.balls-container {
  display: flex;
  gap: 8px;
  flex: 1;
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

.red-ball {
  background: linear-gradient(135deg, #f56c6c, #e74c3c);
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.3);
}

.blue-ball {
  background: linear-gradient(135deg, #409eff, #3498db);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.date-info {
  min-width: 120px;
  font-size: 14px;
  color: #909399;
}
</style>
