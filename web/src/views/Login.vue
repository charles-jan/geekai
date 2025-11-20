<template>
  <div class="flex-center loginPage">
    <div class="left">
      <div class="login-box">
        <AccountTop>
          <template #default>
            <div class="wechatLog flex-center" v-if="wechatLoginURL !== ''">
              <a :href="wechatLoginURL" @click="setRoute(router.currentRoute.value.path)"> <i class="iconfont icon-wechat"></i>使用微信登录 </a>
            </div>
          </template>
        </AccountTop>

        <div class="input-form">
          <el-form ref="ruleFormRef" :model="ruleForm" :rules="currentRules">
            <!-- 手机号+验证码登录 -->
            <template v-if="enableVerify">
              <el-form-item label="" prop="mobile">
                <div class="form-title">手机号码</div>
                <el-input v-model="ruleForm.mobile" size="large" placeholder="请输入手机号码" maxlength="11" @keyup="handleKeyup" />
              </el-form-item>
              <el-form-item label="" prop="code">
                <div class="form-title">验证码</div>
                <div class="flex w100">
                  <el-input v-model="ruleForm.code" size="large" placeholder="请输入验证码" maxlength="6" class="code-input" @keyup="handleKeyup" />
                  <send-msg size="large" :receiver="ruleForm.mobile" type="mobile" />
                </div>
              </el-form-item>
            </template>
            
            <!-- 用户名+密码登录 -->
            <template v-else>
              <el-form-item label="" prop="username">
                <div class="form-title">账号</div>
                <el-input v-model="ruleForm.username" size="large" placeholder="请输入账号" @keyup="handleKeyup" />
              </el-form-item>
              <el-form-item label="" prop="password">
                <div class="flex-between w100">
                  <div class="form-title">密码</div>
                  <div class="form-forget text-color-primary" @click="router.push('/resetpassword')">忘记密码？</div>
                </div>
                <el-input size="large" v-model="ruleForm.password" placeholder="请输入密码" show-password autocomplete="off" @keyup="handleKeyup" />
              </el-form-item>
            </template>
            
            <el-form-item>
              <el-button class="login-btn" size="large" type="primary" @click="login">登录</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
    <account-bg />

    <captcha v-if="enableVerify" @success="doLogin" ref="captchaRef" />
  </div>
</template>

<script setup>
import { onMounted, ref, reactive, computed } from "vue";
import { httpGet, httpPost } from "@/utils/http";
import { useRouter } from "vue-router";
import AccountBg from "@/components/AccountBg.vue";
import { isMobile } from "@/utils/libs";
import { checkSession, getLicenseInfo, getSystemInfo } from "@/store/cache";
import { setUserToken } from "@/store/session";
import { showMessageError } from "@/utils/dialog";
import { setRoute } from "@/store/system";
import { useSharedStore } from "@/store/sharedata";

import AccountTop from "@/components/AccountTop.vue";
import Captcha from "@/components/Captcha.vue";
import SendMsg from "@/components/SendMsg.vue";

const router = useRouter();
const title = ref("Geek-AI");

const logo = ref("");
const licenseConfig = ref({});
const wechatLoginURL = ref("");
const enableVerify = ref(false);
const captchaRef = ref(null);
const ruleFormRef = ref(null);
const ruleForm = reactive({
  username: process.env.VUE_APP_USER || "",
  password: process.env.VUE_APP_PASS || "",
  mobile: "",
  code: "",
});
const currentRules = computed(() => {
  if (enableVerify.value) {
    // 手机号+验证码验证规则
    return {
      mobile: [
        { required: true, trigger: "blur", message: "请输入手机号码" },
        { pattern: /^1[3-9]\d{9}$/, trigger: "blur", message: "请输入正确的手机号码" }
      ],
      code: [{ required: true, trigger: "blur", message: "请输入验证码" }],
    };
  } else {
    // 用户名+密码验证规则
    return {
      username: [{ required: true, trigger: "blur", message: "请输入账号" }],
      password: [{ required: true, trigger: "blur", message: "请输入密码" }],
    };
  }
});
onMounted(() => {
  // 获取系统配置
  getSystemInfo()
    .then((res) => {
      logo.value = res.data.logo;
      title.value = res.data.title;
      enableVerify.value = res.data["enabled_verify"];
    })
    .catch((e) => {
      showMessageError("获取系统配置失败：" + e.message);
    });

  getLicenseInfo()
    .then((res) => {
      licenseConfig.value = res.data;
    })
    .catch((e) => {
      showMessageError("获取 License 配置：" + e.message);
    });

  checkSession()
    .then(() => {
      if (isMobile()) {
        router.push("/mobile");
      } else {
        router.push("/chat");
      }
    })
    .catch(() => {});

  const returnURL = `${location.protocol}//${location.host}/login/callback?action=login`;
  httpGet("/api/user/clogin?return_url=" + returnURL)
    .then((res) => {
      wechatLoginURL.value = res.data.url;
    })
    .catch((e) => {
      console.error(e);
    });
});

const handleKeyup = (e) => {
  if (e.key === "Enter") {
    login();
  }
};

const login = async function () {
  await ruleFormRef.value.validate(async (valid) => {
    if (valid) {
      doLogin({});
    }
  });
};

const store = useSharedStore();
const doLogin = (verifyData) => {
  // 根据登录方式构建不同的请求数据
  let loginData = {
    key: verifyData.key,
    dots: verifyData.dots,
    x: verifyData.x,
  };

  if (enableVerify.value) {
    // 手机号+验证码登录
    loginData.mobile = ruleForm.mobile;
    loginData.code = ruleForm.code;
  } else {
    // 用户名+密码登录
    loginData.username = ruleForm.username;
    loginData.password = ruleForm.password;
  }

  httpPost("/api/user/login", loginData)
    .then((res) => {
      setUserToken(res.data.token);
      store.setIsLogin(true);
      if (isMobile()) {
        router.push("/mobile");
      } else {
        router.push("/chat");
      }
    })
    .catch((e) => {
      showMessageError("登录失败，" + e.message);
    });
};
</script>

<style lang="stylus" scoped>
@import "@/assets/css/login.styl"
</style>
