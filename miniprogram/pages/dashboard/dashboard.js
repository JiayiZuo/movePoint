const api = require('../../services/api')
const { formatDateTime, minutesText } = require('../../utils/format')

const typeMap = {
  bouldering: '抱石',
  sport_climbing: '难度'
}

Page({
  data: {
    user: {},
    profile: {},
    stats: {},
    records: [],
    durationHours: '0'
  },

  onShow() {
    if (!wx.getStorageSync('token')) {
      wx.reLaunch({ url: '/pages/login/login' })
      return
    }
    this.loadData()
  },

  async loadData() {
    try {
      const [profile, stats, recordRes] = await Promise.all([
        api.getProfile(),
        api.getStats(),
        api.getRecords({ page: 1, limit: 3 })
      ])

      const records = (recordRes.data || []).map((item) => ({
        ...item,
        typeText: typeMap[item.type] || item.type,
        startText: formatDateTime(item.start_time),
        durationText: minutesText(item.duration)
      }))

      this.setData({
        user: wx.getStorageSync('user') || {},
        profile,
        stats,
        records,
        durationHours: ((stats.total_duration || 0) / 60).toFixed(1)
      })
    } catch (err) {
      console.warn(err)
    }
  },

  addRecord() {
    wx.navigateTo({ url: '/pages/record-form/record-form' })
  },

  quickAdd(e) {
    wx.navigateTo({ url: `/pages/record-form/record-form?type=${e.currentTarget.dataset.type}` })
  },

  goRecords() {
    wx.switchTab({ url: '/pages/records/records' })
  },

  openRecord(e) {
    wx.navigateTo({ url: `/pages/record-form/record-form?id=${e.currentTarget.dataset.id}` })
  }
})
