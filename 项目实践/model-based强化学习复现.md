# model-based 强化学习调度器复现

复现论文：Horizontal and Vertical Scaling of Container-based Applications using Reinforcement Learning

具体代码：[github](https://github.com/wtysos11/ModelbasedScheduler)

## 算法思想

参考资料：[课件](http://cseweb.ucsd.edu/~gary/190-RL/Lecture_Ch_9.pdf)

核心思想：使用已有过程估算出状态转移概率矩阵P，将未知的马尔科夫过程转换成已知的马尔科夫过程，从而实现快速迭代。

复现的是full backup model-based算法，这里的full backup指的是consider everything that might happen，在中文版中指的是期望更新。

问题的关键是如何准确估算这个状态转移概率矩阵P。

实验环境：
* 微服务容器云环境，kubernetes平台，单服务（仅考虑单个Deployment的伸缩）
* 状态空间为：实例数、响应时间、CPU占用率
* 动作空间：三动作（+1/0/-1）、五动作(+2~-2)
* 代价函数计算为：，目标是让代价函数最小。

实现细节：
* SLA规定，上界为250ms，下界为
* 调度间隔为
* 容器实例数范围为
* 响应时间离散化方法
* CPU占用率离散化方法
* 

## 任务1

阈值法，作为benchmark和数据来源

## 任务2

model-based的想法

第一个做法：在线更新

0. （可选）使用阈值法数据更新/使用阈值法数据初始化Q表
1. 拿到当前的实例数，响应时间和CPU占用率，计算得到当前状态s
2. 使用epsilon-greedy选择动作a，epsilon=1/t
3. 记录状态s'，更新概率转移记录和代价记录
4. 每一次更新概率之后，重新计算一遍所有状态和所有动作

离线更新想法：

使用阈值法做法进行预加载，然后将数据导入

需要建立三个数组进行更新：
1. Q表，表项为state * action个，用来决定最佳的动作，可以使用numpy数组
2. 转移矩阵，state\*action\*state，用来决定在给定state和action后，转移到state的概率为多少。可以使用字典+list的方式进行
3. 对performance penalty的估算函数，给定状态来估算C。