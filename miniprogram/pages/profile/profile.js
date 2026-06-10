const api = require('../../services/api')
const { formatDate } = require('../../utils/format')

const app = getApp()

Page({
  data: {
    avatarLetter: 'M',
    user: {},
    profile: {},
    birthDateText: '',
    editing: false,
    saving: false,
    form: {
      avatar_url: '',
      weight: '',
      height: '',
      birth_date: '',
      bio: ''
    }
  },

  onShow() {
    if (!wx.getStorageSync('token')) {
      wx.reLaunch({ url: '/pages/login/login' })
      return
    }
    this.loadProfile()
  },

  async loadProfile() {
    const profile = await api.getProfile()
    const birthDateText = formatDate(profile.birth_date)
    this.setData({
      user: wx.getStorageSync('user') || {},
      profile,
      birthDateText,
      avatarLetter: (profile.username || 'M').slice(0, 1).toUpperCase(),
      form: {
        avatar_url: profile.avatar_url || '',
        weight: profile.weight || '',
        height: profile.height || '',
        birth_date: birthDateText,
        bio: profile.bio || ''
      }
    })
  },

  startEdit() {
    this.setData({ editing: true })
  },

  cancelEdit() {
    const profile = this.data.profile
    this.setData({
      editing: false,
      form: {
        avatar_url: profile.avatar_url || '',
        weight: profile.weight || '',
        height: profile.height || '',
        birth_date: this.data.birthDateText,
        bio: profile.bio || ''
      }
    })
  },

  onInput(e) {
    this.setData({ [`form.${e.currentTarget.dataset.field}`]: e.detail.value })
  },

  onBirthDate(e) {
    this.setData({ 'form.birth_date': e.detail.value })
  },

  async saveProfile() {
    const form = this.data.form
    this.setData({ saving: true })
    try {
      await api.updateProfile({
        avatar_url: form.avatar_url,
        weight: Number(form.weight || 0),
        height: Number(form.height || 0),
        birth_date: form.birth_date,
        bio: form.bio
      })
      wx.showToast({ title: '已保存' })
      this.setData({ editing: false })
      await this.loadProfile()
    } finally {
      this.setData({ saving: false })
    }
  },

  goAchievements() {
    wx.navigateTo({ url: '/pages/achievements/achievements' })
  },

  logout() {
    wx.removeStorageSync('token')
    wx.removeStorageSync('user')
    app.globalData.user = null
    wx.reLaunch({ url: '/pages/login/login' })
  }
})
