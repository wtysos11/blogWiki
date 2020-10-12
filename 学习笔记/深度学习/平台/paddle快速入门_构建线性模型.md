# paddle快速入门记录_构建线性模型

之前没有使用过百度的飞桨，因为课程要求使用了百度的平台，因此根据百度官方的[波士顿房价](https://www.paddlepaddle.org.cn/tutorials/projectdetail/597188)来学习线性模型的构建。

我快速浏览了一下，与pytorch比较接近，应该比较快能搞定。

## 流程说明

从任务的角度来说是数据处理->模型设计->训练配置->训练过程->模型保存。但是对我们来说，只有模型设计与训练过程是需要关注的，毕竟其他的也不是什么新鲜的东西。

### 基础构造

paddlepaddle的主库使用的是`paddle.fluid`，而且还分为动态图的类库和静态图的类库。从效果上来说，静态图可能会更快一些，但是预先定义网络结构。而动态图会更灵活一些。因此本次先参考教程，采用动态图构建网络(`paddle.fluid.dygraph`)。

### 网络设计

在paddle模型的设计设计与pytorch类似，需要诶网络单独实现一个类（该类需要继承`fluid.dygraph.Layer`）。类中要重写两个方法，一个是`__init__`，可以初始化整个网络，读取一些参数；一个是`forward`，即前向计算过程，负责网络的具体构建。

```python
class Regressor(fluid.dygraph.Layer):
    def __init__(self):
        super(Regressor, self).__init__()
        
        # 定义一层全连接层，输出维度是1，激活函数为None，即不使用激活函数
        self.fc = Linear(input_dim=13, output_dim=1, act=None)
    
    # 网络的前向计算函数
    def forward(self, inputs):
        x = self.fc(inputs)
        return x
```
### 训练配置

包含四步：
1. 指定运行训练的机器资源
2. 声明模型实例
3. 加载训练和测试数据
4. 设置优化算法和学习率

训练配置代码如下，指明了配置的资源、训练的模型、输入的数据、采用的优化算法与参数，与pytorch十分类似。

```python
# 定义飞桨动态图的工作环境
with fluid.dygraph.guard():
    # 声明定义好的线性回归模型
    model = Regressor()
    # 开启模型训练模式
    model.train()
    # 加载数据
    training_data, test_data = load_data()
    # 定义优化算法，这里使用随机梯度下降-SGD
    # 学习率设置为0.01
    opt = fluid.optimizer.SGD(learning_rate=0.01, parameter_list=model.parameters())
```

#### 指定机器资源

paddlepaddle使用`guard`函数指定训练的机器资源，表明在`with`作用域下的程序均执行在本机的CPU资源下。

### 训练过程

训练一般采用双层：
* 内层采用mini-batch模式，`for iter_id, mini_batch in enumerate(mini_batches):`
* 外层定义遍历数据集的次数。

batch的设置要适量，过大会消耗内存、过小会容易陷入局部最优解，每个batch的样本没有统计学意义。

1. 数据准备：将一个批次的数据转变成指定格式。
2. 前向计算：将一个批次的数据输入到网络中，计算输出结果。
3. 计算损失函数：损失函数使用飞桨提供的`square_error_cost`来计算
4. 反向传播：执行梯度反向传播`backward`函数，并根据设置的优化算法更新参数`opt.minimize`

```python
with dygraph.guard(fluid.CPUPlace()):
    EPOCH_NUM = 10   # 设置外层循环次数
    BATCH_SIZE = 10  # 设置batch大小
    
    # 定义外层循环
    for epoch_id in range(EPOCH_NUM):
        # 在每轮迭代开始之前，将训练数据的顺序随机的打乱
        np.random.shuffle(training_data)
        # 将训练数据进行拆分，每个batch包含10条数据
        mini_batches = [training_data[k:k+BATCH_SIZE] for k in range(0, len(training_data), BATCH_SIZE)]
        # 定义内层循环
        for iter_id, mini_batch in enumerate(mini_batches):
            x = np.array(mini_batch[:, :-1]).astype('float32') # 获得当前批次训练数据
            y = np.array(mini_batch[:, -1:]).astype('float32') # 获得当前批次训练标签（真实房价）
            # 将numpy数据转为飞桨动态图variable形式
            house_features = dygraph.to_variable(x)
            prices = dygraph.to_variable(y)
            
            # 前向计算
            predicts = model(house_features)
            
            # 计算损失
            loss = fluid.layers.square_error_cost(predicts, label=prices)
            avg_loss = fluid.layers.mean(loss)
            if iter_id%20==0:
                print("epoch: {}, iter: {}, loss is: {}".format(epoch_id, iter_id, avg_loss.numpy()))
            
            # 反向传播
            avg_loss.backward()
            # 最小化loss,更新参数
            opt.minimize(avg_loss)
            # 清除梯度
            model.clear_gradients()
    # 保存模型
    fluid.save_dygraph(model.state_dict(), 'LR_model')
```


