<script setup lang="tsx">
import { onMounted, ref } from 'vue';
import { NButton, NSpace, NTag } from 'naive-ui';
import { $t } from '@/locales';
import { useAppStore } from '@/store/modules/app';
import { request } from '@/service/request';

const appStore = useAppStore();
const loading = ref(false);
const clearing = ref(false);
const killingId = ref<number | null>(null);
const data = ref<any[]>([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);

const columns = [
  { key: 'id', title: 'ID', align: 'center' as const },
  { key: 'username', title: $t('dataMap.user.username'), align: 'center' as const },
  { key: 'rustdesk_id', title: $t('dataMap.device.rustdesk_id'), align: 'center' as const },
  { key: 'device_os', title: $t('dataMap.token.device_os'), align: 'center' as const },
  { key: 'device_name', title: $t('dataMap.token.device_name'), align: 'center' as const },
  { key: 'token_hash', title: $t('dataMap.token.token_hash'), align: 'center' as const, ellipsis: { tooltip: true } },
  {
    key: 'is_admin',
    title: $t('dataMap.token.is_admin'),
    align: 'center' as const,
    render: (row: any) => (
      <NTag type={row.is_admin ? 'warning' : 'default'} size="small">
        {row.is_admin ? $t('dataMap.token.is_admin') : $t('dataMap.user.statusLabel.normal')}
      </NTag>
    )
  },
  {
    key: 'status',
    title: $t('dataMap.token.status'),
    align: 'center' as const,
    render: (row: any) => (
      <NTag type={row.status === 1 ? 'success' : 'error'} size="small">
        {row.status === 1 ? $t('common.yesOrNo.yes') : $t('common.yesOrNo.no')}
      </NTag>
    )
  },
  { key: 'expired', title: $t('dataMap.session.expired'), align: 'center' as const },
  { key: 'created_at', title: $t('dataMap.audit.created_at'), align: 'center' as const },
  {
    key: 'actions',
    title: $t('common.action'),
    align: 'center' as const,
    render: (row: any) => (
      <NSpace size="small" justify="center">
        <NButton
          type="error"
          size="small"
          quaternary
          disabled={row.is_current}
          loading={killingId.value === row.id}
          onClick={() => confirmKill(row)}
        >
          {row.is_current ? '当前 Token' : $t('page.user.sessions.kill')}
        </NButton>
      </NSpace>
    )
  }
];

async function loadData() {
  loading.value = true;
  try {
    const { data: res, error } = await request({
      url: '/tokens/list',
      params: { current: currentPage.value, size: pageSize.value }
    });
    if (!error && res) {
      data.value = res.records || [];
      total.value = res.total || 0;
    }
  } finally {
    loading.value = false;
  }
}

async function handleKill(row: any) {
  if (row.is_current || killingId.value !== null) return;
  try {
    killingId.value = row.id;
    const { error } = await request({ url: '/tokens/kill', method: 'post', data: { ids: [row.id] } });
    if (!error) {
      window.$message?.success($t('common.deleteSuccess'));
      await loadData();
    }
  } finally {
    killingId.value = null;
  }
}

function confirmKill(row: any) {
  if (row.is_current) return;
  window.$dialog?.warning({
    title: $t('common.tip'),
    content: $t('common.confirmDelete'),
    positiveText: $t('common.confirm'),
    negativeText: $t('common.cancel'),
    onPositiveClick: () => handleKill(row)
  });
}

function handlePageChange(page: number) {
  currentPage.value = page;
  loadData();
}
function handlePageSizeChange(size: number) {
  pageSize.value = size;
  currentPage.value = 1;
  loadData();
}

async function handleClearAll() {
  try {
    clearing.value = true;
    const { data: result, error } = await request<{ cleared: number; retained: number }>({
      url: '/tokens/clear',
      method: 'post'
    });
    if (!error) {
      if ((result?.cleared || 0) > 0) {
        window.$message?.success(`已清除 ${result?.cleared} 个 Token，当前登录 Token 已保留`);
      } else {
        window.$message?.info('没有可清除的其他 Token，当前登录 Token 必须保留');
      }
      await loadData();
    }
  } finally {
    clearing.value = false;
  }
}

function confirmClearAll() {
  window.$dialog?.warning({
    title: $t('common.tip'),
    content: $t('common.confirmClear'),
    positiveText: $t('common.confirm'),
    negativeText: $t('common.cancel'),
    onPositiveClick: handleClearAll
  });
}

onMounted(() => {
  loadData();
});
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <NCard :title="$t('route.system_tokens')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <NButton type="error" size="small" :loading="clearing" @click="confirmClearAll">
          {{ $t('common.clear') }}
        </NButton>
      </template>
      <NDataTable
        :columns="columns"
        :data="data"
        size="small"
        :flex-height="!appStore.isMobile"
        :scroll-x="1200"
        :loading="loading"
        remote
        :row-key="(row: any) => row.id"
        :pagination="{
          page: currentPage,
          pageSize: pageSize,
          itemCount: total,
          showSizePicker: true,
          pageSizes: [10, 20, 50],
          onChange: handlePageChange,
          onUpdatePageSize: handlePageSizeChange
        }"
        class="sm:h-full"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
