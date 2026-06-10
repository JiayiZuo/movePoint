const api = require('../../services/api')

function dateText(date) {
  const yyyy = date.getFullYear()
  const mm = String(date.getMonth() + 1).padStart(2, '0')
  const dd = String(date.getDate()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd}`
}

Page({
  data: {
    from: dateText(new Date(Date.now() - 90 * 24 * 60 * 60 * 1000)),
    to: dateText(new Date()),
    analysis: { summary: {} },
    gradeRows: [],
    monthlyRows: [],
    durationHours: '0',
    calories: '0'
  },

  onShow() {
    if (!wx.getStorageSync('token')) {
      wx.reLaunch({ url: '/pages/login/login' })
      return
    }
    this.loadAnalysis()
  },

  onPullDownRefresh() {
    this.loadAnalysis().finally(() => wx.stopPullDownRefresh())
  },

  onDate(e) {
    this.setData({ [e.currentTarget.dataset.field]: e.detail.value }, () => this.loadAnalysis())
  },

  async loadAnalysis() {
    const analysis = await api.getAnalysis({ from: this.data.from, to: this.data.to })
    const distribution = analysis.grade_distribution || {}
    const rates = analysis.success_rate_by_grade || {}
    const gradeRows = Object.keys(distribution).map((grade) => ({
      grade,
      attempts: distribution[grade].attempts || 0,
      success: distribution[grade].success || 0,
      rate: Math.round(rates[grade] || 0)
    }))

    this.setData({
      analysis,
      gradeRows,
      monthlyRows: analysis.monthly_trends || [],
      durationHours: ((analysis.summary.total_duration || 0) / 60).toFixed(1),
      calories: Math.round(analysis.summary.total_calories || 0)
    })
  }
})
