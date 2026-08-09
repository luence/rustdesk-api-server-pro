<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useAppStore } from '@/store/modules/app';
import { useThemeStore } from '@/store/modules/theme';
import { useWebBackground } from '@/hooks/common/web-background';

defineOptions({ name: 'ClientOAuthResultPage' });

const route = useRoute();
const appStore = useAppStore();
const themeStore = useThemeStore();
const { backgroundStyle } = useWebBackground();

const success = computed(() => route.query.status === 'success');
const errorCode = computed(() => (typeof route.query.error_code === 'string' ? route.query.error_code : ''));

function returnToClient() {
  window.location.href = 'rustdesk://config/';
}

function closePage() {
  window.close();
}

onMounted(() => {
  if (success.value) window.setTimeout(closePage, 1200);
});
</script>

<template>
  <div class="relative size-full flex-center overflow-hidden p-24px" :style="backgroundStyle">
    <NCard :bordered="false" class="relative z-4 w-auto rd-12px text-center">
      <div class="w-320px max-w-[calc(100vw-48px)]">
        <header class="flex-y-center justify-between">
          <SystemLogo class="text-48px text-primary lt-sm:text-36px" />
          <h3 class="text-22px text-primary font-500 lt-sm:text-18px">{{ $t('system.title') }}</h3>
          <div class="i-flex-col">
            <ThemeSchemaSwitch
              :theme-schema="themeStore.themeScheme"
              :show-tooltip="false"
              class="text-20px"
              @switch="themeStore.toggleThemeScheme"
            />
            <LangSwitch
              :lang="appStore.locale"
              :lang-options="appStore.localeOptions"
              :show-tooltip="false"
              @change-lang="appStore.changeLocale"
            />
          </div>
        </header>
        <main class="pt-24px">
          <div
            class="mx-auto size-48px flex-center rounded-full text-26px font-700"
            :class="success ? 'bg-green-500/15 text-green-500' : 'bg-red-500/15 text-red-500'"
          >
            {{ success ? '✓' : '!' }}
          </div>
          <h1 class="mb-0 mt-16px text-20px font-600">
            {{ success ? '第三方登录成功' : '第三方登录失败' }}
          </h1>
          <p class="mb-0 mt-12px text-14px text-gray-500 leading-24px">
            {{ success ? '认证已完成，客户端正在自动获取登录结果。' : '登录失败，请返回客户端后重试。' }}
          </p>
          <p v-if="!success && errorCode" class="mb-0 mt-8px text-13px text-error">错误码：{{ errorCode }}</p>
          <div class="mt-24px flex flex-col gap-10px">
            <NButton v-if="success" type="primary" block @click="returnToClient">返回 RustDesk 客户端</NButton>
            <NButton block @click="closePage">关闭页面</NButton>
          </div>
          <p class="mb-0 mt-16px text-12px text-gray-400 leading-20px">
            若浏览器阻止自动关闭，可使用上方按钮或直接关闭此标签页。
          </p>
        </main>
      </div>
    </NCard>
  </div>
</template>
