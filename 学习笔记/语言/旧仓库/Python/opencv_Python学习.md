# opencv-python学习

标签：opencv python

来源：[OpenCV-Python Tutorials](https://opencv-python-tutroals.readthedocs.io/en/latest/py_tutorials/py_tutorials.html)

## GUI Features in OpenCV

### 基本图像操作

#### 读图

读图像：`cv.imread()`，第二个参数可以指定读取图像的方式

* `cv2.IMREAD_COLOR`：读取彩色图像，透明度参数会被忽视，此为默认。
* `cv2.IMREAD_GRAYSCALE`：读取灰度图
* `cv2.IMREAD_UNCHANGED`：包括alpha通道

除了上述三个标签外，可以使用1、0、-1简单替代。

返回对象是numpy的二维数组，即使不存在图像也不会抛出错误，只会返回None。

#### 展示图像

展示图像：`cv2.imshow(imgName,imgObject)`，窗口大小与图像大小相同。后面必须加上`cv2.waitKey(0)`，否则要么一闪而过，要么卡死。

第一个参数是string，代表窗口名字。第二个是二维numpy数组，代表图像。窗口名字不能够相同。

`cv2.destroyAllWindow()`，关闭所有创建的窗口。如果想要关闭指定窗口，使用`cv2.destroyWindow()`加上想要关闭窗口的名字。

#### 写图像

`cv2.imwrite()`将图像从内存写入到磁盘中。第一个参数为文件名，第二个参数为图像变量，即二维numpy数组。

### 基本视频操作

#### 调用摄像头

下面代码打开摄像头，并逐帧返回实时灰度图像。按q可以关闭。

```python
import numpy as np
import cv2

cap = cv2.VideoCapture(0)

while(True):
    # Capture frame-by-frame
    ret, frame = cap.read()

    # Our operations on the frame come here
    gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)

    # Display the resulting frame
    cv2.imshow('frame',gray)
    if cv2.waitKey(1) & 0xFF == ord('q'):
        break

# When everything done, release the capture
cap.release()
cv2.destroyAllWindows()
```

#### 从文件打开视频

```python
import numpy as np
import cv2

cap = cv2.VideoCapture('vtest.avi')

while(cap.isOpened()):
    ret, frame = cap.read()

    gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)

    cv2.imshow('frame',gray)
    if cv2.waitKey(1) & 0xFF == ord('q'):
        break

cap.release()
cv2.destroyAllWindows()
```

#### 保存视频

使用`VideoWriter`对象，指定输出文件名，指定`FourCC code`，以及帧数fps和帧的大小。最后是`isColor`标签。其中`FourCC code`是一个4字节代码，用于指定视频编码方式。

下面代码从摄像头读取视频并保存

```python
import numpy as np
import cv2

cap = cv2.VideoCapture(0)

# Define the codec and create VideoWriter object
fourcc = cv2.VideoWriter_fourcc(*'XVID')
out = cv2.VideoWriter('output.avi',fourcc, 20.0, (640,480))

while(cap.isOpened()):
    ret, frame = cap.read()
    if ret==True:
        frame = cv2.flip(frame,0)

        # write the flipped frame
        out.write(frame)

        cv2.imshow('frame',frame)
        if cv2.waitKey(1) & 0xFF == ord('q'):
            break
    else:
        break

# Release everything if job is finished
cap.release()
out.release()
cv2.destroyAllWindows()
```

### 画图函数

通用参数：
 
* img：你想要画图的对象
* color：画上去的颜色，是一个三元组。如果是灰度图则是一个标量。
* thickness：厚度，默认为1.
* lineType：线型，通常为8连通。`cv2.LINE_AA`是一种对于曲线十分不错的线型。

#### 画直线

直线需要给出开头和结尾（二元组），以及颜色和线的粗细。

```python
import numpy as np
import cv2

# Create a black image
img = np.zeros((512,512,3), np.uint8)

# Draw a diagonal blue line with thickness of 5 px
img = cv2.line(img,(0,0),(511,511),(255,0,0),5)
```

#### 矩形

矩形需要给出左上角和右下角坐标：`img = cv2.rectangle(img,(384,0),(510,128),(0,255,0),3)`

#### 圆形

画圆需要给出圆心和半径：`img = cv2.circle(img,(447,63), 63, (0,0,255), -1)`

#### 多边形

```python
pts = np.array([[10,5],[20,30],[70,20],[50,10]], np.int32)
pts = pts.reshape((-1,1,2))
img = cv2.polylines(img,[pts],True,(0,255,255))
```

#### 添加文字

```python
font = cv2.FONT_HERSHEY_SIMPLEX
cv2.putText(img,'OpenCV',(10,500), font, 4,(255,255,255),2,cv2.LINE_AA)
```