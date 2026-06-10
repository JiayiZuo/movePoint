const api = require('../../services/api')
const { formatDate, toRFC3339 } = require('../../utils/format')

const typeOptions = [
  {
    label: '抱石',
    value: 'bouldering',
    placeholder: '例如 V3',
    grades: ['V0', 'V1', 'V2', 'V3', 'V4', 'V5', 'V6']
  },
  {
    label: '难度攀爬',
    value: 'sport_climbing',
    placeholder: '例如 5.10b',
    grades: ['5.8', '5.9', '5.10a', '5.10b', '5.10c', '5.11a']
  }
]

const attemptOptions = [
  { label: 'Flash', value: 'flash' },
  { label: '2-3 次', value: '2-3' },
  { label: '4-6 次', value: '4-6' },
  { label: '7 次以上', value: '7+' },
  { label: '未完成', value: 'failed' }
]

const colorOptions = [
  { name: '白色', hex: '#f8fafc' },
  { name: '黄色', hex: '#facc15' },
  { name: '绿色', hex: '#22c55e' },
  { name: '蓝色', hex: '#3b82f6' },
  { name: '红色', hex: '#ef4444' },
  { name: '黑色', hex: '#111827' }
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

function minutesBetween(form) {
  const start = new Date(toRFC3339(form.startDate, form.startTime))
  const end = new Date(toRFC3339(form.endDate, form.endTime))
  const minutes = Math.max(0, Math.round((end.getTime() - start.getTime()) / 60000))
  return minutes
}

function durationText(minutes) {
  if (minutes < 60) return `${minutes} 分`
  const hours = Math.floor(minutes / 60)
  const rest = minutes % 60
  return rest ? `${hours}h${rest}m` : `${hours}h`
}

Page({
  data: {
    id: '',
    loading: false,
    typeOptions,
    attemptOptions,
    colorOptions,
    durationOptions: [
      { label: '45 分', value: 45 },
      { label: '60 分', value: 60 },
      { label: '90 分', value: 90 },
      { label: '120 分', value: 120 }
    ],
    noteTemplates: ['力量不错', '指力疲劳', '脚法要更安静', '下次复盘 beta'],
    ratingOptions: [1, 2, 3, 4, 5],
    typeIndex: 0,
    attemptIndex: 0,
    gradeSuggestions: typeOptions[0].grades,
    durationPreview: '90 分',
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

    if (query.type) {
      const typeIndex = Math.max(0, typeOptions.findIndex((item) => item.value === query.type))
      this.setData({
        typeIndex,
        gradeSuggestions: typeOptions[typeIndex].grades
      })
    }

    if (query.id) {
      this.setData({ id: query.id })
      this.loadRecord(query.id)
    } else {
      this.updateDurationPreview()
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
      gradeSuggestions: typeOptions[typeIndex].grades,
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
    }, () => this.updateDurationPreview())
  },

  onInput(e) {
    this.setData({ [`form.${e.currentTarget.dataset.field}`]: e.detail.value })
  },

  onInputDate(e) {
    this.setData({ [`form.${e.currentTarget.dataset.field}`]: e.detail.value }, () => this.updateDurationPreview())
  },

  selectType(e) {
    const typeIndex = Number(e.currentTarget.dataset.index)
    this.setData({
      typeIndex,
      gradeSuggestions: typeOptions[typeIndex].grades
    })
  },

  onAttemptChange(e) {
    const attemptIndex = Number(e.detail.value)
    this.setData({
      attemptIndex,
      'form.success': attemptOptions[attemptIndex].value !== 'failed'
    })
  },

  onSwitch(e) {
    const success = e.detail.value
    const failedIndex = attemptOptions.findIndex((item) => item.value === 'failed')
    const nextAttemptIndex = success
      ? (this.data.attemptIndex === failedIndex ? 1 : this.data.attemptIndex)
      : failedIndex
    this.setData({
      'form.success': success,
      attemptIndex: nextAttemptIndex
    })
  },

  pickGrade(e) {
    this.setData({ 'form.grade': e.currentTarget.dataset.value })
  },

  pickColor(e) {
    this.setData({ 'form.color': e.currentTarget.dataset.value })
  },

  pickRating(e) {
    this.setData({ 'form.rating': Number(e.currentTarget.dataset.value) })
  },

  applyDuration(e) {
    const minutes = Number(e.currentTarget.dataset.minutes)
    const end = new Date(toRFC3339(this.data.form.endDate, this.data.form.endTime))
    const start = new Date(end.getTime() - minutes * 60000)
    const date = `${start.getFullYear()}-${String(start.getMonth() + 1).padStart(2, '0')}-${String(start.getDate()).padStart(2, '0')}`
    const time = `${String(start.getHours()).padStart(2, '0')}:${String(start.getMinutes()).padStart(2, '0')}`
    this.setData({
      'form.startDate': date,
      'form.startTime': time
    }, () => this.updateDurationPreview())
  },

  appendNote(e) {
    const value = e.currentTarget.dataset.value
    const notes = this.data.form.notes
    this.setData({ 'form.notes': notes ? `${notes}；${value}` : value })
  },

  updateDurationPreview() {
    const minutes = minutesBetween(this.data.form)
    this.setData({ durationPreview: durationText(minutes) })
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
    if (minutesBetween(this.data.form) <= 0) {
      wx.showToast({ title: '结束时间需要晚于开始时间', icon: 'none' })
      return
    }
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
