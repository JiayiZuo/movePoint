const api = require('../../services/api')
const { formatDate, toRFC3339 } = require('../../utils/format')

const typeOptions = [
  { label: '抱石', value: 'bouldering' },
  { label: '难度攀爬', value: 'sport_climbing' }
]

const attemptOptions = [
  { label: 'Flash', value: 'flash' },
  { label: '2-3 次', value: '2-3' },
  { label: '4-6 次', value: '4-6' },
  { label: '7 次以上', value: '7+' },
  { label: '未完成', value: 'failed' }
]

function nowParts(offsetMinutes = 0) {
  const date = new Date(Date.now() + offsetMinutes * 60000)
  const yyyy = date.getFullYear()
  const mm = String(date.getMonth() + 1).padStart(2, '0')
  const dd = String(date.getDate()).padStart(2, '0')
  const hh = String(date.getHours()).padStart(2, '0')
  const mi = String(date.getMinutes()).padStart(2, '0')
  return {
    date: `${yyyy}-${mm}-${dd}`,
    time: `${hh}:${mi}`
  }
}

Page({
  data: {
    id: '',
    loading: false,
    typeOptions,
    attemptOptions,
    typeIndex: 0,
    attemptIndex: 0,
    form: {
      startDate: nowParts(-90).date,
      startTime: nowParts(-90).time,
      endDate: nowParts().date,
      endTime: nowParts().time,
      grade: '',
      color: '',
      success: true,
      rating: 3,
      location: '',
      notes: '',
      media_urls: ''
    }
  },

  onLoad(query) {
    if (!wx.getStorageSync('token')) {
      wx.reLaunch({ url: '/pages/login/login' })
      return
    }
    if (query.id) {
      this.setData({ id: query.id })
      this.loadRecord(query.id)
    }
  },

  async loadRecord(id) {
    const record = await api.getRecord(id)
    const startDate = formatDate(record.start_time)
    const endDate = formatDate(record.end_time)
    const start = new Date(record.start_time)
    const end = new Date(record.end_time)
    const typeIndex = Math.max(0, typeOptions.findIndex((item) => item.value === record.type))
    const attemptIndex = Math.max(0, attemptOptions.findIndex((item) => item.value === record.attempts))

    this.setData({
      typeIndex,
      attemptIndex,
      form: {
        startDate,
        startTime: `${String(start.getHours()).padStart(2, '0')}:${String(start.getMinutes()).padStart(2, '0')}`,
        endDate,
        endTime: `${String(end.getHours()).padStart(2, '0')}:${String(end.getMinutes()).padStart(2, '0')}`,
        grade: record.grade || '',
        color: record.color || '',
        success: !!record.success,
        rating: record.rating || 3,
        location: record.location || '',
        notes: record.notes || '',
        media_urls: record.media_urls || ''
      }
    })
  },

  onInput(e) {
    this.setData({ [`form.${e.currentTarget.dataset.field}`]: e.detail.value })
  },

  onInputDate(e) {
    this.setData({ [`form.${e.currentTarget.dataset.field}`]: e.detail.value })
  },

  onTypeChange(e) {
    this.setData({ typeIndex: Number(e.detail.value) })
  },

  onAttemptChange(e) {
    this.setData({ attemptIndex: Number(e.detail.value) })
  },

  onSwitch(e) {
    this.setData({ 'form.success': e.detail.value })
  },

  onSlider(e) {
    this.setData({ 'form.rating': e.detail.value })
  },

  buildPayload() {
    const form = this.data.form
    return {
      type: typeOptions[this.data.typeIndex].value,
      start_time: toRFC3339(form.startDate, form.startTime),
      end_time: toRFC3339(form.endDate, form.endTime),
      grade: form.grade,
      color: form.color,
      attempts: attemptOptions[this.data.attemptIndex].value,
      success: form.success,
      rating: Number(form.rating),
      location: form.location,
      notes: form.notes,
      media_urls: form.media_urls
    }
  },

  async submit() {
    const payload = this.buildPayload()
    if (!payload.grade || !payload.location) {
      wx.showToast({ title: '请填写难度和地点', icon: 'none' })
      return
    }

    this.setData({ loading: true })
    try {
      if (this.data.id) {
        await api.updateRecord(this.data.id, payload)
      } else {
        await api.createRecord(payload)
      }
      wx.showToast({ title: '已保存' })
      wx.navigateBack()
    } finally {
      this.setData({ loading: false })
    }
  },

  remove() {
    wx.showModal({
      title: '删除记录',
      content: '删除后无法恢复，确定继续吗？',
      success: async (res) => {
        if (!res.confirm) return
        await api.deleteRecord(this.data.id)
        wx.showToast({ title: '已删除' })
        wx.navigateBack()
      }
    })
  }
})
