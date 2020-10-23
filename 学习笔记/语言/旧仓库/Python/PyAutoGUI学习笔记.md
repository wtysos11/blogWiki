# PyAutoGUI学习笔记

标签：python 第三方库 学习笔记 鼠标控制 键盘控制

## 简介

PyAutoGUI是一个Python的库，可以跨平台支持鼠标和键盘的控制（支持Windows,Mac和Linux），包括热键和各种组合键。它在Windows上不依赖于任何的库，可以直接使用pip进行安装，非常方便。

所有的函数都是基于导入的pyautogui实体。

## Official Cheat Sheat

[官方总结](https://pyautogui.readthedocs.io/en/latest/cheatsheet.html)

## 自己的总结

### 通用函数

* `position()`，返回二元组(x,y)，表示当前鼠标的坐标。
* `size()`，返回二元组（width,height），表示当前显示器宽度。
* `onScreen(x,y)`，接受两个参数，返回布尔值，表示坐标是否在显示器上

值得一提的是如果同时存在多个显示器（windows，显示器扩展）的话，它也只会显示当前分辨率大小，同时副显示器的鼠标坐标会变成负数，而且不是从0开始。（右边为-388与左边为-1920，原因未知。分辨率两边都是能够显示1920\*1080而且系统分辨率也是如此）

同样，onScreen也无法在有扩展显示器的情况下正常进行判定（辅显示器上的坐标全部为False）

### 鼠标控制函数

这个库使用的坐标系统与通常的计算机软件一样，以屏幕的左上角作为定点，横向向右为X轴正方向，纵向向下为y轴正方向。

* `moveTo(x, y, duration=num_seconds)`，鼠标移动函数，移动鼠标向坐标x,y，其中duration表示持续时间。moveTo的第四个参数可以接受鼠标移动速度控制参数，包括`pyautogui.easeInQuad`等，详见Doc。
* `moveRel(xOffset,yOffset,duration=num_seconds)`，作用和move()相同，鼠标移动函数，移动鼠标向当前坐标的偏移值，duration表示持续时间。
* dragTo和drag函数，与moveTo和move相似，表示拖动。
* `click(x=moveToX, y=moveToY, clicks=num_of_clicks, interval=secs_between_clicks, button='left')`，此外还有简化的doubleClick(),tripleClick(),rightClick()等。
* `scroll(amount_to_scroll, x=moveToX, y=moveToY)`，滚动滚轮的模拟。
* `mouseDown(x=moveToX, y=moveToY, button='left'),mouseUp(x=moveToX, y=moveToY, button='left')`，鼠标按下和鼠标抬起函数，可以分开。

duration不能小于最小延迟`pyautogui.MINIMUM_DURATION`，也就是0.1，不然不会有效果。

### 键盘控制

* `typewrite(str,interval=duration)`函数，直接输出字符串到文本框内，因为是模拟键盘上的动作，键盘上没有的将不会输出。（不是操控剪贴板之类的）
* press、keyDown和keyUp，press是由keyDown和keyUp组合起来的。此外还接受类似于`'enter' 'f1' 'left'(left arrow)`等按键，详见[按键查询](https://pyautogui.readthedocs.io/en/latest/keyboard.html#keyboard-keys)
* `hotkey`函数，等价于同时keyDown若干个键后抬起。

### Message Box

* `alert(text='', title='', button='OK')`，弹出警示框
* `confirm(text='', title='', buttons=['OK', 'Cancel'])`，弹出确认框。
* `prompt(text='', title='' , default='')`，提示框，可以输入。
* `password(text='', title='', default='', mask='*')`，密码框，同样可以输入，会进行遮掩。

### 截图函数

需要使用pillow模块，OSX下使用screencapture命令，linux下使用scrot命令

* screenshot()函数，获得一个PIL.Image.Image对象，传递一个string参数作为filename可以保存。对于1920\*1080的分辨率，全屏截图花费大约100ms，速度还可以。同样接受一个关键字region，是一个四元组，表示截图的区域。

### 定位函数

* `locateOnScreen(filename)`，接受一个图像的文件名，返回这个文件在图像中的位置。如果不能找到，会抛出`ImageNotFoundException`。不然，会返回一个四元组`Box(left,top,width,height)`，可以通过属性或者下标访问。
* 上述函数返回的四元组可以使用`pyautogui.center(Box)`返回一个二元组`Point(x,y)`，表示中心点的位置。
* `locateCenterOnScreen`，此函数为`locateOnScreen`和`center`的结合。

locate系列函数

* `locateOnScreen(image, grayscale=False)`，对于找到的第一个符合的对象，返回四元组。不然抛出`ImageNotFoundException`异常。
* `locateCenterOnScreen(image, grayscale=False)`，返回找到的第一个符合的对象的中心点，返回二元组。不然抛出`ImageNotFoundException`异常。
* `locateAllOnScreen(image, grayscale=False)`，返回一个迭代器，产生四元组。
* `locate(needleImage, haystackImage, grayscale=False)`，返回在haystackImage中找到的第一个needleImage，一个四元组。不然抛出`ImageNotFoundException`异常
* `locateAll(needleImage, haystackImage, grayscale=False)`，找到指定图像haystackImage中所有的needleImage，返回一个产生四元组的迭代器。

迭代器可以使用list()转换后使用。

一些参数

* 值得一提的是confidence参数，置信度参数，是一个0~1的值，需要安装opencv-python才能够正常运行。比较神奇的是本身在函数定义中并没有出现这个参数，不知道是如何实现的。
* grayscale=True，locate系列函数在灰度图情况进行匹配，会得到30%左右的加速。
* pyautogui.pixel(x,y)，拿去指定位置的rgb值，返回一个RGB对象，可以使用属性或者下标访问数据。
* pixelMatchesColor(x,y,(r,g,b))，返回一个布尔值，表示指定位置的RGB值是否与指定的RGB值相匹配。可以接受关键字tolerance，表示每个属性值允许的差值。