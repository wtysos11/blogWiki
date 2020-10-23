# gRPC学习笔记

标签：gRPC golang 微服务 学习笔记

## 来源

官方教程：

* [gRPC Basics-GO](https://grpc.io/docs/tutorials/basic/go.html)
* [官方quick-start](https://grpc.io/docs/quickstart/go.html#update-the-client)

中文例子：

* [Go-gRPC实践指南](https://www.golang123.com/book/36?chapterID=854)
* [Golang中的微服务](https://studygolang.com/articles/12060)

下面例子来自gRPC Basics-GO以及其[中文翻译](http://doc.oschina.net/grpc?t=60133)

下一步打算研究这本[Go-gRPC实践指南](https://www.bookstack.cn/read/go-grpc/readme.md)，里面有很多实践。

## 服务端

第一步是用protocol buffers定义request和response方法。

首先在.proto文件中定义一个service，然后在内部定义rpc方法，指定它们的request和response类型。例如

```proto
service RouteGuide{
    rpc GetFeature(Point) returns (Feature) {}
    rpc ListFeatures(Rectangle) returns (stream Feature) {}
    rpc RecordRoute(stream Point) returns (RouteSummary) {}
    rpc RouteChat(stream RouteNote) returns (stream RouteNote){}
}
```

使用message来定义protocol buffer message类型定义。例如下面的Point

```proto
message Point{
    int32 latitude = 1;
    int32 longitude = 2;
}
```

### 基本的RPC

`GetFeature(Point)`就是一个基本的RPC，它的go实现如下：

```go
func (s *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	// No feature was found, return an unnamed feature
	return &pb.Feature{"", point}, nil
}
```

### 服务端streaming RPC

如`ListFeatures`，接受一个参数，需要向客户端返回多个Feature。

```go
func (s *routeGuideServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	for _, feature := range s.savedFeatures {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}
```

### 客户端streaming RPC

如`RecordRoute`，函数接受一个流类型，返回错误类型。参数中的流即可读，也可写。

```go
func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	var pointCount, featureCount, distance int32
	var lastPoint *pb.Point
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointCount++
		for _, feature := range s.savedFeatures {
			if proto.Equal(feature.Location, point) {
				featureCount++
			}
		}
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}
		lastPoint = point
	}
}
```

### 双向流RPC

```go
func (s *routeGuideServer) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := serialize(in.Location)
                ... // look for notes to be sent to client
		for _, note := range s.routeNotes[key] {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}
```

### 开始工作

服务器启动代码及工作流程：

1. 指定工作端口信息
2. 创建gRPC实例
3. 将服务实现在gRPC服务器上进行注册
4. 调用服务器的Serve方法进行阻塞式等待，直到进程结束或者Stop方法被调用。

```golang
flag.Parse()
lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
if err != nil {
        log.Fatalf("failed to listen: %v", err)
}
grpcServer := grpc.NewServer()
pb.RegisterRouteGuideServer(grpcServer, &routeGuideServer{})
... // determine whether to use TLS
grpcServer.Serve(lis)
```

## 产生代码

使用protoc命令来产生代码，例如`protoc -I routeguide/ routeguide/route_guide.proto --go_out=plugins=grpc:routeguide`

## 创建客户端

### Creating a stub

为了调用服务的方法，需要创建一个gRPC通道来和服务器进行通信。如下，通过传递服务器地址和端口到grpc.Dial()来实现。

```golang
conn, err := grpc.Dial(*serverAddr)
if err != nil {
    ...
}
defer conn.Close()
```

一旦gRPC通道建立，我们需要一个client stub来进行RPC。使用产生的pb包中提供的`NewRouteGuideClient`来完成这一点：`client := pb.NewRouteGuideClient(conn)`

### 调用服务方法

#### 简单的RPC

调用简单RPC的方法和调用本地方法几乎是一样的。

```golang
feature, err := client.GetFeature(context.Background(), &pb.Point{409146138, -746188906})
if err != nil {
        ...
}
```

#### 服务端的Streaming RPC

即服务端会返回一个流的RPC。这里ListFeatures会返回一个Feature的流。

```golang
rect := &pb.Rectangle{ ... }  // initialize a pb.Rectangle
stream, err := client.ListFeatures(context.Background(), rect)
if err != nil {
    ...
}
for {
    feature, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
    }
    log.Println(feature)
}
```

#### 客户端的streaming RPC

客户端的streaming RPC接受一个流作为参数，实现上与服务端的RPC类似。

```golang
// Create a random number of random points
r := rand.New(rand.NewSource(time.Now().UnixNano()))
pointCount := int(r.Int31n(100)) + 2 // Traverse at least two points
var points []*pb.Point
for i := 0; i < pointCount; i++ {
	points = append(points, randomPoint(r))
}
log.Printf("Traversing %d points.", len(points))
stream, err := client.RecordRoute(context.Background())
if err != nil {
	log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
}
for _, point := range points {
	if err := stream.Send(point); err != nil {
		log.Fatalf("%v.Send(%v) = %v", stream, point, err)
	}
}
reply, err := stream.CloseAndRecv()
if err != nil {
	log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
}
log.Printf("Route summary: %v", reply)
```

#### 双向streaming RPC

```golang
stream, err := client.RouteChat(context.Background())
waitc := make(chan struct{})
go func() {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			close(waitc)
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
		}
		log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
	}
}()
for _, note := range notes {
	if err := stream.Send(note); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}
}
stream.CloseSend()
<-waitc
```