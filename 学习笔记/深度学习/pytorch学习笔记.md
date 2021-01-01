# pytorch学习笔记

参考资料：
* [莫凡](https://mofanpy.com/tutorials/machine-learning/torch/)，优点是与视频结合非常方便理解
* [中文-gitbook](https://github.com/zergtant/pytorch-handbook)，很好的入门资料


## Pytorch基础-快速入门

### pytorch与numpy

#### 相互转换

numpy是非常常用的数据保存方式，pytorch可以用以下的方式进行转换

```python
import torch
import numpy as np

np_data = np.arange(6).reshape((2, 3))
torch_data = torch.from_numpy(np_data)
tensor2array = torch_data.numpy()
print(
    '\nnumpy array:', np_data,          # [[0 1 2], [3 4 5]]
    '\ntorch tensor:', torch_data,      #  0  1  2 \n 3  4  5    [torch.LongTensor of size 2x3]
    '\ntensor to array:', tensor2array, # [[0 1 2], [3 4 5]]
)
```

#### Torch数学运算

Torch中的数学运算与numpy是一样的。更多的话应该参考[API doc](https://pytorch.org/docs/stable/torch.html)

* 矩阵乘法：`np.matmul(a,b)`,`torch.mm(a,b)`

常用操作：
* `torch.unsqueeze`
* `torch.squeeze`
* `torch.gather`

### 变量

pytorch与tensorflow类似，计算梯度的时候可以直接根据变量之间的操作来进行计算。

```python
import torch
from torch.autograd import Variable # torch 中 Variable 模块

# 先生鸡蛋
tensor = torch.FloatTensor([[1,2],[3,4]])
# 把鸡蛋放到篮子里, requires_grad是参不参与误差反向传播, 要不要计算梯度
variable = Variable(tensor, requires_grad=True)

print(tensor)
"""
 1  2
 3  4
[torch.FloatTensor of size 2x2]
"""

print(variable)
"""
Variable containing:
 1  2
 3  4
[torch.FloatTensor of size 2x2]
"""
```

所有torch数学操作返回的即是变量，变量完成了计算图的搭建，在最终误差反传的时候就能够计算出梯度

```python
t_out = torch.mean(tensor*tensor)       # x^2
v_out = torch.mean(variable*variable)   # x^2
print(t_out)
print(v_out)    # 7.5
v_out.backward()    # 模拟 v_out 的误差反向传递

# 下面两步看不懂没关系, 只要知道 Variable 是计算图的一部分, 可以用来传递误差就好.
# v_out = 1/4 * sum(variable*variable) 这是计算图中的 v_out 计算步骤
# 针对于 v_out 的梯度就是, d(v_out)/d(variable) = 1/4*2*variable = variable/2

print(variable.grad)    # 初始 Variable 的梯度
'''
 0.5000  1.0000
 1.5000  2.0000
'''
```

如果要利用变量内的数据，应当使用
* `variable.data`，返回tensor形式数据
* `variable.numpy()`，返回numpy形式数据

### 网络示例

自然，有一个网络用一下很多东西很快就明白了

pytorch与Keras比较大的不同就在于新建网络的时候要从`torch.nn.Module`来继承。`init`来完成网络的构建，并用`forward`来描述数据的操作。

```python
import torch
import torch.nn.functional as F     # 激励函数都在这

class Net(torch.nn.Module):  # 继承 torch 的 Module
    def __init__(self, n_feature, n_hidden, n_output):
        super(Net, self).__init__()     # 继承 __init__ 功能
        # 定义每层用什么样的形式
        self.hidden = torch.nn.Linear(n_feature, n_hidden)   # 隐藏层线性输出
        self.predict = torch.nn.Linear(n_hidden, n_output)   # 输出层线性输出

    def forward(self, x):   # 这同时也是 Module 中的 forward 功能
        # 正向传播输入值, 神经网络分析出输出值
        x = F.relu(self.hidden(x))      # 激励函数(隐藏层的线性值)
        x = self.predict(x)             # 输出值
        return x

net = Net(n_feature=1, n_hidden=10, n_output=1)

print(net)  # net 的结构
"""
Net (
  (hidden): Linear (1 -> 10)
  (predict): Linear (10 -> 1)
)
"""
```

`optimizer`和`loss_func`自然不必多说，训练的过程就是拿到标签`val_target`和现在的已有值`val_eval`，使用`optimizer.zero_grad()`清空梯度，然后就可以使用`loss = loss_func(val_target,val_eval)`计算出损失函数的值，并用`loss.backward()`反传，然后用`optimizer.step()`来应用梯度下降器。

可以说非常方便了。

```python
# optimizer 是训练的工具
optimizer = torch.optim.SGD(net.parameters(), lr=0.2)  # 传入 net 的所有参数, 学习率
loss_func = torch.nn.MSELoss()      # 预测值和真实值的误差计算公式 (均方差)

for t in range(100):
    prediction = net(x)     # 喂给 net 训练数据 x, 输出预测值

    loss = loss_func(prediction, y)     # 计算两者的误差

    optimizer.zero_grad()   # 清空上一步的残余更新参数值
    loss.backward()         # 误差反向传播, 计算参数更新值
    optimizer.step()        # 将参数更新值施加到 net 的 parameters 上
```

```python
import matplotlib.pyplot as plt

plt.ion()   # 画图
plt.show()

for t in range(200):

    ...
    loss.backward()
    optimizer.step()

    # 接着上面来
    if t % 5 == 0:
        # plot and show learning process
        plt.cla()
        plt.scatter(x.data.numpy(), y.data.numpy())
        plt.plot(x.data.numpy(), prediction.data.numpy(), 'r-', lw=5)
        plt.text(0.5, 0, 'Loss=%.4f' % loss.data.numpy(), fontdict={'size': 20, 'color':  'red'})
        plt.pause(0.1)
```