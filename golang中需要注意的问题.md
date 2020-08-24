```
_________     ______
__  ____/________  /_____ ______________ _
_  / __ _  __ \_  /_  __ `/_  __ \_  __ `/
/ /_/ / / /_/ /  / / /_/ /_  / / /  /_/ /
\____/  \____//_/  \__,_/ /_/ /_/_\__, /
                                 /____/
```
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

通道类型本身的值就是并发安全的，是 Go 语言 ** 自带的唯一一个可以满足并发安全性 ** 的类型。

在声明 chan 的时候，当容量为 0 时，我们可以称通道为非缓冲通道，也就是不带缓冲的通道。而当容量大于 0 时，我们可以称为缓冲通道，也就是带有缓冲的通道。

** 一个通道相当于一个先进先出 (FIFO) 的队列。**

### 对通道的发送和接收操作的基本特性?
1. 对于同一个通道，发送操作之间是互斥的，接收操作之间也是互斥的。
2. 发送操作和接收操作中对元素值的处理都是不可分割的。
3. 发送操作在完全完成之前会被阻塞。接收操作也是如此。

要注意的一个细节是，元素值从外界进入通道时会被复制。更具体地说，进入通道的并 不是在接收操作符右边的那个元素值，而是它的副本。

### 发送操作和接收操作在什么时候可能被长时间的阻塞?
- 缓冲通道: 如果通道已满，那么对它的所有发送操作都会被阻塞，直到通道 中有元素值被接收走。
- 非缓冲通道: 无论是发送操作还是接收操作，一开始执行就会被阻塞，直到配对的操作也开始执行，才会继续传递。(同步传递消息)

### 发送操作和接收操作在什么时候会引发 panic?
- 通道一旦关闭的情况下, 在对其进行操作会发生 panic。
- 如果我们试图关闭一个已经关闭了的通道，也会引发 panic。

### 通道底层存储数据的是链表还是数组? 环形链表

### 单项通道
- `chan<-` send
- `<-chan` revc

> 感觉 chanel 有点像 socket 的同步阻塞模式，只不过 channel 的发送端和接收端共享一个缓 冲，套接字则是发送这边有发送缓冲，接收这边有接收缓冲，而且 socket 接收端如果先 close 的话，发送端再发送数据的也会引发 panic(linux 上会触发 SIG_PIPE 信号，不处理程序就崩溃了)。
> Go 语言里没有深层复制。数组是值类型，所以会被完全复制
> 通道必须要手动关闭吗? go 会自动清理吗?   需要手动关闭，这是个很好的习惯，而且也可以利用关的动作来给接收方传递一个信号。Go 的 GC 只会清理被分配到堆上的、不再有任何引用的对象。
> 不要从接受端关闭 channel 算是基本原则了，另外如果有多个并发发送者，1 个或多个接收 者，有什么普适选择可以分享吗? 可以用另外的标志位做，比如 context。
> 浅拷贝只是拷贝值以及值中直接包含的东西，深拷贝就是把所有深层次的结构一并拷贝.
</details>







## func
在 Go 语言中，函数可是一等的 (first-class) 公民，函数类型也是一等的数据类型。

简单来说，这意味着函数不但可以用于封装代码、分割功能、解耦逻辑，还可以化身为普通 的值，在其他函数间传递、赋予变量、做类型判断和转换等等，就像切片和字典的值那样。

而更深层次的含义就是: 函数值可以由此成为能够被随意传播的独立逻辑组件 (或者说功能 模块)。
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

### 闭包 TODO TODO

</details>

## struct

<details><summary>doc</summary>

> 函数是独立的程序实体。我们可以声明有名字的函数，也可以声明没名字的函数，还可以把 它们当做普通的值传来传去。我们能把具有相同签名的函数抽象成独立的函数类型，以作为 一组输入、输出 (或者说一类逻辑组件) 的代表。

> 方法却不同，它需要有名字，不能被当作值来看待，最重要的是，它必须隶属于某一个类 型。方法所属的类型会通过其声明中的接收者 (receiver) 声明体现出来。

> 接收者声明就是在关键字 func 和方法名称之间的圆括号包裹起来的内容，其中必须包含确 切的名称和类型字面量。

> 接收者的类型其实就是当前方法所属的类型，而接收者的名称，则用于在当前方法中引用它 所属的类型的当前值。

