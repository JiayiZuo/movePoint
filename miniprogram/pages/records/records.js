const api = require('../../services/api')
const { formatDateTime, minutesText } = require('../../utils/format')

const typeMap = {
  bouldering: '抱石',
  sport_climbing: '难度'
}

const attemptsMap = {
  flash: 'Flash',
  '2-3': '2-3 次',
  '4-6': '4-6 次',
  '7+': '7 次以上',
  failed: '未完成'
}

Page({
  data: {
    records: [],
    page: 1,
    limit: 20,
    total: 0,
    hasMore: false,
    from: '',
    to: ''
  },

  onShow() {
    if (!wx.getStorageSync('token')) {
      wx.reLaunch({ url: '/pages/login/login' })
      return
    }
    this.refresh()
  },

  onPullDownRefresh() {
    this.refresh().finally(() => wx.stopPullDownRefresh())
  },

  async refresh() {
    this.setData({ page: 1 })
    await this.loadRecords(true)
  },

  async loadRecords(reset = false) {
    const { page, limit, from, to } = this.data
    const res = await api.getRecords({ page, limit, from, to })
    const list = (res.data || []).map((item) => ({
      ...item,
      typeText: typeMap[item.type] || item.type,
      attemptsText: attemptsMap[item.attempts] || item.attempts,
      startText: formatDateTime(item.start_time),
      durationText: minutesText(item.duration)
    }))
    const records = reset ? list : this.data.records.concat(list)
    this.setData({
      records,
      total: res.total || 0,
      hasMore: records.length < (res.total || 0)
    })
  },

  loadMore() {
    if (!this.data.hasMore) return
    this.setData({ page: this.data.page + 1 }, () => this.loadRecords())
  },

  onFilterDate(e) {
    this.setData({ [e.currentTarget.dataset.field]: e.detail.value }, () => this.refresh())
  },

  addRecord() {
    wx.navigateTo({ url: '/pages/record-form/record-form' })
  },

  openRecord(e) {
    wx.navigateTo({ url: `/pages/record-form/record-form?id=${e.currentTarget.dataset.id}` })
  }
})
