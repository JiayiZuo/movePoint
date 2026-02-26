// pages/index/index.js
Page({
  data: {
    userInfo: null,
    hasUserInfo: false,
    canIUseGetUserProfile: false,
    dailyStats: {
      sessions: 0,
      duration: 0,
      calories: 0
    },
    recentRecords: []
  },

  onLoad() {
    // 检测是否支持getUserProfile
    if (wx.getUserProfile) {
      this.setData({
        canIUseGetUserProfile: true
      })
    }
    
    // 检查登录状态
    this.checkLoginStatus();
  },

  onShow() {
    // 页面显示时刷新数据
    this.checkLoginStatus();
  },

  checkLoginStatus() {
    const token = wx.getStorageSync('token');
    if (token) {
      this.setData({
        hasUserInfo: true
      });
      // 获取用户信息和统计数据
      this.fetchUserData();
    }
  },

  fetchUserData() {
    const token = wx.getStorageSync('token');
    if (!token) {
      wx.navigateTo({
        url: '/pages/login/login'
      });
      return;
    }

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
            dailyStats: {
              sessions: res.data.weekly_sessions || 0,
              duration: Math.round((res.data.total_duration || 0) / 60) || 0, // 转换为小时
              calories: Math.round(res.data.total_calories || 0)
            }
          });
        }
      },
      fail: (err) => {
        console.error('获取用户统计数据失败:', err);
      }
    });

    // 获取最近的攀岩记录
    wx.request({
      url: `${getApp().globalData.serverUrl}/records?page=1&limit=5`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.statusCode === 200) {
          this.setData({
            recentRecords: res.data.data || []
          });
        }
      },
      fail: (err) => {
        console.error('获取最近记录失败:', err);
      }
    });
  },

  getUserProfile(e) {
    // 推荐使用wx.getUserProfile获取用户信息，开发者每次通过该接口获取用户个人信息均需用户确认
    wx.getUserProfile({
      desc: '用于完善会员资料', // 声明获取用户个人信息后的用途，后续会展示在弹窗中，请谨慎修改
      success: (res) => {
        console.log(res);
        this.setData({
          userInfo: res.userInfo,
          hasUserInfo: true
        });
      }
    });
  },

  getUserInfo(e) {
    // 不推荐使用getUserInfo获取用户信息，预计自2021年4月13日起，getUserInfo将不再弹出弹窗，并直接返回匿名的用户个人信息
    console.log(e);
    this.setData({
      userInfo: e.detail.userInfo,
      hasUserInfo: true
    });
  },

  login() {
    wx.navigateTo({
      url: '/pages/login/login'
    });
  },

  navigateToAddRecord() {
    wx.navigateTo({
      url: '/pages/add-record/add-record'
    });
  },

  navigateToAnalysis() {
    wx.navigateTo({
      url: '/pages/analysis/analysis'
    });
  },

  navigateToRecords() {
    wx.navigateTo({
      url: '/pages/records/records'
    });
  }
});