更宽泛地讲，如果结构体类型的某个字段声明中只有一个类型名，那么该字段代表了什么?

Go 语言规范规定，如果一个字段的声明中只有字段的类型名而没有字段的名称，那么
它就是一个嵌入字段，也可以被称为匿名字段。我们可以通过此类型变量的名称后 跟 “`.`”，再后跟嵌入字段类型的方式引用到该字段。也就是说，嵌入字段的类型既是类型也是名称。

问题 1: Go 语言是用嵌入字段实现了继承吗?

> 这里强调一下，Go 语言中根本没有继承的概念，它所做的是通过嵌入字段的方式实现了类 型之间的组合。这样做的具体原因和理念请见 Go 语言官网的 FAQ 中的 Why is there no type inheritance?。

问题 2: 值方法和指针方法都是什么意思，有什么区别?

值方法:
```go
func (cat Cat) String() string {
	return fmt.Sprintf("%s (category: %s, name: %q)",
	cat.scientificName, cat.Animal.AnimalCategory, cat.name)
}
```
指针方法:
```go
func (cat *Cat) SetName(name string) {
	cat.name = name
}
```


> 1. 值方法的接收者是该方法所属的那个类型值的一个副本。我们在该方法内对该副本的修 改一般都不会体现在原值上，除非这个类型本身是某个引用类型 (比如切片或字典) 的 别名类型。而指针方法的接收者，是该方法所属的那个基本类型值的指针值的一个副本。我们在这样的方法内对该副本指向的值进行修改，却一定会体现在原值上。

* 一个自定义数据类型的方法集合中仅会包含它的所有值方法，而该类型的指针类型的方法集合却囊括了前者的所有方法，包括所有值方法和所有指针方法。*

> 2. ** 一个自定义数据类型的方法集合中仅会包含它的所有值方法，而该类型的指针类型的方法集合却囊括了前者的所有方法，包括所有值方法和所有指针方法。** 严格来讲，我们在这样的基本类型的值上只能调用到它的值方法。但是，Go 语言会适时 地为我们进行自动地转译，使得我们在这样的值上也能调用到它的指针方法。比如，在 Cat 类型的变量 cat 之上，之所以我们可以通过 `cat.SetName("monster")` 修改猫的名字，是因为 Go 语言把它自动转译为了 (&cat).SetName("monster")， 即: 先取 cat 的指针值，然后在该指针值上调用 SetName 方法。
> 3. 在后边你会了解到，一个类型的方法集合中有哪些方法与它能实现哪些接口类型是息息 相关的。如果一个基本类型和它的指针类型的方法集合是不同的，那么它们具体实现的 接口类型的数量就也会有差异，除非这两个数量都是零。比如，一个指针类型实现了某某接口类型，但它的基本类型却不一定能够作为该接口的 实现类型。能够体现值方法和指针方法之间差异的小例子我放在 demo30.go 文件里了，你可以参照 一下。


** 最后，再次强调，嵌入字段是实现类型间组合的一种方式，这与继承没有半点儿关系。Go 语言虽然支持面向对象编程，但是根本就没有 “继承” 这个概念。**

</details>

## interface
<details><summary>doc</summary>

在 Go 语言的语境中，当我们在谈论 “接口” 的时候，一定指的是接口类型。因为接口类型与其他数据类型不同，它是没法被实例化的。

更具体地说，我们既不能通过调用 new 函数或 make 函数创建出一个接口类型的值，也无法用字面量来表示一个接口类型的值。

接口类型声明中的这些方法所代表的就是该接口的方法集合。一个接口的方法集合就是它的全部特征。

e.g. 声明的类型 Dog 附带了 3 个方法。其中有 2 个值方法，分别是 Name 和 Category，另外 还有一个指针方法 SetName。

接口:
```go
type Pet interface {
	SetName(name string)
	Name() string
	Category() string
}
```

类型:
```go
type Dog struct {
	name string // 名字。
}
func (dog *Dog) SetName(name string) {
	dog.name = name
}
func (dog Dog) Name() string {
	return dog.name
}
func (dog Dog) Category() string {
	return "dog"
}
```
这就意味着，Dog 类型本身的方法集合中只包含了 2 个方法，也就是所有的值方法。而它的 指针类型 `*Dog` 方法集合却包含了 3 个方法，

