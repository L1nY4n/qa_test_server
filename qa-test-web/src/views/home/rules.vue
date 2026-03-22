<template>
  <div class="page">
    <a-card title="规则策略（规划版）" size="small">
      <a-alert
        type="warning"
        show-icon
        message="建议将规则拆分为采集规则、告警规则、联动规则三层，支持启停与版本化。"
      />
      <a-divider />
      <a-space direction="vertical" :size="6">
        <a-tag color="purple">采集规则：采样周期、数据清洗、字段映射</a-tag>
        <a-tag color="magenta">告警规则：阈值、连续触发次数、恢复条件</a-tag>
        <a-tag color="volcano">联动规则：告警分派、通知渠道、自动恢复动作</a-tag>
      </a-space>
      <a-divider />
      <a-empty description="规则编辑器可在下一迭代接入" />
    </a-card>
  </div>
</template>

<style scoped>
.page {
  padding: 12px;
}

@media (max-width: 900px) {
  .page {
    padding: 8px;
  }
}
</style>
