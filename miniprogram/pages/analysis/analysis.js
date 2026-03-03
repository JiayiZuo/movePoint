Page({
  data: {
    statistics: {
      total: 0,
      today: 0,
      categories: {},
      priorities: {}
    },
    chartData: []
  },

  onLoad() {
    this.loadAndAnalyzeRecords();
  },

  loadAndAnalyzeRecords() {
    // 从本地存储获取记录
    const records = wx.getStorageSync('records') || [];
    const userInfo = wx.getStorageSync('userInfo');
    
    // 如果未登录，不显示提示，只设置默认值
    if (!userInfo || !userInfo.openId) {
      this.setData({
        statistics: {
          total: 0,
          today: 0,
          categories: {},
          priorities: {}
        },
        chartData: []
      });
      return;
    }
    
    // 过滤当前用户的记录
    const userRecords = records.filter(record => record.openId === userInfo.openId);
    
    // 统计数据
    const stats = this.analyzeRecords(userRecords);
    
    this.setData({
      statistics: stats,
      chartData: this.generateChartData(stats)
    });
  },

  analyzeRecords(records) {
    const stats = {
      total: records.length,
      today: 0,
      categories: {},
      priorities: {}
    };
    
    const today = new Date().toISOString().split('T')[0];
    
    records.forEach(record => {
      // 统计今日记录数
      if (record.date === today) {
        stats.today++;
      }
      
      // 统计分类
      if (stats.categories[record.category]) {
        stats.categories[record.category]++;
      } else {
        stats.categories[record.category] = 1;
      }
      
      // 统计优先级
      const priorityKey = `P${record.priority}`;
      if (stats.priorities[priorityKey]) {
        stats.priorities[priorityKey]++;
      } else {
        stats.priorities[priorityKey] = 1;
      }
    });
    
    return stats;
  },

  generateChartData(stats) {
    // 生成图表数据，这里简单返回统计结果
    const categoryData = Object.keys(stats.categories).map(key => ({
      name: key,
      count: stats.categories[key]
    }));
    
    const priorityData = Object.keys(stats.priorities).map(key => ({
      name: key,
      count: stats.priorities[key]
    }));
    
    return [
      { type: 'category', title: '分类统计', data: categoryData },
      { type: 'priority', title: '优先级统计', data: priorityData }
    ];
  },

  onRefresh() {
    this.loadAndAnalyzeRecords();
  }
})