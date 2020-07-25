# Keras注意力机制详解与实现

标签：视频学习 B站 短篇

[视频源](https://www.bilibili.com/video/BV1C7411k7Wg?from=search&seid=4206239689304403160)

对应的文章：[csdn](https://blog.csdn.net/weixin_44791964/article/details/104000722)

## 实验

设计了一个数据，输出只与timestep为2时的两个特征相关，两个值为0的时候输出为0，两个值为1的时候输出为1。

因此，如果要更好地学习，则模型应该更专注于timestep=2时的数据。

