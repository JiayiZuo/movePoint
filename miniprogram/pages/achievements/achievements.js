const api = require('../../services/api')

Page({
  data: {
    achievements: []
  },

  onShow() {
    if (!wx.getStorageSync('token')) {
      wx.reLaunch({ url: '/pages/login/login' })
      return
    }
    this.loadAchievements()
  },

  async loadAchievements() {
    const achievements = await api.getAchievements()
    this.setData({
      achievements: (achievements || []).map((item) => ({
        ...item,
        progress: Math.round(item.progress || 0)
      }))
    })
  },

  async checkAchievements() {
    await api.checkAchievements()
    wx.showToast({ title: '已刷新' })
    this.loadAchievements()
  }
})
