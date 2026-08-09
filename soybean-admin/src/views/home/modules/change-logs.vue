<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { $t } from '@/locales';
import { fetchStat } from '@/service/api/home';
import { getBuildTime, getVersionTag } from '@/utils/version';

defineOptions({
  name: 'ChangeLogs'
});

const compatVersion = ref('latest');
const serverVersion = ref('');
const buildTime = ref('');
const appVersion = getVersionTag();
const frontendBuildTime = getBuildTime();

const currentServerVersion = computed(() => serverVersion.value || appVersion);
const clientVersion = computed(() => {
  const match = compatVersion.value.match(/client-([\w.-]+?)-server-/i);
  return match?.[1] || compatVersion.value;
});
const updatedAt = computed(() => buildTime.value || frontendBuildTime || new Date().toISOString());

onMounted(async () => {
  const stat = (await fetchStat()).data;
  if (stat?.compatVersion) {
    compatVersion.value = stat.compatVersion;
  }
  if (stat?.serverVersion) {
    serverVersion.value = stat.serverVersion;
  }
  if (stat?.buildTime) {
    buildTime.value = stat.buildTime;
  }
});
</script>

<template>
  <NCard :title="$t('page.home.changeLogs')" :bordered="false" size="small" segmented class="card-wrapper">
    <div class="compat-summary">
      <SoybeanAvatar class="compat-avatar size-48px!" />
      <div class="compat-content">
        <div class="compat-title">RustDesk API 兼容状态</div>
        <div class="compat-fields">
          <div><span>客户端兼容版本</span><strong>{{ clientVersion }}</strong></div>
          <div><span>服务端版本</span><strong>{{ currentServerVersion }}</strong></div>
          <div><span>构建时间</span><strong>{{ updatedAt }}</strong></div>
        </div>
      </div>
    </div>
  </NCard>
</template>

<style scoped>
.compat-summary {
  display: grid;
  grid-template-columns: 48px minmax(0, 1fr);
  gap: 16px;
  align-items: start;
}

.compat-avatar {
  margin-top: 2px;
}

.compat-content {
  min-width: 0;
}

.compat-title {
  margin-bottom: 12px;
  font-size: 16px;
  font-weight: 600;
}

.compat-fields {
  display: grid;
  gap: 8px;
}

.compat-fields > div {
  display: grid;
  grid-template-columns: 112px minmax(0, 1fr);
  gap: 12px;
  font-size: 13px;
  line-height: 20px;
}

.compat-fields span {
  color: var(--n-text-color-3);
}

.compat-fields strong {
  min-width: 0;
  overflow-wrap: anywhere;
  font-weight: 500;
}

@media (max-width: 480px) {
  .compat-summary {
    grid-template-columns: 36px minmax(0, 1fr);
    gap: 12px;
  }

  .compat-avatar {
    width: 36px !important;
    height: 36px !important;
  }

  .compat-fields > div {
    grid-template-columns: 1fr;
    gap: 2px;
  }
}
</style>