也就是说，它拥有 Dog 类型附带的所有值方法和指针方法。又由于这 3 个方法恰恰分别是 Pet 接口中某个方法的实现，所以 `*Dog` 类型就成为了 Pet 接口的实现类型。

正因为如此，我可以声明并初始化一个 Dog 类型的变量 dog，然后把它的指针值赋给类型为 Pet 的变量 pet。

```go
dog := Dog{"little pig"}
var pet Pet = &dog
```
这里有几个名词需要你先记住。对于一个接口类型的变量来说，例如上面的变量 pet，我们赋给它的值可以被叫做它的 ** 实际值 (也称动态值)**，而该值的类型可以被叫做这个变量的 ** 实际类型 (也称动态类型)**。

比如，我们把取址表达式 `&dog` 的结果值赋给了变量 `pet`，这时这个结果值就是变量 `pet` 的 动态值，而此结果值的类型 `*Dog` 就是该变量的动态类型。

动态类型这个叫法是相对于静态类型而言的。对于变量 pet 来讲，** 它的静态类型就是 Pet， 并且永远是 Pet，但是它的动态类型却会随着我们赋给它的动态值而变化。**

接口 类型本身是无法被值化的。在我们赋予它实际的值之前，它的值一定会是 nil，这也是它的 零值。

反过来讲，一旦它被赋予了某个实现类型的值，它的值就不再是 nil 了。不过要注意，即使我们像前面那样把 dog 的值赋给了 pet，pet 的值与 dog 的值也是不同的。这不仅仅是副本与原值的那种不同。

> 当我们给一个接口变量赋值的时候，该变量的动态类型会与它的动态值一起被存储在一个专 用的数据结构中。

> 严格来讲，这样一个变量的值其实是这个专用数据结构的一个实例，而不是我们赋给该变量 的那个实际的值。所以我才说，pet 的值与 dog 的值肯定是不同的，无论是从它们存储的内 容，还是存储的结构上来看都是如此。不过，我们可以认为，这时 pet 的值中包含了 dog 值 的副本。

> 我们就把这个专用的数据结构叫做 `iface` 吧，在 Go 语言的 runtime 包中它其实就叫这个 名字。

> `iface` 的实例会包含两个指针，一个是指向类型信息的指针，另一个是指向动态值的指针。 这里的类型信息是由另一个专用数据结构的实例承载的，其中包含了动态值的类型，以及使 它实现了接口的方法和调用它们的途径，等等。

只要我们把一个有类型的 nil 赋给接口变量，那么这个变量的值就一定不会是那个真正的 nil。因此，当我们使用判等符号 == 判断 pet 是否与字面量 nil 相等的时候，答案一定会是 false。

那么，怎样才能让一个接口变量的值真正为 nil 呢? 要么只声明它但不做初始化，要么直接 把字面量 nil 赋给它。

### 接口组合和 struct 组合类似
接口类型间的嵌入要更简单一些，因为它不会涉及方法间的 “屏蔽”。只要组合的接口之间 有同名的方法就会产生冲突，从而无法通过编译，即使同名方法的签名彼此不同也会是如 此。因此，接口的组合根本不可能导致 “屏蔽” 现象的出现。

Go 语言团队鼓励我们声明体量较小的接口，并建议我们通过这种接口间的组合来扩展程 序、增加程序的灵活性。

</details>

## pointer
<details><summary>doc</summary>
再来看 Go 语言标准库中的 unsafe 包。unsafe 包中有一个类型叫做 Pointer，也代表 了 “指针”。

unsafe.Pointer 可以表示任何指向可寻址的值的指针，同时它也是前面提到的指针值和 uintptr 值之间的桥梁。也就是说，通过它，我们可以在这两种值之上进行双向的转换。 这里有一个很关键的词——可寻址的 (addressable)。在我们继续说 unsafe.Pointer 之前，需要先要搞清楚这个词的确切含义。

