Page({
  data: {
    title: '',
    content: '',
    grade: '5.8',
    grades: ['5.5', '5.6', '5.7', '5.8', '5.9', '5.10a', '5.10b', '5.10c', '5.10d', '5.11a', '5.11b', '5.11c', '5.11d', '5.12a', '5.12b', '5.12c', '5.12d'],
    type: '抱石',
    types: ['抱石', '运动攀登', '传统攀登', '冰雪攀登', '室内攀岩', '户外自然岩壁'],
    date: new Date().toISOString().split('T')[0],
    duration: '',
    completed: true
  },

  onLoad() {},

  onTitleInput(e) {
    this.setData({
      title: e.detail.value
    });
  },

  onContentInput(e) {
    this.setData({
      content: e.detail.value
    });
  },

  onGradeChange(e) {
    this.setData({
      grade: this.data.grades[e.detail.value]
    });
  },

  onTypeChange(e) {
    this.setData({
      type: this.data.types[e.detail.value]
    });
  },

  onDateChange(e) {
    this.setData({
      date: e.detail.value
    });
  },

  onDurationInput(e) {
    this.setData({
      duration: e.detail.value
    });
  },

  onCompletedChange(e) {
    this.setData({
      completed: e.detail.value
    });
  },

  saveRecord() {
    const { title, content, grade, type, date, duration, completed } = this.data;
    
    if (!title.trim()) {
      wx.showToast({
        title: '请输入攀岩路线名称',
        icon: 'none'
      });
      return;
    }

    // 验证用时是否为有效数字
    if (duration && isNaN(duration)) {
      wx.showToast({
        title: '请输入有效的用时',
        icon: 'none'
      });
      return;
    }

    // 获取当前用户信息
    const userInfo = wx.getStorageSync('userInfo');
    if (!userInfo || !userInfo.openId) {
      wx.showToast({
        title: '请先登录',
        icon: 'none'
      });
      return;
    }
    
    // 构造记录数据
    const recordData = {
      id: Date.now().toString(),
      openId: userInfo.openId,
      title: title.trim(),
      content: content.trim(),
      grade,
      type,
      date,
      duration: parseInt(duration) || 0,
      completed,
      createTime: new Date().toISOString(),
      updateTime: new Date().toISOString()
    };

    // 从本地存储获取现有记录
    let records = wx.getStorageSync('climbing_records') || [];
    
    // 添加新记录
    records.unshift(recordData);
    
    // 保存到本地存储
    wx.setStorageSync('climbing_records', records);

    // 显示成功提示
    wx.showToast({
      title: '攀岩记录保存成功',
      icon: 'success'
    });

    // 返回上一页
    setTimeout(() => {
      wx.navigateBack();
    }, 1500);
  }
})