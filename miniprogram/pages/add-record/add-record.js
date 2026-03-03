Page({
  data: {
    title: '',
    content: '',
    category: '日常',
    categories: ['日常', '工作', '学习', '健康', '财务', '其他'],
    date: new Date().toISOString().split('T')[0],
    time: new Date().toTimeString().substring(0, 5),
    priority: 1
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

  onCategoryChange(e) {
    this.setData({
      category: this.data.categories[e.detail.value]
    });
  },

  onDateChange(e) {
    this.setData({
      date: e.detail.value
    });
  },

  onTimeChange(e) {
    this.setData({
      time: e.detail.value
    });
  },

  onPriorityChange(e) {
    this.setData({
      priority: parseInt(e.detail.value) + 1
    });
  },

  saveRecord() {
    const { title, content, category, date, time, priority } = this.data;
    
    if (!title.trim()) {
      wx.showToast({
        title: '请输入标题',
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
      category,
      date,
      time,
      priority,
      createTime: new Date().toISOString(),
      updateTime: new Date().toISOString()
    };

    // 从本地存储获取现有记录
    let records = wx.getStorageSync('records') || [];
    
    // 添加新记录
    records.unshift(recordData);
    
    // 保存到本地存储
    wx.setStorageSync('records', records);

    // 显示成功提示
    wx.showToast({
      title: '记录保存成功',
      icon: 'success'
    });

    // 返回上一页
    setTimeout(() => {
      wx.navigateBack();
    }, 1500);
  }
})