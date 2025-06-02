// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import {createRouter, createWebHistory} from "vue-router";

const routes = [
  {
    name: "Index",
    path: "/",
    meta: { title: "首页" },
    component: () => import("@/views/Index.vue"),
  },
  {
    name: "home",
    path: "/home",
    redirect: "/chat",
    component: () => import("@/views/Home.vue"),
    children: [
      {
        name: "chat",
        path: "/chat",
        meta: { title: "创作中心" },
        component: () => import("@/views/ChatPlus.vue"),
      },
      {
        name: "chat-id",
        path: "/chat/:id",
        meta: { title: "创作中心" },
        component: () => import("@/views/ChatPlus.vue"),
      },
      {
        name: "image-mj",
        path: "/mj",
        meta: { title: "MidJourney 绘画中心" },
        component: () => import("@/views/ImageMj.vue"),
      },
      {
        name: "image-sd",
        path: "/sd",
        meta: { title: "stable diffusion 绘画中心" },
        component: () => import("@/views/ImageSd.vue"),
      },
      {
        name: "member",
        path: "/member",
        meta: { title: "会员充值中心" },
        component: () => import("@/views/Member.vue"),
      },
      {
        name: "chat-app",
        path: "/apps",
        meta: { title: "应用中心" },
        component: () => import("@/views/ChatApps.vue"),
      },
      {
        name: "images",
        path: "/images-wall",
        meta: { title: "作品展示" },
        component: () => import("@/views/ImagesWall.vue"),
      },
      {
        name: "user-invitation",
        path: "/invite",
        meta: { title: "推广计划" },
        component: () => import("@/views/Invitation.vue"),
      },
      {
        name: "powerLog",
        path: "/powerLog",
        meta: { title: "消费日志" },
        component: () => import("@/views/PowerLog.vue"),
      },
      {
        name: "xmind",
        path: "/xmind",
        meta: { title: "思维导图" },
        component: () => import("@/views/MarkMap.vue"),
      },
      {
        name: "dalle",
        path: "/dalle",
        meta: { title: "DALLE-3" },
        component: () => import("@/views/Dalle.vue"),
      },
      {
        name: "suno",
        path: "/suno",
        meta: { title: "Suno音乐创作" },
        component: () => import("@/views/Suno.vue"),
      },
      {
        name: "ExternalLink",
        path: "/external",
        component: () => import("@/views/ExternalPage.vue"),
      },
      {
        name: "song",
        path: "/song/:id",
        meta: { title: "Suno音乐播放" },
        component: () => import("@/views/Song.vue"),
      },
      {
        name: "luma",
        path: "/luma",
        meta: { title: "Luma视频创作" },
        component: () => import("@/views/Luma.vue"),
      },
    ],
  },
  {
    name: "chat-export",
    path: "/chat/export",
    meta: { title: "导出会话记录" },
    component: () => import("@/views/ChatExport.vue"),
  },
  {
    name: "login",
    path: "/login",
    meta: { title: "用户登录" },
    component: () => import("@/views/Login.vue"),
  },

  {
    name: "login-callback",
    path: "/login/callback",
    meta: { title: "用户登录" },
    component: () => import("@/views/LoginCallback.vue"),
  },
  {
    name: "register",
    path: "/register",

    meta: { title: "用户注册" },
    component: () => import("@/views/Register.vue"),
  },
  {
    name: "resetpassword",
    path: "/resetpassword",
    meta: { title: "重置密码" },
    component: () => import("@/views/Resetpassword.vue"),
  },
  {
    path: "/admin/login",
    name: "admin-login",
    meta: { title: "控制台登录" },
    component: () => import("@/views/admin/Login.vue"),
  },
  {
    path: "/payReturn",
    name: "pay-return",
    meta: { title: "支付回调" },
    component: () => import("@/views/PayReturn.vue"),
  },
  {
    name: "admin",
    path: "/admin",
    redirect: "/admin/dashboard",
    component: () => import("@/views/admin/Home.vue"),
    meta: { title: "Geek-AI 控制台" },
    children: [
      {
        path: "/admin/dashboard",
        name: "admin-dashboard",
        meta: { title: "仪表盘" },
        component: () => import("@/views/admin/Dashboard.vue"),
      },
      {
        path: "/admin/system",
        name: "admin-system",
        meta: { title: "系统设置" },
        component: () => import("@/views/admin/SysConfig.vue"),
      },
      {
        path: "/admin/user",
        name: "admin-user",
        meta: { title: "用户管理" },
        component: () => import("@/views/admin/Users.vue"),
      },
      {
        path: "/admin/app",
        name: "admin-app",
        meta: { title: "应用列表" },
        component: () => import("@/views/admin/Apps.vue"),
      },
      {
        path: "/admin/app/type",
        name: "admin-app-type",
        meta: { title: "应用分类" },
        component: () => import("@/views/admin/AppType.vue"),
      },
      {
        path: "/admin/apikey",
        name: "admin-apikey",
        meta: { title: "API-KEY 管理" },
        component: () => import("@/views/admin/ApiKey.vue"),
      },
      {
        path: "/admin/chat/model",
        name: "admin-chat-model",
        meta: { title: "语言模型" },
        component: () => import("@/views/admin/ChatModel.vue"),
      },
      {
        path: "/admin/product",
        name: "admin-product",
        meta: { title: "充值产品" },
        component: () => import("@/views/admin/Product.vue"),
      },
      {
        path: "/admin/order",
        name: "admin-order",
        meta: { title: "充值订单" },
        component: () => import("@/views/admin/Order.vue"),
      },
      {
        path: "/admin/redeem",
        name: "admin-redeem",
        meta: { title: "兑换码管理" },
        component: () => import("@/views/admin/Redeem.vue"),
      },
      {
        path: "/admin/loginLog",
        name: "admin-loginLog",
        meta: { title: "登录日志" },
        component: () => import("@/views/admin/LoginLog.vue"),
      },
      {
        path: "/admin/functions",
        name: "admin-functions",
        meta: { title: "函数管理" },
        component: () => import("@/views/admin/Functions.vue"),
      },
      {
        path: "/admin/chats",
        name: "admin-chats",
        meta: { title: "对话管理" },
        component: () => import("@/views/admin/ChatList.vue"),
      },
      {
        path: "/admin/images",
        name: "admin-images",
        meta: { title: "绘图管理" },
        component: () => import("@/views/admin/ImageList.vue"),
      },
      {
        path: "/admin/medias",
        name: "admin-medias",
        meta: { title: "音视频管理" },
        component: () => import("@/views/admin/Medias.vue"),
      },
      {
        path: "/admin/powerLog",
        name: "admin-power-log",
        meta: { title: "算力日志" },
        component: () => import("@/views/admin/PowerLog.vue"),
      },
      {
        path: "/admin/manger",
        name: "admin-manger",
        meta: { title: "管理员" },
        component: () => import("@/views/admin/Manager.vue"),
      },
    ],
  },

  {
    name: "mobile-login",
    path: "/mobile/login",
    meta: { title: "用户登录" },
    component: () => import("@/views/mobile/Login.vue"),
  },
  {
    name: "mobile",
    path: "/mobile",
    meta: { title: "首页" },
    component: () => import("@/views/mobile/Home.vue"),
    redirect: "/mobile/index",
    children: [
      {
        path: "/mobile/index",
        name: "mobile-index",
        component: () => import("@/views/mobile/Index.vue"),
      },
      {
        path: "/mobile/chat",
        name: "mobile-chat",
        component: () => import("@/views/mobile/ChatList.vue"),
      },
      {
        path: "/mobile/image",
        name: "mobile-image",
        component: () => import("@/views/mobile/Image.vue"),
      },
      {
        path: "/mobile/profile",
        name: "mobile-profile",
        component: () => import("@/views/mobile/Profile.vue"),
      },
      {
        path: "/mobile/imgWall",
        name: "mobile-img-wall",
        component: () => import("@/views/mobile/pages/ImgWall.vue"),
      },
      {
        path: "/mobile/chat/session",
        name: "mobile-chat-session",
        component: () => import("@/views/mobile/ChatSession.vue"),
      },
      {
        path: "/mobile/chat/export",
        name: "mobile-chat-export",
        component: () => import("@/views/mobile/ChatExport.vue"),
      },
    ],
  },

  {
    name: "test",
    path: "/test",
    meta: { title: "测试页面" },
    component: () => import("@/views/Test.vue"),
  },
  {
    name: "test2",
    path: "/test2",
    meta: { title: "测试页面" },
    component: () => import("@/views/RealtimeTest.vue"),
  },
  {
    name: "NotFound",
    path: "/:all(.*)",
    meta: { title: "页面没有找到" },
    component: () => import("@/views/404.vue"),
  },
];

// console.log(MY_VARIABLE)
const router = createRouter({
  history: createWebHistory(),
  routes: routes,
});

let prevRoute = null;
// dynamic change the title when router change
router.beforeEach((to, from, next) => {
  let suffix = to.meta.suffix || "上海睿邦米企业管理咨询有限公司";
  document.title = `${to.meta.title} - ${suffix}`;
  prevRoute = from;
  next();
});

export { router, prevRoute };
