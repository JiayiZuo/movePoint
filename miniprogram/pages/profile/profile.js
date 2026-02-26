// pages/profile/profile.js
Page({
  data: {
    userInfo: null,
    profile: null,
    achievements: [],
    stats: {}
  },

  onLoad() {
    this.loadProfileData();
  },

  onShow() {
    // 页面显示时刷新数据
    this.loadProfileData();
  },

  loadProfileData() {
    const token = wx.getStorageSync('token');
    if (!token) {
      wx.navigateTo({
        url: '/pages/login/login'
      });
      return;
    }

    // 获取用户个人信息
    wx.request({
      url: `${getApp().globalData.serverUrl}/profile`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.statusCode === 200) {
          this.setData({
            profile: res.data
          });
        }
      },
      fail: (err) => {
        console.error('获取用户信息失败:', err);
      }
    });

    // 获取用户统计数据
    wx.request({
      url: `${getApp().globalData.serverUrl}/profile/stats`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.statusCode === 200) {
          this.setData({
            stats: res.data
          });
        }
      },
      fail: (err) => {
        console.error('获取统计信息失败:', err);
      }
    });

    // 获取用户成就
    wx.request({
      url: `${getApp().globalData.serverUrl}/profile/achievements`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.statusCode === 200) {
          this.setData({
            achievements: res.data
          });
        }
      },
      fail: (err) => {
        console.error('获取成就失败:', err);
      }
    });
  },

  updateProfile() {
    wx.navigateTo({
      url: '/pages/update-profile/update-profile'
    });
  },

  checkAchievements() {
    const token = wx.getStorageSync('token');
    if (!token) {
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
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.statusCode === 200) {
          wx.showToast({
            title: '成就检查完成',
            icon: 'success'
          });
          // 重新加载成就数据
          this.loadProfileData();
        }
      },
      fail: (err) => {
        console.error('检查成就失败:', err);
      }
    });
  }
});