// pages/achievements/achievements.js
Page({
  data: {
    achievements: [],
    unlockedCount: 0,
    totalCount: 0
  },

  onLoad() {
    this.loadAchievements();
  },

  loadAchievements() {
    const app = getApp();
    const token = app.getGlobalToken();
    const userInfo = app.globalData.userInfo;
    
    if (!token || !userInfo || !userInfo.openId) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      });
      return;
    }
    
    // 请求用户成就数据
    wx.request({
      url: `${app.globalData.serverUrl}/profile/achievements`,
      method: 'GET',
      header: {
        'Authorization': token,
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.statusCode === 200) {
          const achievements = res.data.achievements || [];
          const unlockedCount = achievements.filter(item => item.completed).length;
          
          this.setData({
            achievements: achievements,
            unlockedCount: unlockedCount,
            totalCount: achievements.length
          });
        } else {
          wx.showToast({
            title: '获取成就失败',
            icon: 'none'
          });
        }
      },
      fail: (err) => {
        console.error('获取成就失败:', err);
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        });
      }
    });
  },

  onAchievementTap(e) {
    const index = e.currentTarget.dataset.index;
    const achievement = this.data.achievements[index];
    
    // 显示成就详情
    wx.showModal({
      title: achievement.name,
      content: achievement.description,
      showCancel: false,
      confirmText: '确定'
    });
  }
})