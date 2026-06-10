const { request } = require('../utils/request')

const api = {
  login(data) {
    return request({ url: '/login', method: 'POST', data })
  },

  register(data) {
    return request({ url: '/register', method: 'POST', data })
  },

  getProfile() {
    return request({ url: '/profile' })
  },

  updateProfile(data) {
    return request({ url: '/profile', method: 'PUT', data })
  },

  getStats() {
    return request({ url: '/profile/stats' })
  },

  getAchievements() {
    return request({ url: '/profile/achievements' })
  },

  checkAchievements() {
    return request({ url: '/profile/check-achievements', method: 'POST' })
  },

  getRecords(params = {}) {
    const query = Object.keys(params)
      .filter((key) => params[key] !== undefined && params[key] !== '')
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`)
      .join('&')
    return request({ url: `/records${query ? `?${query}` : ''}` })
  },

  getRecord(id) {
    return request({ url: `/records/${id}` })
  },

  createRecord(data) {
    return request({ url: '/records', method: 'POST', data })
  },

  updateRecord(id, data) {
    return request({ url: `/records/${id}`, method: 'PUT', data })
  },

  deleteRecord(id) {
    return request({ url: `/records/${id}`, method: 'DELETE' })
  },

  getAnalysis(params = {}) {
    const query = Object.keys(params)
      .filter((key) => params[key] !== undefined && params[key] !== '')
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`)
      .join('&')
    return request({ url: `/analysis/climbing${query ? `?${query}` : ''}` })
  }
}

module.exports = api
