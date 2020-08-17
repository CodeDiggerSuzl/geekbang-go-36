# Golang 中需要注意的地方

## map 注意点

<details> <summary> doc </summary>

map 在大多数语言中都是基于 hash table 来实现的。

hash 表通常会把键转换成 hash 值（一个 unsigned int 值），哈希表会持有一定数量的桶（bucket 哈希桶），会均匀的储存哈希表收纳的键 - 值对儿。hash table 会先用键的 hash 值的低几位去定位到一个 bucket，然后去这个 bucket 中查找这个键。然后在用这个键去找到对应的值。

### map 的 key 不能为哪些类型？函数、字典和切片
在 go 中字典的键类型不能是那些类型：** 函数、字典类型和切片类型 **。这是因为键的类型必须可以是可以施加 `=` 或 `!=` 符号类型的。上面三种类型不支持判等操作。

另外如过键的类型是接口类型的，那么键值的类型实际类型也不能是上面 3 中类型，** 否则会在程序运行时候引发 panic** `panic: runtime error: hash of unhashable type []int`
```go
func Test_InterfaceKey(t *testing.T) {
	var mapT = map[interface{}]int{
		[]int{4}: 4,   // panic: runtime error: hash of unhashable type []int
		3:        3,
	}
}
```
另外，如果键的类型为：** 数组类型 **，那么应该确定该类型的元素类型不是 ** 函数类型、字典或切片类型 **。
### 为什么键必须支持判等操作呢？
每个 bucket 会把自己包含的所有键的 hash 值存起来，在 bucket 中定位的时候，会使用被查找的 hash 值和和 bucket 中 hash 值逐个进行比较。只有键的哈希值和键值都相等，才能说明查找 到了匹配的键 - 元素对。

### 优先选用哪些类型作为字典的键类型 ？
求哈希和判等操作的速度越快，对应的类型就越适合作为键类型。

对于基本类型来说：
- 所有的基本类型，指针类型，以及数组类型、结构体类型和接口类型，go 中都有一套 hash 算法来算对应的 hash 值，、
- 宽度越小的类型速度越快，`bool`、`int`、`float`、`complex` 和 `pointer` 类型。对于字符串来说，需要看他的具体长度，长度越短求 hash 越快。
- 类型的宽度是指它的单个值需要占用的字节数。比如，bool、int8 和 uint8 类型的一个值 需要占用的字节数都是 1，因此这些类型的宽度就都是 1。

对于高级类型来说：
- 对结构体类型的值求哈希实际上就是对它的所有字段值求哈希并进行合并，所以 关键在于它的各个字段的类型以及字段的数量。
- 对于接口类型，具体的哈希算法，则由值的实际类型决定。
- 把接口类型作为字典的键类型最危险。

** 优先选用数值类型和指针类型，通常情况下类型的宽度越小越好。**

### 在值为 `nil` 的字典上执行读和写操作会成功吗？
除了 ** 添加 ** 键 - 元素对，我们在一个值为 nil 的 map 上做任何操作都不会引 起错误。当我们试图在一个值为 nil 的字典中添加键 - 元素对的时候，Go 语言的运行时系统就会立即抛出一个 panic。

```go
func Test_Map(t *testing.T) {
	var m map[string]int
	key := "second"
	// get from a nil map. OK
	elem, ok := m["second"]
	// del from a nil map. OK
	delete(m, key)

	elem = 2
	// set in a nil map. NOT OK: panic: assignment to entry in nil map [recovered]
	m["second"] = elem
}
```

</details>


## channel
> Rob Pike: Don't communicate by sharing memory,share memory by communicating.
<details><summary>doc</summary>

通道类型本身的值就是并发安全的，是 Go 语言**自带的唯一一个可以满足并发安全性** 的类型。

在声明chan的时候，当容量为0时，我们可以称通道为非缓冲通道，也就是不带缓冲的通道。而当容量大于0 时，我们可以称为缓冲通道，也就是带有缓冲的通道。

**一个通道相当于一个先进先出(FIFO)的队列。**

### 对通道的发送和接收操作的基本特性?
1. 对于同一个通道，发送操作之间是互斥的，接收操作之间也是互斥的。
2. 发送操作和接收操作中对元素值的处理都是不可分割的。
3. 发送操作在完全完成之前会被阻塞。接收操作也是如此。

要注意的一个细节是，元素值从外界进入通道时会被复制。更具体地说，进入通道的并 不是在接收操作符右边的那个元素值，而是它的副本。

### 发送操作和接收操作在什么时候可能被长时间的阻塞?
- 缓冲通道: 如果通道已满，那么对它的所有发送操作都会被阻塞，直到通道 中有元素值被接收走。
- 非缓冲通道: 无论是发送操作还是接收操作，一开始执行就会被阻塞，直到配对的操作也开始执行，才会继续传递。(同步传递消息)

### 发送操作和接收操作在什么时候会引发 panic?
- 通道一旦关闭的情况下,在对其进行操作会发生 panic。
- 如果我们试图关闭一个已经关闭了的通道，也会引发 panic。

### 通道底层存储数据的是链表还是数组? 环形链表

### 单项通道
- `chan<-` send
- `<-chan` revc

> 感觉chanel有点像socket的同步阻塞模式，只不过channel的发送端和接收端共享一个缓 冲，套接字则是发送这边有发送缓冲，接收这边有接收缓冲，而且socket接收端如果先 close的话，发送端再发送数据的也会引发panic(linux上会触发SIG_PIPE信号，不处理程序就崩溃了)。
> Go语言里没有深层复制。数组是值类型，所以会被完全复制
> 通道必须要手动关闭吗?go会自动清理吗?   需要手动关闭，这是个很好的习惯，而且也可以利用关的动作来给接收方传递一个信号。Go的GC只会清理被分配到堆上的、不再有任何引用的对象。
> 不要从接受端关闭channel算是基本原则了，另外如果有多个并发发送者，1个或多个接收 者，有什么普适选择可以分享吗? 可以用另外的标志位做，比如context。
> 浅拷贝只是拷贝值以及值中直接包含的东西，深拷贝就是把所有深层次的结构一并拷贝.
</details>







## func
在 Go 语言中，函数可是一等的(first-class)公民，函数类型也是一等的数据类型。

简单来说，这意味着函数不但可以用于封装代码、分割功能、解耦逻辑，还可以化身为普通 的值，在其他函数间传递、赋予变量、做类型判断和转换等等，就像切片和字典的值那样。

而更深层次的含义就是:函数值可以由此成为能够被随意传播的独立逻辑组件(或者说功能 模块)。
<details><summary>doc</summary>

对于函数类型来说，它是一种对一组输入、输出进行模板化的重要工具，它比接口类型更加 轻巧、灵活，它的值也借此变成了可被热替换的逻辑组件。

```go
import "fmt"

type Printer func(content string) (n int, err error)

func printToStd(s string) (n int, err error) {
	return fmt.Println(s)
}

func main() {
	var p Printer
	p = printToStd
	p("hello")
}

```
### 高阶函数
什么是高阶函数:
1. 接受其他的函数作为参数传入;
2. 把其他的函数作为结果返回。

```go
type operate func(x, y int) int

func calc(x, y int, op operate) (int, error) {
	if op == nil {
		return 0, errors.New("invalid operation")
	}
	return op(x, y), nil
}

func main() {
	x, y := 12, 23
	op := func(x, y int) int {
		return x * y
	}
	result, err := calc(x, y, op)
	fmt.Printf("The result: %d (error: %v)\n", result, err)
}
```

### 闭包


</details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>