### 你能列举出 Go 语言中的哪些值是不可寻址的吗?
1. 第一个关键词: 不可变的。由于 Go 语言中的字符串值也是不可变的，所以对于一个字符串 类型的变量来说，基于它的索引或切片的结果值也都是不可寻址的，因为即使拿到了这种值 的内存地址也改变不了什么。
2. 算术操作的结果值属于一种临时结果。在我们把这种结果值赋给任何变量或常量之前，即使 能拿到它的内存地址也是没有任何意义的。第二个关键词: 临时结果。这个关键词能被用来解释很多现象。我们可以把各种对值字面量 施加的表达式的求值结果都看做是临时结果。
3. 第三个关键词: 不安全的。“不安全的” 操作很可能会破坏程序的一致性，引发不可预知的 错误，从而严重影响程序的功能和稳定性。
### 怎样通过 unsafe.Pointer 操纵可寻址的值?

```go

dog := Dog{"little pig"}
dogP := &dog
dogPtr := uintptr(unsafe.Pointer(dogP))
namePtr := dogPtr + unsafe.Offsetof(dogP.name)
nameP := (*string)(unsafe.Pointer(namePtr))

```
这里需要与 unsafe.Offsetof 函数搭配使用才能看出端倪。unsafe.Offsetof 函数用 于获取两个值在内存中的起始存储地址之间的偏移量，以字节为单位。

这两个值一个是某个字段的值，另一个是该字段值所属的那个结构体值。我们在调用这个函 数的时候，需要把针对字段的选择表达式传给它，比如 dogP.name。
</details>

## goroutine

<details><summary>doc</summary>
Go 语言不但有着独特的并发编程模型，以及用户级线程 goroutine，还拥 有强大的用于调度 goroutine、对接系统级线程的调度器。这个调度器是 Go 语言运行时系统的重要组成部分，它主要负责统筹调配 Go 并发编程模型 中的三个主要元素，即: G(goroutine 的缩写)、P(processor 的缩写) 和 M(machine 的缩写)。
其中的 M 指代的就是系统级线程。而 P 指的是一种可以承载若干个 G，且能够使这些 G 适时地与 M 进行对接，并得到真正运行的中介。

从宏观上说，G 和 M 由于 P 的存在可以呈现出多对多的关系。当一个正在与某个 M 对接 并运行着的 G，需要因某个事件 (比如等待 I/O 或锁的解除) 而暂停运行的时候，调度器
总会及时地发现，并把这个 G 与那个 M 分离开，以释放计算资源供那些等待运行的 G 使用。

而当一个 G 需要恢复运行的时候，调度器又会尽快地为它寻找空闲的计算资源 (包括 M) 并安排运行。另外，当 M 不够用时，调度器会帮我们向操作系统申请新的系统级线程，而 当某个 M 已无用时，调度器又会负责把它及时地销毁掉。

正因为调度器帮助我们做了很多事，所以我们的 Go 程序才总是能高效地利用操作系统和计 算机资源。程序中的所有 goroutine 也都会被充分地调度，其中的代码也都会被并发地运 行，即使这样的 goroutine 有数以十万计，也仍然可以如此。
### 什么是主 goroutine，它与我们启用的其他 goroutine 有什么 不同?
与一个进程总会有一个主线程类似，每一个独立的 Go 程序在运行时也总会有一个主 goroutine。这个主 goroutine 会在 Go 程序的运行准备工作完成后被自动地启用，并不 需要我们做任何手动的操作。

**** imp

想必你已经知道，每条 go 语句一般都会携带一个函数调用，这个被调用的函数常常被称为 go 函数。而主 goroutine 的 go 函数就是那个作为程序入口的 main 函数。


### 怎样才能让主 goroutine 等待其他 goroutine?
方式 1 `time.Sleep()`

方式 2 使用 channel
```go
func main() {
	num := 10
	sign := make(chan struct{}, num)
	for i := 0; i < num; i++ {
		go func(i int) {
			fmt.Println(i)
			sign <- struct{}{}
		}(i)
	}
	for j := 0; j < num; j++ {
		<-sign
	}
}
```
- 我在声明通道 sign 的时候是以 chan struct{} 作为其类型 的。其中的类型字面量 struct{} 有些类似于空接口类型 interface{}，它代表了既不包 含任何字段也不拥有任何方法的空结构体类型。
- 注意，struct{} 类型值的表示法只有一个，即: struct{}{}。并且，它占用的内存空间 是 0 字节。确切地说，这个值在整个 Go 程序中永远都只会存在一份。虽然我们可以无数次 地使用这个值字面量，但是用到的却都是同一个值.

