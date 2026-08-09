<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { $t } from '@/locales';
import { useNaiveForm } from '@/hooks/common/form';
import { useAuthStore } from '@/store/modules/auth';
import { fetchCaptcha, fetchConfirmClientWebauth, fetchLoginConfig, fetchOAuthLoginUrl, fetchOAuthProviders } from '@/service/api/auth';

defineOptions({
  name: 'PwdLogin'
});

const authStore = useAuthStore();
const route = useRoute();
const { formRef, validate } = useNaiveForm();
const oauthProviders = ref<Api.Auth.OAuthProvider[]>([]);
const activeProvider = ref('');
const clientWebauthCompleted = ref(false);
const clientWebauthLoading = ref(false);
const captchaEnabled = ref(true);
const clientPollToken = computed(() => {
  if (route.name !== 'client-webauth') return '';
  return typeof route.query.poll_token === 'string' ? route.query.poll_token : '';
});

const model: Api.Form.LoginForm = reactive({
  username: '',
  password: '',
  code: '',
  captchaId: ''
});

const captcha: Api.Auth.Captcha = reactive({
  id: '',
  img: ''
});

const rules = computed<Record<keyof Api.Form.LoginForm, App.Global.FormRule[]>>(() => {
  return {
    username: [
      {
        required: true,
        message: $t('page.user.list.inputUsername')
      }
    ],
    password: [
      {
        required: true,
        message: $t('page.user.list.inputPassword')
      }
    ],
    code: captchaEnabled.value
      ? [{ required: true, message: $t('page.login.common.codePlaceholder') }]
      : [],
    captchaId: captchaEnabled.value
      ? [{ required: true, message: $t('page.login.common.codePlaceholder') }]
      : []
  };
});

async function handleSubmit() {
  await validate();
  if (clientPollToken.value) {
    if (clientWebauthLoading.value) return;
    clientWebauthLoading.value = true;
    try {
      const { data, error } = await fetchConfirmClientWebauth(clientPollToken.value, model);
      if (!error && data?.ok) {
        clientWebauthCompleted.value = true;
        window.setTimeout(closeClientWebauthPage, 1200);
        return;
      }
      if (error?.response?.data?.message?.includes('CaptchaError')) handleCaptcha();
    } finally {
      clientWebauthLoading.value = false;
    }
    return;
  }
  const err = await authStore.login(model);
  if (err?.response?.data?.message?.includes('CaptchaError')) {
    handleCaptcha();
  }
}

function launchRustDesk() {
  window.location.href = 'rustdesk://config/';
}

function closeClientWebauthPage() {
  window.close();
}

async function handleCaptcha() {
  if (!captchaEnabled.value) return;
  const c = await fetchCaptcha();
  captcha.id = c.data?.id || '';
  captcha.img = c.data?.img || '';
  model.captchaId = captcha.id || '';
}

async function loadOAuthProviders() {
  try {
    const { data } = await fetchOAuthProviders();
    oauthProviders.value = data || [];
  } catch {
    oauthProviders.value = [];
  }
}

async function handleOAuthLogin(provider: Api.Auth.OAuthProvider) {
  if (activeProvider.value) return;
  activeProvider.value = provider.name;
  try {
    const redirect =
      typeof route.query.redirect === 'string' && route.query.redirect.startsWith('/') && !route.query.redirect.startsWith('//') && route.query.redirect !== '/'
        ? route.query.redirect
        : '/';
    const { data, error } = await fetchOAuthLoginUrl(provider.name, redirect);
    if (!error && data?.enabled && data.url) {
      window.location.href = data.url;
      return;
    }
    window.$message?.error($t('page.login.common.providerUnavailable', { provider: provider.displayName }));
  } finally {
    activeProvider.value = '';
  }
}

function providerIcon(type: string) {
  const icons: Record<string, string> = {
    github: 'mdi:github',
    qq: 'ri:qq-fill',
    google: 'logos:google-icon',
    apple: 'ic:baseline-apple',
    oidc: 'mdi:shield-account'
  };
  return icons[type] || 'mdi:login-variant';
}

onMounted(() => {
  fetchLoginConfig().then(({ data }) => {
    captchaEnabled.value = data?.captchaEnabled !== false;
    if (captchaEnabled.value) handleCaptcha();
  });
  if (!clientPollToken.value) loadOAuthProviders();
});
</script>

<template>
  <NResult v-if="clientWebauthCompleted" status="success" :title="$t('page.login.common.loginSuccess')">
    <template #footer>
      <NSpace vertical :size="12">
        <p class="text-center text-14px text-gray-500">认证信息已发送，RustDesk 客户端会自动完成登录。</p>
        <p class="text-center text-12px text-gray-400">页面将尝试自动关闭；若浏览器阻止关闭，可直接关闭此标签页。</p>
        <NButton type="primary" round block @click="launchRustDesk">
          <template #icon><SvgIcon icon="mdi:remote-desktop" /></template>
          返回 RustDesk
        </NButton>
        <NButton round block @click="closeClientWebauthPage">关闭页面</NButton>
      </NSpace>
    </template>
  </NResult>
  <NForm v-else ref="formRef" :model="model" :rules="rules" size="medium" :show-label="false">
    <NFormItem path="username">
      <NInput v-model:value="model.username" :placeholder="$t('page.login.common.userNamePlaceholder')" />
    </NFormItem>
    <NFormItem path="password">
      <NInput
        v-model:value="model.password"
        type="password"
        show-password-on="click"
        :placeholder="$t('page.login.common.passwordPlaceholder')"
      />
    </NFormItem>
    <NFormItem v-if="captchaEnabled" path="code">
      <NInput v-model:value="model.code" :clearable="true" :placeholder="$t('page.login.common.codePlaceholder')" />
      <div class="flex-shrink-0 pl-8px">
        <img
          width="120"
          height="40"
          class="cursor-pointer lt-sm:h-28px lt-sm:w-80px"
          :src="captcha.img"
          @click="handleCaptcha"
        />
      </div>
    </NFormItem>
    <NSpace vertical :size="16">
      <div class="flex-y-center justify-between">
        <NCheckbox>{{ $t('page.login.pwdLogin.rememberMe') }}</NCheckbox>
      </div>
      <NButton
        attr-type="submit"
        type="primary"
        size="medium"
        round
        block
        :loading="clientPollToken ? clientWebauthLoading : authStore.loginLoading"
        @click="handleSubmit"
      >
        {{ $t('common.confirm') }}
      </NButton>
      <NDivider v-if="oauthProviders.length > 0" class="!mt-0 !mb-0">{{ $t('page.login.common.thirdPartyLogin') }}</NDivider>
      <div class="grid grid-cols-2 gap-8px">
        <NButton
          v-for="provider in oauthProviders"
          :key="provider.name"
          tertiary
          size="small"
          :loading="activeProvider === provider.name"
          @click="handleOAuthLogin(provider)"
        >
          <template #icon><SvgIcon :icon="providerIcon(provider.type)" /></template>
          {{ $t('page.login.common.continueWith', { provider: provider.displayName }) }}
        </NButton>
      </div>
    </NSpace>
  </NForm>
</template>

<style scoped></style>
