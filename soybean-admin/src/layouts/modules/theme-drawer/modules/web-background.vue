<script setup lang="ts">
import { computed } from 'vue';
import type { UploadFileInfo } from 'naive-ui';
import { clearUploadedBackground, saveUploadedBackground, useWebBackground } from '@/hooks/common/web-background';
import SettingItem from '../components/setting-item.vue';

defineOptions({ name: 'WebBackground' });
const { backgroundMode, uploadedBackground, globalBackgroundEnabled, setBackgroundMode, setGlobalBackgroundEnabled } = useWebBackground();
const mode = computed({ get: () => backgroundMode.value, set: setBackgroundMode });
const globalEnabled = computed({ get: () => globalBackgroundEnabled.value, set: setGlobalBackgroundEnabled });
const modeOptions = computed(() => [
  { label: '默认星空图片', value: 'fixed' },
  { label: 'Bing 每日聚焦', value: 'bing' },
  { label: '上传的图片', value: 'upload', disabled: !uploadedBackground.value }
]);

async function handleUpload(options: { file: UploadFileInfo }) {
  if (!options.file.file) return;
  try {
    await saveUploadedBackground(options.file.file);
    window.$message?.success('背景图片已保存到当前浏览器');
  } catch (error) {
    window.$message?.error(error instanceof Error ? error.message : 'ERR-1019: 背景图片保存失败');
  }
}
</script>

<template>
  <NDivider>Web 背景</NDivider>
  <div class="flex-col-stretch gap-12px">
    <SettingItem label="背景来源"><NSelect v-model:value="mode" :options="modeOptions" size="small" class="w-150px" /></SettingItem>
    <SettingItem label="全部页面显示"><NSwitch v-model:value="globalEnabled" /></SettingItem>
    <NUpload :show-file-list="false" accept="image/*" :custom-request="handleUpload"><NButton block secondary>上传背景图片</NButton></NUpload>
    <NButton v-if="uploadedBackground" block quaternary type="error" @click="clearUploadedBackground">清除上传图片</NButton>
    <NText depth="3" class="text-12px">登录页始终显示所选背景；全部页面开关只影响登录后的 Web 页面。上传图片仅保存在当前浏览器。</NText>
  </div>
</template>