方式 3 使用 sync.WaitGroup

### 问题 2: 怎样让我们启用的多个 goroutine 按照既定的顺序运行?

```go
for i := 0; i < num; i++ {
	go func() {
		fmt.Println(i)
	}()
}
```

只有这样，Go 语言才能保证每个 goroutine 都可以拿到一个唯一的整数。其原因与 go 函 数的执行时机有关。

我在前面已经讲过了。在 go 语句被执行时，我们传给 go 函数的参数 i 会先被求值，如此就得 到了当次迭代的序号。之后，无论 go 函数会在什么时候执行，这个参数值都不会变。也就 是说，go 函数中调用的 fmt.Println 函数打印的一定会是那个当次迭代的序号。

```go
package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	var count uint32
	// trigger 函数会不断地获取一个名叫 count 的变量的值，并判断该值是否与参数 i 的值相同。
	// 如果相同，那么就立即调用 fn 代表的函数，然后把 count 变量的值加 1，最后显式地退出当前的循环。
	// 否则，我们就先让当前的 goroutine“睡眠” 一个纳秒再进入下一个迭代。
	trigger := func(i uint32, fn func()) {
		for {
			if n := atomic.LoadUint32(&count); n == i {
				fn()
				// 操作变量 count 的时候使用的都是原子操作。
				// 这是由于 trigger 函数会被多个 goroutine 并发地调用，所以它用到的非本地变量 count，就被多个用户级线程共用了。
				// 因此，对它的操作就产生了竞态条件 (race condition)，破坏了程序的并发安全性。
				atomic.AddUint32(&count, 1)
				break
			}
			time.Sleep(time.Nanosecond)
		}
	}

	// 在 go 函数中先声明了一个匿名的函数，并把它赋给了变量 fn。
	// 这个匿名函数做的事情很 简单，只是调用 fmt.Println 函数以打印 go 函数的参数 i 的值。
	for i := uint32(0); i < 10; i++ {
		go func(i uint32) {
			fn := func() {
				fmt.Println(i)
			}
			// 调用了一个名叫 trigger 的函数，并把 go 函数的参数 i 和刚刚声明的变量 fn 作为参数传给了它。
			// 注意，for 语句声明的局部变量 i 和 go 函数的参数 i 的类型都变了，都由 int 变为了 uint32
			trigger(i, fn)
		}(i)
	}
	// 该函数接受两个参数，一个是 uint32 类型的参数 i, 另一个是 func() 类型的参数 fn。
	// 你应该记得，func() 代表的是既无参数声明也无结果声明的函数类型。
	trigger(10, func() {})
}
```

主 goroutine 的运行若过早结束，那么我们的并发程序的功能就很可能无法全部完成。所 以我们往往需要通过一些手段去进行干涉，比如调用 time.Sleep 函数或者使用通道。我 们在后面的文章中还会讨论更高级的手段。

另外，go 函数的实际执行顺序往往与其所属的 go 语句的执行顺序 (或者说 goroutine 的启 用顺序) 不同，而且默认情况下的执行顺序是不可预知的。

### runtime 包中提供了哪些与模型三要素 G、P 和 M 相关的函数?(模型三要素内容在上 一篇)

</details>

## if for & switch TODO2
### for
1. range 表达式只会在 for 语句开始执行时被求值一次，无论后边会有多少次迭代;
2. range 表达式的求值结果会被复制，也就是说，被迭代的对象是 range 表达式结果值的
副本而不是原值。
```go
func Test_Range_Cpy(t *testing.T) {
	b := [3]int{1, 2, 3}
	for i, v := range b { //i,v 从 a 复制的对象里提取出
		if i == 0 {
			b[1], b[2] = 200, 300
			fmt.Println(b) // 输出 [1 200 300]
		}
		b[i] = v + 100 //v 是复制对象里的元素 [1, 2, 3]
	}
	fmt.Println(b) // 输出 [101, 102, 103]
}

func Test_Range_Cpy(t *testing.T) {
	b := []int{1, 2, 3}
	for i, v := range b { //i,v 从 a 复制的对象里提取出
		if i == 0 {
			b[1], b[2] = 200, 300
			fmt.Println(b) // 输出 [1 200 300]
		}
		b[i] = v + 100
	}
	fmt.Println(b) // [101 300 400]
}
```
## error
<details><summary>doc</summary>

