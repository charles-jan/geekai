<template>
  <div class="agreement-page">
    <div class="agreement-container custom-scroll">
      <div class="agreement-header">
        <h1>服务条款、隐私政策与免责声明</h1>
        <p class="last-update">最后更新时间：{{ lastUpdateTime }}</p>
      </div>

      <el-tabs v-model="activeTab" class="agreement-tabs">
        <el-tab-pane label="服务条款" name="terms">
          <div class="agreement-content" v-html="termsContent"></div>
        </el-tab-pane>
        <el-tab-pane label="隐私政策" name="privacy">
          <div class="agreement-content" v-html="privacyContent"></div>
        </el-tab-pane>
        <el-tab-pane label="免责声明" name="disclaimer">
          <div class="agreement-content" v-html="disclaimerContent"></div>
        </el-tab-pane>
      </el-tabs>

      <div class="agreement-footer">
        <el-button type="primary" @click="goBack">返回</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import MarkdownIt from "markdown-it";

const router = useRouter();
const activeTab = ref("terms");
const lastUpdateTime = ref(new Date().toLocaleDateString("zh-CN"));
const termsContent = ref("");
const privacyContent = ref("");
const disclaimerContent = ref("");

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
});

onMounted(async () => {
  try {
    // 加载服务条款
    const termsResponse = await fetch("/text/terms-of-service.txt");
    const termsText = await termsResponse.text();
    termsContent.value = md.render(termsText);

    // 加载隐私政策
    const privacyResponse = await fetch("/text/privacy.txt");
    const privacyText = await privacyResponse.text();
    privacyContent.value = md.render(privacyText);

    // 加载免责声明
    const disclaimerResponse = await fetch("/text/disclaimer.txt");
    const disclaimerText = await disclaimerResponse.text();
    disclaimerContent.value = md.render(disclaimerText);
  } catch (error) {
    console.error("加载协议内容失败:", error);
  }
});

const goBack = () => {
  router.back();
};
</script>

<style lang="stylus" scoped>
.agreement-page
  min-height 100vh
  background #f5f7fa
  padding 20px

.agreement-container
  max-width 1200px
  margin 0 auto
  background white
  border-radius 8px
  box-shadow 0 2px 12px rgba(0, 0, 0, 0.1)
  padding 40px
  max-height calc(100vh - 40px)
  overflow-y auto

.agreement-header
  text-align center
  margin-bottom 30px
  border-bottom 2px solid #e4e7ed
  padding-bottom 20px

  h1
    font-size 32px
    font-weight 600
    color #303133
    margin 0 0 10px 0

  .last-update
    font-size 14px
    color #909399
    margin 0

.agreement-tabs
  margin-bottom 30px

  :deep(.el-tabs__item)
    font-size 16px
    font-weight 500

  :deep(.el-tabs__nav-wrap::after)
    background-color #e4e7ed

.agreement-content
  font-size 15px
  line-height 1.8
  color #606266
  padding 20px 0

  :deep(h1)
    font-size 28px
    font-weight 600
    color #303133
    margin 30px 0 20px 0
    padding-bottom 10px
    border-bottom 2px solid #409eff

  :deep(h2)
    font-size 22px
    font-weight 600
    color #303133
    margin 25px 0 15px 0

  :deep(h3)
    font-size 18px
    font-weight 600
    color #606266
    margin 20px 0 10px 0

  :deep(p)
    margin 12px 0
    text-indent 2em

  :deep(blockquote)
    background #ecf5ff
    border-left 4px solid #409eff
    padding 10px 15px
    margin 15px 0
    color #606266

    p
      margin 5px 0
      text-indent 0

  :deep(ul), :deep(ol)
    margin 10px 0
    padding-left 30px

  :deep(li)
    margin 8px 0
    line-height 1.6

  :deep(strong)
    color #303133
    font-weight 600

  :deep(code)
    background #f4f4f5
    padding 2px 6px
    border-radius 3px
    font-family monospace
    color #e6a23c

.agreement-footer
  text-align center
  margin-top 30px
  padding-top 20px
  border-top 1px solid #e4e7ed

  .el-button
    padding 12px 40px
    font-size 16px

@media (max-width: 768px)
  .agreement-page
    padding 10px

  .agreement-container
    padding 20px
    border-radius 0

  .agreement-header
    h1
      font-size 24px

  .agreement-content
    font-size 14px

    :deep(h1)
      font-size 22px

    :deep(h2)
      font-size 18px

    :deep(h3)
      font-size 16px
</style>
