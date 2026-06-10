const api = require('../../services/api')
const { formatDate } = require('../../utils/format')

const app = getApp()

Page({
  data: {
    avatarLetter: 'M',
    user: {},
    profile: {},
    achievements: [],
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
    const [profile, achievements] = await Promise.all([
      api.getProfile(),
      api.getAchievements()
    ])
    this.setData({
      user: wx.getStorageSync('user') || {},
      profile,
      avatarLetter: (profile.username || 'M').slice(0, 1).toUpperCase(),
      achievements: (achievements || []).map((item) => ({
        ...item,
        progress: Math.round(item.progress || 0)
      })),
      form: {
        avatar_url: profile.avatar_url || '',
        weight: profile.weight || '',
        height: profile.height || '',
        birth_date: formatDate(profile.birth_date),
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
      this.loadProfile()
    } finally {
      this.setData({ saving: false })
    }
  },

  async checkAchievements() {
    await api.checkAchievements()
    wx.showToast({ title: '已刷新' })
    this.loadProfile()
  },

  logout() {
    wx.removeStorageSync('token')
    wx.removeStorageSync('user')
    app.globalData.user = null
    wx.reLaunch({ url: '/pages/login/login' })
  }
})
