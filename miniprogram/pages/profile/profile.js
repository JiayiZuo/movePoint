// pages/profile/profile.js
Page({
  data: {
    userInfo: null,
    profile: null,
    achievements: [],
    stats: {}
  },

  onLoad() {
    this.checkLoginStatus();
  },

  onShow() {
    // 页面显示时检查登录状态
    this.checkLoginStatus();
  },

  checkLoginStatus() {
    // 检查本地是否有用户信息
    const userInfo = wx.getStorageSync('userInfo');
    if (!userInfo || !userInfo.openId) {
      // 未登录，不执行后续操作，等待用户手动登录
      this.setData({
        profile: null
      });
      return;
    }

    // 已登录，加载用户数据
    this.setData({
      profile: userInfo,
      stats: wx.getStorageSync('userStats') || {}
    });
  },

  goToLogin() {
    wx.navigateTo({
      url: '/pages/login/login'
    });
  },

  updateProfile() {
    const userInfo = wx.getStorageSync('userInfo');
    if (!userInfo || !userInfo.openId) {
      wx.navigateTo({
        url: '/pages/login/login'
      });
      return;
    }
    wx.navigateTo({
      url: '/pages/update-profile/update-profile'
    });
  },

  checkAchievements() {
    const userInfo = wx.getStorageSync('userInfo');
    if (!userInfo || !userInfo.openId) {
      wx.navigateTo({
        url: '/pages/login/login'
      });
      return;
    }

    // 检查并更新成就
    wx.request({
      url: `${getApp().globalData.serverUrl}/profile/check-achievements`,
      method: 'POST',
      header: {
        'Authorization': wx.getStorageSync('token'),
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.statusCode === 200) {
          wx.showToast({
            title: '成就检查完成',
            icon: 'success'
          });
          // 重新加载成就数据
          this.checkLoginStatus();
        }
      },
      fail: (err) => {
        console.error('检查成就失败:', err);
      }
    });
  }
});