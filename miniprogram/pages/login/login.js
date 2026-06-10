const api = require('../../services/api')

const app = getApp()

Page({
  data: {
    email: '',
    password: '',
    loading: false
  },

  onLoad() {
    if (wx.getStorageSync('token')) {
      wx.switchTab({ url: '/pages/dashboard/dashboard' })
    }
  },

  onInput(e) {
    this.setData({ [e.currentTarget.dataset.field]: e.detail.value })
  },

  async submit() {
    const { email, password } = this.data
    if (!email || !password) {
      wx.showToast({ title: '请输入邮箱和密码', icon: 'none' })
      return
    }

    this.setData({ loading: true })
    try {
      const res = await api.login({ email, password })
      const user = {
        id: res.user_id,
        username: res.username,
        email: res.email
      }
      wx.setStorageSync('token', res.token)
      wx.setStorageSync('user', user)
      app.globalData.user = user
      wx.switchTab({ url: '/pages/dashboard/dashboard' })
    } finally {
      this.setData({ loading: false })
    }
  },

  goRegister() {
    wx.navigateTo({ url: '/pages/register/register' })
  }
})