</details>

# sync
<details><summary>doc</summary>

一旦数据被多个线程共享，那么就很可能会产生争用和冲突的情况。这种情况也被称为竞态 条件 (race condition)，这往往会破坏共享数据的一致性。

共享数据的一致性代表着某种约定，即: 多个线程对共享数据的操作总是可以达到它们各自预期的效果。

** 同步的用途有两个，一个是避免多个线程在同一时刻操作同一个数据块，另一个 是协调多个线程，以避免它们在同一时刻执行同一个代码块。**

一个线程在想要访问某一个共享资源的时候，需要先申请对该资源的访问权限，并且只有在 申请成功之后，访问才能真正开始。而当线程对共享资源的访问结束时，它还必须归还对该资源的访问权限，若要再次访问仍需 申请。

在 Go 语言中，可供我们选择的同步工具并不少。其中，最重要且最常用的同步工具当属 互斥量 (mutual exclusion，简称 mutex)。sync 包中的 Mutex 就是与其对应的类型， 该类型的值可以被称为互斥量或者互斥锁。

最后，需要特别注意的是，无论是互斥锁还是读写锁，我们都不要试图去解锁未锁定的锁， 因为这样会引发不可恢复的 panic。


### 竞态条件、临界区与同步工具

### 条件变量
条件变量提供的方法有三个: 等待通知 (wait)、单发通知 (signal) 和广播通知 (broadcast)。

条件变量是基于互斥锁的，它必须有互斥锁的支撑才能够起作 用。因此，这里的参数值是不可或缺的，它会参与到条件变量的方法实现当中。


```go
var mailbox uint8
var lock sync.RWMutex
sendCond := sync.NewCond(&lock)
recvCond := sync.NewCond(lock.RLocker())
```
条件变量的 Wait 方法主要做了四件事。

1. 把调用它的 goroutine(也就是当前的 goroutine) 加入到当前条件变量的通知队列中。
2. 解锁当前的条件变量基于的那个互斥锁。
3. 让当前的 goroutine 处于等待状态，等到通知到来时再决定是否唤醒它。此时，这个
goroutine 就会阻塞在调用这个 Wait 方法的那行代码上。
4. 如果通知到来并且决定唤醒这个 goroutine，那么就在唤醒它之后重新锁定当前条件变
量基于的互斥锁。自此之后，当前的 goroutine 就会继续执行后面的代码了。

因为条件变量的 Wait 方法在阻塞当前的 goroutine 之前，会解锁它基于的互斥锁，所以在 调用该 Wait 方法之前，我们必须先锁定那个互斥锁，否则在调用这个 Wait 方法时，就会引发一个不可恢复的 panic。

### 2. 为什么要用 for 语句来包裹调用其 Wait 方法的表达式，用 if 语句不行吗?
这主要是为了保险起见。如果一个 goroutine 因收到通知而被唤醒，但却发现共享资源的状态，依然不符合它的要求，那么就应该再次调用条件变量的 Wait 方法，并继续等待下次 通知的到来。

### 条件变量的 Signal 方法和 Broadcast 方法有哪些异同?
条件变量的 Signal 方法和 Broadcast 方法都是被用来发送通知的，不同的是，前者的通 知只会唤醒一个因此而等待的 goroutine，而后者的通知却会唤醒所有为此等待的 goroutine。

_条件变量的 Wait 方法总会把当前的 goroutine 添加到通知队列的队尾_，而它的 Signal 方法总会从通知队列的队首开始，查找可被唤醒的 goroutine。所以，因 Signal 方法的通 知，而被唤醒的 goroutine 一般都是最早等待的那一个。

这两个方法的行为决定了它们的适用场景。如果你确定只有一个 goroutine 在等待通知， 或者只需唤醒任意一个 goroutine 就可以满足要求，那么使用条件变量的 Signal 方法就好了。

此外，再次强调一下，与 Wait 方法不同，_条件变量的 Signal 方法和 Broadcast 方法并不需要在互斥锁的保护下执行_。恰恰相反，我们最好在解锁条件变量基于的那个互斥锁之后， 再去调用它的这两个方法。这更有利于程序的运行效率。

