<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { getBuildTime, getVersionTag } from '@/utils/version';
import { $t } from '@/locales';
import { useAuthStore } from '@/store/modules/auth';
import { fetchLoginCaptchaConfig, fetchSaveLoginCaptchaConfig } from '@/service/api/home';

const defaultCheckUrl = 'https://raw.githubusercontent.com/liyan-lucky/rustdesk-api-server-pro/main/VERSION';
const storageKey = 'rustdesk-api-update-check-url';
const commandStorageKey = 'rustdesk-api-update-command-template';
const updateScript = './update-rustdesk-api.sh';
const defaultCommandTemplate = updateScript;
const prepareUpdateScript =
  'docker cp rustdesk-api-server-pro:/usr/local/share/rustdesk-api-server-pro/update-container.sh ./update-rustdesk-api.sh && chmod 700 ./update-rustdesk-api.sh';
const authStore = useAuthStore();
const isAdmin = computed(() => authStore.userInfo.roles.includes('R_SUPER'));
const checkUrl = ref(localStorage.getItem(storageKey) || defaultCheckUrl);
const commandTemplate = ref(localStorage.getItem(commandStorageKey) || defaultCommandTemplate);
const runningVersion = ref(getVersionTag().replace(/^v/, ''));
const latestVersion = ref('');
const compatVersion = ref('');
const checking = ref(false);
const errorMessage = ref('');
const captchaEnabled = ref(true);
const captchaConfigLoading = ref(false);
const hasUpdate = computed(() => latestVersion.value && compareVersions(latestVersion.value, runningVersion.value) > 0);
const resolvedUpdateCommand = computed(() =>
  commandTemplate.value.replaceAll('{version}', latestVersion.value || runningVersion.value)
);
const cleanupCommand =
  'docker image prune -f --filter "label=org.opencontainers.image.source=https://github.com/liyan-lucky/rustdesk-api-server-pro" 2>/dev/null; docker image prune -f 2>/dev/null';
const latestUpdateCommand = computed(() => `${prepareUpdateScript} && ${updateScript} && ${cleanupCommand}`);
const pinnedUpdateCommand = computed(
  () =>
    `${prepareUpdateScript} && IMAGE=ghcr.io/liyan-lucky/rustdesk-api-server-pro:${latestVersion.value || runningVersion.value} EXPECTED_VERSION=${latestVersion.value || runningVersion.value} ${updateScript} && ${cleanupCommand}`
);

function normalizeVersion(value: string) {
  const match = value.trim().match(/v?(\d+\.\d+\.\d+)/);
  return match?.[1] || '';
}

function compareVersions(a: string, b: string) {
  const left = a.split('.').map(Number);
  const right = b.split('.').map(Number);
  for (let i = 0; i < 3; i += 1) {
    if ((left[i] || 0) !== (right[i] || 0)) return (left[i] || 0) - (right[i] || 0);
  }
  return 0;
}

async function loadRuntimeVersion() {
  try {
    const response = await fetch('/api/version', { cache: 'no-store' });
    const payload = await response.json();
    runningVersion.value =
      payload?.compat_target?.server?.version || normalizeVersion(payload?.version || '') || runningVersion.value;
    compatVersion.value = payload?.compat_target?.client?.version || '';
  } catch {
    /* build-injected version remains available */
  }
}

async function checkUpdate() {
  checking.value = true;
  errorMessage.value = '';
  latestVersion.value = '';
  try {
    const url = new URL(checkUrl.value, window.location.origin);
    if (!['http:', 'https:'].includes(url.protocol)) throw new Error($t('page.about.invalidUrl'));
    const response = await fetch(url.toString(), {
      cache: 'no-store',
      headers: { Accept: 'application/json, text/plain' }
    });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const text = await response.text();
    let candidate = text;
    try {
      const json = JSON.parse(text);
      candidate = json.version || json.latest_version || json.tag_name || json.server?.version || text;
    } catch {
      /* plain VERSION file */
    }
    latestVersion.value = normalizeVersion(candidate);
    if (!latestVersion.value) throw new Error($t('page.about.invalidResponse'));
    localStorage.setItem(storageKey, url.toString());
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error);
  } finally {
    checking.value = false;
  }
}

function restoreDefault() {
  checkUrl.value = defaultCheckUrl;
  localStorage.setItem(storageKey, defaultCheckUrl);
}

function saveCommandTemplate() {
  localStorage.setItem(commandStorageKey, commandTemplate.value);
  window.$message?.success($t('common.updateSuccess'));
}

function restoreDefaultCommand() {
  commandTemplate.value = defaultCommandTemplate;
  localStorage.setItem(commandStorageKey, defaultCommandTemplate);
}

async function loadLoginCaptchaConfig() {
  if (!isAdmin.value) return;
  const { data } = await fetchLoginCaptchaConfig();
  captchaEnabled.value = data?.enabled !== false;
}

