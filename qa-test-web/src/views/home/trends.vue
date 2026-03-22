<template>
  <div class="page">
    <a-card title="趋势分析（规划版）" size="small">
      <a-alert
        type="success"
        show-icon
        message="建议按成熟监控平台实践，分为实时趋势、历史对比、异常波动三类。"
      />

      <a-row :gutter="[12, 12]" class="row">
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-card size="small">
            <a-statistic title="实时窗口" value="30s / 5m / 15m" />
          </a-card>
        </a-col>
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-card size="small">
            <a-statistic title="历史窗口" value="1h / 24h / 7d" />
          </a-card>
        </a-col>
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-card size="small">
            <a-statistic title="分析维度" value="设备 / 通道 / 参数" />
          </a-card>
        </a-col>
      </a-row>

      <a-divider />
      <a-space direction="vertical" :size="6">
        <a-tag color="blue">趋势图支持叠加阈值线（上限/下限）</a-tag>
        <a-tag color="geekblue">支持按设备 SN 与参数类型筛选</a-tag>
        <a-tag color="cyan">支持导出 CSV 供现场追溯</a-tag>
      </a-space>
    </a-card>
  </div>
</template>

<style scoped>
.page {
  padding: 12px;
}

.row {
  margin-top: 12px;
}

@media (max-width: 900px) {
  .page {
    padding: 8px;
  }
}
</style>