最后，请注意，条件变量的通知具有即时性。也就是说，如果发送通知的时候没有 goroutine 为此等待，那么该通知就会被直接丢弃。在这之后才开始等待的 goroutine 只 可能被后面的通知唤醒。

</details>

## 原子包

<details><summary>doc</summary>
```go
go func() {
	defer func() {
		sign <- struct{}{}
	}()
	// 使用 CAS 实现自旋锁 spinlock
	for {
		if atomic.CompareAndSwapInt32(&num, 10, 0) {
			fmt.Println("The number has go to 0")
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
}()
```

使用 CAS 实现自旋锁 spinlock

- 在 for 语句中的 CAS 操作可以不停地检查某个需要满足的条件，一旦条件满足就退出 for 循环。这就相当于，只要条件未被满足，当前的流程就会被一直 “阻塞” 在这里。

- 这在效果上与互斥锁有些类似。不过，它们的适用场景是不同的。_我们在使用互斥锁的时候，总是假设共享资源的状态会被其他的 goroutine 频繁地改变。_

- 而 for 语句加 CAS 操作的假设往往是: _共享资源状态的改变并不频繁，或者，它的状态总会变成期望的那样。_ 这是一种更加乐观，或者说更加宽松的做法。



</details>

## WaitGroup and Once

<details><summary>doc</summary>

```go
func Test_Once3(t *testing.T) {
	once := sync.Once{}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer func() {
			if p := recover(); p != nil {
				fmt.Printf("fatal error:%v\n", p)
			}
		}()
		once.Do(func() {
			fmt.Println("Do task. [4]")
			panic(errors.New("something wrong"))
		})
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 500)
		once.Do(func() {
			fmt.Println("Do task. [5]") //  will not run here cause the Once
		})
		fmt.Println("Done. [5]")
	}()
	wg.Wait()
}
```

</details>

# Context
<details><summary>doc</summary>

### 撤销信号是如何在上下文树中传播的?

1. context 包的 WithCancel 函数在被调用后会产生两个结果值。第 1 个结果值就是那个可 撤销的 Context 值，而第 2 个结果值则是用于触发撤销信号的函数。
2. 在撤销函数被调用之后，对应的 Context 值会先关闭它内部的接收通道，也就是它的 Done 方法会返回的那个通道。
3. 然后，它会向它的所有子值 (或者说子节点) 传达撤销信号。这些子值会如法炮制，把撤销 信号继续传播下去。最后，这个 Context 值会断开它与其父值之间的关联。
4. 我们通过调用 context 包的 WithDeadline 函数或者 WithTimeout 函数生成的 Context 值也是可撤销的。它们不但可以被手动撤销，还会依据在生成时被给定的过期时间，自动地 进行定时撤销。这里定时撤销的功能是借助它们内部的计时器来实现的。


### 怎样通过 Context 值携带数据? 怎样从中获取数据?
WithValue 函数在产生新的 Context 值 (以下简称含数据的 Context 值) 的时候需要三个 参数，即: 父值、键和值。与 “字典对于键的约束” 类似，这里键的类型必须是可判等的。

原因很简单，当我们从中获取数据的时候，它需要根据给定的键来查找对应的值。不过，这 种 Context 值并不是用字典来存储键和值的，后两者只是被简单地存储在前者的相应字段 中而已。

Context 类型的 Value 方法就是被用来获取数据的。在我们调用含数据的 Context 值的 Value 方法时，它会先判断给定的键，是否与当前值中存储的键相等，如果相等就把该值中存储的值直接返回，否则就到其父值中继续查找。

如果其父值中仍然未存储相等的键，那么该方法就会沿着上下文根节点的方向一路查找下去。

注意，除了含数据的Context值以外，其他几种Context值都是无法携带数据的。因此， Context值的Value方法在沿路查找的时候，会直接跨过那几种值。
</details>

## 临时对象池 sync.Pool TODO
sync.Pool类型可以被称为临时对象池，它的值可以被用来存储临时的对象。与 Go 语言 的很多同步工具一样，sync.Pool类型也属于结构体类型，它的值在被真正使用之后，就不应该再被复制了。

_你可能已经想到了，我们可以把临时对象池当作针对某种数据的缓存来用。实际上，在我看 来，临时对象池最主要的用途就在于此。_
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
<details><summary>doc</summary></details>