async function saveLoginCaptchaConfig(enabled: boolean) {
  captchaConfigLoading.value = true;
  try {
    const { error } = await fetchSaveLoginCaptchaConfig(enabled);
    if (error) {
      captchaEnabled.value = !enabled;
      return;
    }
    window.$message?.success($t('common.updateSuccess'));
  } finally {
    captchaConfigLoading.value = false;
  }
}

async function copyUpdateCommand(value = resolvedUpdateCommand.value) {
  try {
    await navigator.clipboard.writeText(value);
  } catch {
    const node = document.createElement('textarea');
    node.value = value;
    document.body.appendChild(node);
    node.select();
    document.execCommand('copy');
    node.remove();
  }
  window.$message?.success($t('common.updateSuccess'));
}

onMounted(() => {
  loadRuntimeVersion();
  checkUpdate();
  loadLoginCaptchaConfig();
});
</script>

<template>
  <NSpace vertical size="large" class="mt-16px">
    <NCard :title="$t('page.about.versionInfo')" :bordered="false">
      <NDescriptions bordered label-placement="left" :column="1">
        <NDescriptionsItem :label="$t('page.about.runningVersion')">
          <NTag type="success">{{ runningVersion }}</NTag>
        </NDescriptionsItem>
        <NDescriptionsItem :label="$t('page.about.buildTime')">{{ getBuildTime() || '-' }}</NDescriptionsItem>
        <NDescriptionsItem :label="$t('page.about.compatVersion')">{{ compatVersion || '-' }}</NDescriptionsItem>
        <NDescriptionsItem :label="$t('page.about.latestVersion')">
          {{ latestVersion || '-' }}
          <NTag v-if="latestVersion" class="ml-8px" :type="hasUpdate ? 'warning' : 'success'">
            {{ hasUpdate ? $t('page.about.updateAvailable') : $t('page.about.upToDate') }}
          </NTag>
        </NDescriptionsItem>
      </NDescriptions>
    </NCard>
    <NCard :title="$t('page.about.updateCheck')" :bordered="false">
      <NAlert type="info" class="mb-16px">{{ $t('page.about.urlTip') }}</NAlert>
      <NSpace vertical>
        <NInput v-model:value="checkUrl" :placeholder="$t('page.about.urlPlaceholder')" />
        <NSpace>
          <NButton type="primary" :loading="checking" @click="checkUpdate">{{ $t('page.about.checkNow') }}</NButton>
          <NButton @click="restoreDefault">{{ $t('page.about.restoreDefault') }}</NButton>
        </NSpace>
        <NAlert v-if="errorMessage" type="error">{{ $t('page.about.checkFailed') }}: {{ errorMessage }}</NAlert>
      </NSpace>
    </NCard>
    <NCard v-if="isAdmin" :title="$t('page.about.updateCommand')" :bordered="false">
      <NAlert type="info" class="mb-16px">{{ $t('page.about.commandTip') }}</NAlert>
      <NSpace vertical>
        <strong>{{ $t('page.about.latestCommand') }}</strong>
        <NCode :code="latestUpdateCommand" language="shell" word-wrap />
        <NButton class="self-start" type="primary" @click="copyUpdateCommand(latestUpdateCommand)">
          {{ $t('page.about.copyCommand') }}
        </NButton>
        <strong class="mt-12px">{{ $t('page.about.pinnedCommand') }}</strong>
        <NCode :code="pinnedUpdateCommand" language="shell" word-wrap />
        <NButton class="self-start" @click="copyUpdateCommand(pinnedUpdateCommand)">
          {{ $t('page.about.copyCommand') }}
        </NButton>
        <strong class="mt-12px">{{ $t('page.about.customCommand') }}</strong>
        <NInput v-model:value="commandTemplate" type="textarea" :autosize="{ minRows: 3, maxRows: 8 }" />
        <NCode :code="resolvedUpdateCommand" language="shell" word-wrap />
        <NSpace>
          <NButton type="primary" @click="copyUpdateCommand()">{{ $t('page.about.copyCommand') }}</NButton>
          <NButton @click="saveCommandTemplate">{{ $t('common.modify') }}</NButton>
          <NButton @click="restoreDefaultCommand">{{ $t('page.about.restoreDefault') }}</NButton>
        </NSpace>
      </NSpace>
    </NCard>
    <NCard v-if="isAdmin" title="登录安全" :bordered="false">
      <div class="flex-y-center justify-between">
        <div>
          <div class="font-medium">登录验证码</div>
          <div class="mt-4px text-13px text-gray-500">控制后台、用户和客户端 WebAuth 密码登录时是否验证验证码。</div>
        </div>
        <NSwitch
          v-model:value="captchaEnabled"
          :loading="captchaConfigLoading"
          @update:value="saveLoginCaptchaConfig"
        />
      </div>
    </NCard>
  </NSpace>
</template>
