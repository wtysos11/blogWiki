# LSTM Multi-step forward

[文章](https://machinelearningmastery.com/how-to-develop-lstm-models-for-multi-step-time-series-forecasting-of-household-power-consumption/)

[ConvLSTM](https://arxiv.org/abs/1506.04214v1)


## model evaluation

我们如何开发，已经评价一个模型

### Problem Framing

Given recent power consumption, what is the expected power consumption for the week ahead?

### Evaluation Metric

选择使用RMSE，其中为了判断究竟预测几步会比较好，对各个预测步分别计算RMSE，并在最后加权求一个总和

```python
# evaluate one or more weekly forecasts against expected values
# 对未来一周的预测数据进行估算
def evaluate_forecasts(actual, predicted):
	scores = list()
	# calculate an RMSE score for each day
	for i in range(actual.shape[1]):
		# calculate mse
		mse = mean_squared_error(actual[:, i], predicted[:, i])
		# calculate rmse
		rmse = sqrt(mse)
		# store
		scores.append(rmse)
	# calculate overall RMSE
	s = 0
	for row in range(actual.shape[0]):
		for col in range(actual.shape[1]):
			s += (actual[row, col] - predicted[row, col])**2
	score = sqrt(s / (actual.shape[0] * actual.shape[1]))
	return score, scores
```

### Train and Test Sets

使用头三年的数据来训练预测模型，并用最后一年的数据来评估模型。

数据集中的数据会被切分成标准星期，即从周日开始，周六结束的数据。这样可以预测一周前的数据，并且可以预测特定某一天的数据。

将这些数据组装成159个完整的standard weeks作为训练集，用46周的数据作为测试集合。

会使用一种名为walk-forward validation的方式来指导预测。模型会需要来进行下一周数据的预测，然后下一周的数据就会在预测后变得可用，用来预测再下一周的数据。

### LSTM for Multi-Step Forecasting

LSTM的优点：1.天然支持序列预测。2.允许多特征输入 3.允许向量输出
使用seq2seq来进行多对多预测。

seq2seq模型将原始模型分为了两个部分：
* encoder:将输入序列压缩成一个定长的内部表示
* decoder:将定长的内部表示用来预测外部序列

一般的，LSTM不是非常适合用于auto-regression类型的问题上，参考[On the Suitability of LSTMs for Time Series Forecasting](https://machinelearningmastery.com/suitability-long-short-term-memory-networks-time-series-forecasting/)，此时可以使用一维的卷积神经网络来提取特征，可能会更好。即CNN-LSTM模型，参考[这里](https://machinelearningmastery.com/cnn-long-short-term-memory-networks/)。其中TCN模型是其中的一种。

在本教程中，我们会探索以下内容：

* LSTM模型，处理多步预测，单输入数据
* seq2seq的LSTM模型，处理多步预测，单输入数据
* seq2seq的LSTM模型，处理多步预测，多输入数据
* CNN-LSTM Encoder-Decoder模型，多输入多步预测
* ConvLSTM Encoder-Decoder模型，多输入多步预测

[LSTM时间预测-入门](https://machinelearningmastery.com/how-to-develop-lstm-models-for-time-series-forecasting/)

### 单输入多输出的LSTM

首先以一个简单的LSTM模型开始，读取一系列以日为单位的日能源消耗，并预测下一个标准周的日能源消耗。

输入的数据为一个一维数组，输入数据包括如下数据：
* 所有之前的天粒度数据
* 前七天数据
* 前两个星期数据
* 前一个月数据
* 前一年数据

并没有一个泛用的方法，每种方法都是可以的。

一个LSTM模型希望数据有如下的格式：`[samples,timesteps,features]`。
在本文中，一个训练数据会有7个时间步与一个特征点，因此训练集的数据为`[159,7,1]`

这是一个好的开始，用之前的数据来训练之后的数据是可行的。问题是159个实例对于训练一个神经网络而言并不是很多。

一种增加更多训练数据的方式是改变问题定义，每七天预测后七天的数据，忽视调标准周的规定。然而在测试集合中，问题定义并没有改变，仍然是对每一个标准周预测下一个标准周的能源消耗。

训练数据格式如下：`[159,7,8]`，即159个标准周，每周7个时间点，每个时间点8个特征。

处理：
1. 首先flatten the data：`data = train.reshape((train.shape[0]*train.shape[1], train.shape[2]))`
2. 然后遍历时间点，对数据进行切分。期望切分情况如下：
```
Input, Output
[d01, d02, d03, d04, d05, d06, d07], [d08, d09, d10, d11, d12, d13, d14]
[d02, d03, d04, d05, d06, d07, d08], [d09, d10, d11, d12, d13, d14, d15]
...
```

这样就构成了函数`to_supervised()`的逻辑，接收一个以周为单位的数据列表，Input和output的timestep作为参数，计算出有交叉部分的数据，用作神经网络的训练部分。

这样，就可以将159个样本转换为1100个样本。

下一步，我们定义处理训练数据所用的LSTM神经网络。

这个多步时间序列预测问题是一个自回归问题，这意味着需要建立一个前七天到预测七天数据之间的一个模型。LSTM层后面会连接一个200个结点的全连接层，来解释LSTM学习到的特征。最终，会导入到一个输出层中，直接预测出后七天所需要的结果向量。同时，我们使用RMSE作为损失评估函数，使用Adam作为SGD方法，70次训练，batch-size为16。

小的batch-size和算法的随机梯度下降背景意味着同样的模型每次输入同样的输入会有着略微不同的输出。即结果会在模型被评估的时候会有所变化，可以尝试多次训练模型，计算一个平均的表现。

下面的`build_model()`方法准备训练数据，定义LSTM模型，并用训练数据去训练定义的LSTM模型，将训练好的模型返回。

得到模型之后，下一步我们来看一下该如何使用这个模型来进行预测

一般来说，LSTM模型期望一个三维的输入。在本例子中，这个输入应该是一个样本，七天，每天一个特征值的输入，即`[1,7,1]`

因为我们使用了walk-forward validation，即每次使用过去七天来预测未来七天。对新的七天，使用真实数据而非预测数据。因此每次多预测七天后，我们就多了七天的历史数据。

为了去预测下一个七天的数据，我们需要取得已经知道数据中最后的七天。下面的`forecast()`函数实现了这个，依据模型和历史数据来进行预测。

univariate LSTM，全部代码（调整输入从7->14可以较大的提高最终的结果）

```python
# univariate multi-step lstm
from math import sqrt
from numpy import split
from numpy import array
from pandas import read_csv
from sklearn.metrics import mean_squared_error
from matplotlib import pyplot
from keras.models import Sequential
from keras.layers import Dense
from keras.layers import Flatten
from keras.layers import LSTM

# split a univariate dataset into train/test sets
def split_dataset(data):
	# split into standard weeks
	train, test = data[1:-328], data[-328:-6]
	# restructure into windows of weekly data
	train = array(split(train, len(train)/7))
	test = array(split(test, len(test)/7))
	return train, test

# evaluate one or more weekly forecasts against expected values
def evaluate_forecasts(actual, predicted):
	scores = list()
	# calculate an RMSE score for each day
	for i in range(actual.shape[1]):
		# calculate mse
		mse = mean_squared_error(actual[:, i], predicted[:, i])
		# calculate rmse
		rmse = sqrt(mse)
		# store
		scores.append(rmse)
	# calculate overall RMSE
	s = 0
	for row in range(actual.shape[0]):
		for col in range(actual.shape[1]):
			s += (actual[row, col] - predicted[row, col])**2
	score = sqrt(s / (actual.shape[0] * actual.shape[1]))
	return score, scores

# summarize scores
def summarize_scores(name, score, scores):
	s_scores = ', '.join(['%.1f' % s for s in scores])
	print('%s: [%.3f] %s' % (name, score, s_scores))

# convert history into inputs and outputs
def to_supervised(train, n_input, n_out=7):
	# flatten data
	data = train.reshape((train.shape[0]*train.shape[1], train.shape[2]))
	X, y = list(), list()
	in_start = 0
	# step over the entire history one time step at a time
	for _ in range(len(data)):
		# define the end of the input sequence
		in_end = in_start + n_input
		out_end = in_end + n_out
		# ensure we have enough data for this instance
		if out_end <= len(data):
			x_input = data[in_start:in_end, 0]
			x_input = x_input.reshape((len(x_input), 1))
			X.append(x_input)
			y.append(data[in_end:out_end, 0])
		# move along one time step
		in_start += 1
	return array(X), array(y)

# train the model
def build_model(train, n_input):
	# prepare data
	train_x, train_y = to_supervised(train, n_input)
	# define parameters
	verbose, epochs, batch_size = 0, 70, 16
	n_timesteps, n_features, n_outputs = train_x.shape[1], train_x.shape[2], train_y.shape[1]
	# define model
	model = Sequential()
	model.add(LSTM(200, activation='relu', input_shape=(n_timesteps, n_features)))
	model.add(Dense(100, activation='relu'))
	model.add(Dense(n_outputs))
	model.compile(loss='mse', optimizer='adam')
	# fit network
	model.fit(train_x, train_y, epochs=epochs, batch_size=batch_size, verbose=verbose)
	return model

# make a forecast
def forecast(model, history, n_input):
	# flatten data
	data = array(history)
	data = data.reshape((data.shape[0]*data.shape[1], data.shape[2]))
	# retrieve last observations for input data
	input_x = data[-n_input:, 0]
	# reshape into [1, n_input, 1]
	input_x = input_x.reshape((1, len(input_x), 1))
	# forecast the next week
	yhat = model.predict(input_x, verbose=0)
	# we only want the vector forecast
	yhat = yhat[0]
	return yhat

# evaluate a single model
def evaluate_model(train, test, n_input):
	# fit model
	model = build_model(train, n_input)
	# history is a list of weekly data
	history = [x for x in train]
	# walk-forward validation over each week
	predictions = list()
	for i in range(len(test)):
		# predict the week
		yhat_sequence = forecast(model, history, n_input)
		# store the predictions
		predictions.append(yhat_sequence)
		# get real observation and add to history for predicting the next week
		history.append(test[i, :])
	# evaluate predictions days for each week
	predictions = array(predictions)
	score, scores = evaluate_forecasts(test[:, :, 0], predictions)
	return score, scores

# load the new file
dataset = read_csv('household_power_consumption_days.csv', header=0, infer_datetime_format=True, parse_dates=['datetime'], index_col=['datetime'])
# split into train and test
train, test = split_dataset(dataset.values)
# evaluate model and get scores
n_input = 7
score, scores = evaluate_model(train, test, n_input)
# summarize scores
summarize_scores('lstm', score, scores)
# plot scores
days = ['sun', 'mon', 'tue', 'wed', 'thr', 'fri', 'sat']
pyplot.plot(days, scores, marker='o', label='lstm')
pyplot.show()
```

### Encoder-Decoder LSTM Model With Univariate Input

update the vanilla LSTM to use an encoder-decoder model.

更新：模型不再直接产生一个vector sequence，而是分成了两个子模型，编码器(encoder)读取并编码输入序列，解码器(decoder)读取编码的输入序列，并对每一个输出序列进行一步的预测。

比较重要的区别是，LSTM模型被用在了decoder中。我们会实现一个与LSTM autoencoder类似的encoder-decoder结构。
定义模型：
1. 首先，输入序列的内部表示会被多次重复，每一个time step都要有一个输出序列，可以这样做`model.add(RepeatVector(7))`
2. 然后定义一个LSTM隐层作为解码器。`model.add(LSTM(200,activaion='relu',return_sequences=True))`
3. 然后使用一个全连接层来在每一个time step解释。这个输出层每次只预测一步，而不是直接预测七天。`model.add(TimeDistributed(Dense(100,activation='relu'))) model.add(TimeDistributed(Dense(1)))`。这个wrapper可以让Dense层分别在每一个时间序列上工作。

这个网络的输出是一个三维的结构`[samples,timesteps,features]`，因此一个一星期的预测结构如下：`[1,7,1]`

使用了完整的Encoder-Decoder的LSTM如下：

```python
# univariate multi-step encoder-decoder lstm
from math import sqrt
from numpy import split
from numpy import array
from pandas import read_csv
from sklearn.metrics import mean_squared_error
from matplotlib import pyplot
from keras.models import Sequential
from keras.layers import Dense
from keras.layers import Flatten
from keras.layers import LSTM
from keras.layers import RepeatVector
from keras.layers import TimeDistributed

# split a univariate dataset into train/test sets
def split_dataset(data):
	# split into standard weeks
	train, test = data[1:-328], data[-328:-6]
	# restructure into windows of weekly data
	train = array(split(train, len(train)/7))
	test = array(split(test, len(test)/7))
	return train, test

# evaluate one or more weekly forecasts against expected values
def evaluate_forecasts(actual, predicted):
	scores = list()
	# calculate an RMSE score for each day
	for i in range(actual.shape[1]):
		# calculate mse
		mse = mean_squared_error(actual[:, i], predicted[:, i])
		# calculate rmse
		rmse = sqrt(mse)
		# store
		scores.append(rmse)
	# calculate overall RMSE
	s = 0
	for row in range(actual.shape[0]):
		for col in range(actual.shape[1]):
			s += (actual[row, col] - predicted[row, col])**2
	score = sqrt(s / (actual.shape[0] * actual.shape[1]))
	return score, scores

# summarize scores
def summarize_scores(name, score, scores):
	s_scores = ', '.join(['%.1f' % s for s in scores])
	print('%s: [%.3f] %s' % (name, score, s_scores))

# convert history into inputs and outputs
def to_supervised(train, n_input, n_out=7):
	# flatten data
	data = train.reshape((train.shape[0]*train.shape[1], train.shape[2]))
	X, y = list(), list()
	in_start = 0
	# step over the entire history one time step at a time
	for _ in range(len(data)):
		# define the end of the input sequence
		in_end = in_start + n_input
		out_end = in_end + n_out
		# ensure we have enough data for this instance
		if out_end <= len(data):
			x_input = data[in_start:in_end, 0]
			x_input = x_input.reshape((len(x_input), 1))
			X.append(x_input)
			y.append(data[in_end:out_end, 0])
		# move along one time step
		in_start += 1
	return array(X), array(y)

# train the model
def build_model(train, n_input):
	# prepare data
	train_x, train_y = to_supervised(train, n_input)
	# define parameters
	verbose, epochs, batch_size = 0, 20, 16
	n_timesteps, n_features, n_outputs = train_x.shape[1], train_x.shape[2], train_y.shape[1]
	# reshape output into [samples, timesteps, features]
	train_y = train_y.reshape((train_y.shape[0], train_y.shape[1], 1))
	# define model
	model = Sequential()
	model.add(LSTM(200, activation='relu', input_shape=(n_timesteps, n_features)))
	model.add(RepeatVector(n_outputs))
	model.add(LSTM(200, activation='relu', return_sequences=True))
	model.add(TimeDistributed(Dense(100, activation='relu')))
	model.add(TimeDistributed(Dense(1)))
	model.compile(loss='mse', optimizer='adam')
	# fit network
	model.fit(train_x, train_y, epochs=epochs, batch_size=batch_size, verbose=verbose)
	return model

# make a forecast
def forecast(model, history, n_input):
	# flatten data
	data = array(history)
	data = data.reshape((data.shape[0]*data.shape[1], data.shape[2]))
	# retrieve last observations for input data
	input_x = data[-n_input:, 0]
	# reshape into [1, n_input, 1]
	input_x = input_x.reshape((1, len(input_x), 1))
	# forecast the next week
	yhat = model.predict(input_x, verbose=0)
	# we only want the vector forecast
	yhat = yhat[0]
	return yhat

# evaluate a single model
def evaluate_model(train, test, n_input):
	# fit model
	model = build_model(train, n_input)
	# history is a list of weekly data
	history = [x for x in train]
	# walk-forward validation over each week
	predictions = list()
	for i in range(len(test)):
		# predict the week
		yhat_sequence = forecast(model, history, n_input)
		# store the predictions
		predictions.append(yhat_sequence)
		# get real observation and add to history for predicting the next week
		history.append(test[i, :])
	# evaluate predictions days for each week
	predictions = array(predictions)
	score, scores = evaluate_forecasts(test[:, :, 0], predictions)
	return score, scores

# load the new file
dataset = read_csv('household_power_consumption_days.csv', header=0, infer_datetime_format=True, parse_dates=['datetime'], index_col=['datetime'])
# split into train and test
train, test = split_dataset(dataset.values)
# evaluate model and get scores
n_input = 14
score, scores = evaluate_model(train, test, n_input)
# summarize scores
summarize_scores('lstm', score, scores)
# plot scores
days = ['sun', 'mon', 'tue', 'wed', 'thr', 'fri', 'sat']
pyplot.plot(days, scores, marker='o', label='lstm')
pyplot.show()
```

### Encoder-Decoder LSTM Model With Multivariate Input

会更新前面所提到的模型，每次使用8个特征来预测未来一个星期的日总能源消耗量。

完整代码：（其实没有什么太大的区别）

```python
# multivariate multi-step encoder-decoder lstm
from math import sqrt
from numpy import split
from numpy import array
from pandas import read_csv
from sklearn.metrics import mean_squared_error
from matplotlib import pyplot
from keras.models import Sequential
from keras.layers import Dense
from keras.layers import Flatten
from keras.layers import LSTM
from keras.layers import RepeatVector
from keras.layers import TimeDistributed

# split a univariate dataset into train/test sets
def split_dataset(data):
	# split into standard weeks
	train, test = data[1:-328], data[-328:-6]
	# restructure into windows of weekly data
	train = array(split(train, len(train)/7))
	test = array(split(test, len(test)/7))
	return train, test

# evaluate one or more weekly forecasts against expected values
def evaluate_forecasts(actual, predicted):
	scores = list()
	# calculate an RMSE score for each day
	for i in range(actual.shape[1]):
		# calculate mse
		mse = mean_squared_error(actual[:, i], predicted[:, i])
		# calculate rmse
		rmse = sqrt(mse)
		# store
		scores.append(rmse)
	# calculate overall RMSE
	s = 0
	for row in range(actual.shape[0]):
		for col in range(actual.shape[1]):
			s += (actual[row, col] - predicted[row, col])**2
	score = sqrt(s / (actual.shape[0] * actual.shape[1]))
	return score, scores

# summarize scores
def summarize_scores(name, score, scores):
	s_scores = ', '.join(['%.1f' % s for s in scores])
	print('%s: [%.3f] %s' % (name, score, s_scores))

# convert history into inputs and outputs
def to_supervised(train, n_input, n_out=7):
	# flatten data
	data = train.reshape((train.shape[0]*train.shape[1], train.shape[2]))
	X, y = list(), list()
	in_start = 0
	# step over the entire history one time step at a time
	for _ in range(len(data)):
		# define the end of the input sequence
		in_end = in_start + n_input
		out_end = in_end + n_out
		# ensure we have enough data for this instance
		if out_end <= len(data):
			X.append(data[in_start:in_end, :])
			y.append(data[in_end:out_end, 0])
		# move along one time step
		in_start += 1
	return array(X), array(y)

# train the model
def build_model(train, n_input):
	# prepare data
	train_x, train_y = to_supervised(train, n_input)
	# define parameters
	verbose, epochs, batch_size = 0, 50, 16
	n_timesteps, n_features, n_outputs = train_x.shape[1], train_x.shape[2], train_y.shape[1]
	# reshape output into [samples, timesteps, features]
	train_y = train_y.reshape((train_y.shape[0], train_y.shape[1], 1))
	# define model
	model = Sequential()
	model.add(LSTM(200, activation='relu', input_shape=(n_timesteps, n_features)))
	model.add(RepeatVector(n_outputs))
	model.add(LSTM(200, activation='relu', return_sequences=True))
	model.add(TimeDistributed(Dense(100, activation='relu')))
	model.add(TimeDistributed(Dense(1)))
	model.compile(loss='mse', optimizer='adam')
	# fit network
	model.fit(train_x, train_y, epochs=epochs, batch_size=batch_size, verbose=verbose)
	return model

# make a forecast
def forecast(model, history, n_input):
	# flatten data
	data = array(history)
	data = data.reshape((data.shape[0]*data.shape[1], data.shape[2]))
	# retrieve last observations for input data
	input_x = data[-n_input:, :]
	# reshape into [1, n_input, n]
	input_x = input_x.reshape((1, input_x.shape[0], input_x.shape[1]))
	# forecast the next week
	yhat = model.predict(input_x, verbose=0)
	# we only want the vector forecast
	yhat = yhat[0]
	return yhat

# evaluate a single model
def evaluate_model(train, test, n_input):
	# fit model
	model = build_model(train, n_input)
	# history is a list of weekly data
	history = [x for x in train]
	# walk-forward validation over each week
	predictions = list()
	for i in range(len(test)):
		# predict the week
		yhat_sequence = forecast(model, history, n_input)
		# store the predictions
		predictions.append(yhat_sequence)
		# get real observation and add to history for predicting the next week
		history.append(test[i, :])
	# evaluate predictions days for each week
	predictions = array(predictions)
	score, scores = evaluate_forecasts(test[:, :, 0], predictions)
	return score, scores

# load the new file
dataset = read_csv('household_power_consumption_days.csv', header=0, infer_datetime_format=True, parse_dates=['datetime'], index_col=['datetime'])
# split into train and test
train, test = split_dataset(dataset.values)
# evaluate model and get scores
n_input = 14
score, scores = evaluate_model(train, test, n_input)
# summarize scores
summarize_scores('lstm', score, scores)
# plot scores
days = ['sun', 'mon', 'tue', 'wed', 'thr', 'fri', 'sat']
pyplot.plot(days, scores, marker='o', label='lstm')
pyplot.show()
```

### CNN-LSTM Encoder-Decoder Model With Univariate Input

卷积神经网络可以用来作为编码-解码器中的编码器。CNN本身并不支持序列输入，而1维CNN可以读取序列输入并自动学习salient features，这个特征结果可以被LSTM解码器解释。

CNN期望的输入数据与LSTM模型相同。为了简化模型，我们使用单输入，虽然多输入其实也挺简单的。

与之前一样，我们使用14天的天总能量消耗来进行预测。

CNN的结构由两个卷积层构成，最后有一个最大池化层，结果会被平摊。第一个卷积层会读取整个输入序列，并将结果映射到特征空间中。第二个卷积层会在第一层创建出的特征空间里执行同样的操作，尝试去放大任何salient features。我们在每个卷积层使用64个特征，使用一个长度为3的核来读取输入序列。

最大池化层简化了feature map，通过将1/4的数据用最大信号代替。然后这个会被平摊开。

```python
model.add(Conv1D(filters=64, kernel_size=3, activation='relu', input_shape=(n_timesteps,n_features)))
model.add(Conv1D(filters=64, kernel_size=3, activation='relu'))
model.add(MaxPooling1D(pool_size=2))
model.add(Flatten())
```

解码器部分与上面一样，唯一的区别是训练周期变为20.

完整代码：
```python
# univariate multi-step encoder-decoder cnn-lstm
from math import sqrt
from numpy import split
from numpy import array
from pandas import read_csv
from sklearn.metrics import mean_squared_error
from matplotlib import pyplot
from keras.models import Sequential
from keras.layers import Dense
from keras.layers import Flatten
from keras.layers import LSTM
from keras.layers import RepeatVector
from keras.layers import TimeDistributed
from keras.layers.convolutional import Conv1D
from keras.layers.convolutional import MaxPooling1D

# split a univariate dataset into train/test sets
def split_dataset(data):
	# split into standard weeks
	train, test = data[1:-328], data[-328:-6]
	# restructure into windows of weekly data
	train = array(split(train, len(train)/7))
	test = array(split(test, len(test)/7))
	return train, test

# evaluate one or more weekly forecasts against expected values
def evaluate_forecasts(actual, predicted):
	scores = list()
	# calculate an RMSE score for each day
	for i in range(actual.shape[1]):
		# calculate mse
		mse = mean_squared_error(actual[:, i], predicted[:, i])
		# calculate rmse
		rmse = sqrt(mse)
		# store
		scores.append(rmse)
	# calculate overall RMSE
	s = 0
	for row in range(actual.shape[0]):
		for col in range(actual.shape[1]):
			s += (actual[row, col] - predicted[row, col])**2
	score = sqrt(s / (actual.shape[0] * actual.shape[1]))
	return score, scores

# summarize scores
def summarize_scores(name, score, scores):
	s_scores = ', '.join(['%.1f' % s for s in scores])
	print('%s: [%.3f] %s' % (name, score, s_scores))

# convert history into inputs and outputs
def to_supervised(train, n_input, n_out=7):
	# flatten data
	data = train.reshape((train.shape[0]*train.shape[1], train.shape[2]))
	X, y = list(), list()
	in_start = 0
	# step over the entire history one time step at a time
	for _ in range(len(data)):
		# define the end of the input sequence
		in_end = in_start + n_input
		out_end = in_end + n_out
		# ensure we have enough data for this instance
		if out_end <= len(data):
			x_input = data[in_start:in_end, 0]
			x_input = x_input.reshape((len(x_input), 1))
			X.append(x_input)
			y.append(data[in_end:out_end, 0])
		# move along one time step
		in_start += 1
	return array(X), array(y)

# train the model
def build_model(train, n_input):
	# prepare data
	train_x, train_y = to_supervised(train, n_input)
	# define parameters
	verbose, epochs, batch_size = 0, 20, 16
	n_timesteps, n_features, n_outputs = train_x.shape[1], train_x.shape[2], train_y.shape[1]
	# reshape output into [samples, timesteps, features]
	train_y = train_y.reshape((train_y.shape[0], train_y.shape[1], 1))
	# define model
	model = Sequential()
	model.add(Conv1D(filters=64, kernel_size=3, activation='relu', input_shape=(n_timesteps,n_features)))
	model.add(Conv1D(filters=64, kernel_size=3, activation='relu'))
	model.add(MaxPooling1D(pool_size=2))
	model.add(Flatten())
	model.add(RepeatVector(n_outputs))
	model.add(LSTM(200, activation='relu', return_sequences=True))
	model.add(TimeDistributed(Dense(100, activation='relu')))
	model.add(TimeDistributed(Dense(1)))
	model.compile(loss='mse', optimizer='adam')
	# fit network
	model.fit(train_x, train_y, epochs=epochs, batch_size=batch_size, verbose=verbose)
	return model

# make a forecast
def forecast(model, history, n_input):
	# flatten data
	data = array(history)
	data = data.reshape((data.shape[0]*data.shape[1], data.shape[2]))
	# retrieve last observations for input data
	input_x = data[-n_input:, 0]
	# reshape into [1, n_input, 1]
	input_x = input_x.reshape((1, len(input_x), 1))
	# forecast the next week
	yhat = model.predict(input_x, verbose=0)
	# we only want the vector forecast
	yhat = yhat[0]
	return yhat

# evaluate a single model
def evaluate_model(train, test, n_input):
	# fit model
	model = build_model(train, n_input)
	# history is a list of weekly data
	history = [x for x in train]
	# walk-forward validation over each week
	predictions = list()
	for i in range(len(test)):
		# predict the week
		yhat_sequence = forecast(model, history, n_input)
		# store the predictions
		predictions.append(yhat_sequence)
		# get real observation and add to history for predicting the next week
		history.append(test[i, :])
	# evaluate predictions days for each week
	predictions = array(predictions)
	score, scores = evaluate_forecasts(test[:, :, 0], predictions)
	return score, scores

# load the new file
dataset = read_csv('household_power_consumption_days.csv', header=0, infer_datetime_format=True, parse_dates=['datetime'], index_col=['datetime'])
# split into train and test
train, test = split_dataset(dataset.values)
# evaluate model and get scores
n_input = 14
score, scores = evaluate_model(train, test, n_input)
# summarize scores
summarize_scores('lstm', score, scores)
# plot scores
days = ['sun', 'mon', 'tue', 'wed', 'thr', 'fri', 'sat']
pyplot.plot(days, scores, marker='o', label='lstm')
pyplot.show()
```

### ConvLSTM Encoder-Decoder Model With Univariate Input

一个更大的扩展是，将卷积操作用在每一个时间步的输入中。这个结合被称为Convolutional LSTM，或简称为ConvLSTM

不像LSTM直接读取数据并计算内部状态以及状态转移，也不像CNN-LSTM来解释CNN模型产生的的结果，ConvLSTM直接使用卷积层读取数据，读入到LSTM单元内部。

[参考论文](https://arxiv.org/abs/1506.04214v1)

Keras提供了ConvLSTM2D类，期望有着如下的输入：`[samples,timesteps,rows,cols,channels]`，而each time step of data可以被定义成一个`(rows*columns)`的数据点集合。

需要注意，我们所面对的数据是一个一维的总能耗数据，我们可以解释为一行与14列的数据，如果我们使用两周的数据作为输入的话。

对于ConvLSTM，其一次读取可以如下：每次读取一个14天的时间步长，并在这些时间步长上进行卷积。但这并不理想，与之相对，我们将14天切分成两个子序列，ConvLSTM可以再这两个time step上读取特征，使用CNN来处理每一个的七天数据。

输入格式为`[n,2,1,7,1]`。

* Samples:n，代表训练集合的样本数
* Time：2，代表我们切分一个14天的数据为两个子序列（两周）
* Rows:1，每一个子序列是一维的
* Columns：7，每个子序列有7天
* Channels:1，每个输入只有一个特征