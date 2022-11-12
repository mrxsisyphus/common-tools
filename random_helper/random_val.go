package random_helper

import (
	"fmt"
	"github.com/mrxtryagin/common-tools/collection_helper"
	"github.com/mrxtryagin/common-tools/error_helper"
	"github.com/mrxtryagin/common-tools/search_helper"
	"golang.org/x/exp/maps"
	"math"
	"math/rand"
	"time"
)

/**
关于随机数:
https://blog.csdn.net/aslackers/article/details/78548738
https://developer.aliyun.com/article/720190
https://wangbjun.site/2020/coding/golang/random.html
使用seed只需要在全局调用一次即可，如果多次调用则有可能取到相同随机数。
高并发下 不推荐使用 每次都seed的方式(考虑到高并发下会有相同的,所以可以自由的使用) crypto/rand真随机
*/

// InitSeed 可以用来全局注册seed的值
func InitSeed() {
	rand.Seed(time.Now().UnixNano())
}

// Shuffle 打乱数组,不改变原来的值,返回新的值
func Shuffle[T any](input *[]T) *[]T {
	raw := *input
	if len(raw) <= 0 {
		return input
	}
	// 数组拷贝
	res := make([]T, len(raw), len(raw))
	copy(res, raw)
	//rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(res), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})
	return &res
}

// Shuffled 打乱数组,改变原来的值
func Shuffled[T any](input *[]T) {
	raw := *input
	if len(raw) <= 0 {
		return
	}
	//rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(raw), func(i, j int) {
		raw[i], raw[j] = raw[j], raw[i]
	})
}

// Choice 从一个元素数组中随机选择一个元素
func Choice[T any](seq *[]T) T {
	raw := *seq
	if len(raw) <= 0 {
		panic("seq is empty")
	}
	//rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(raw))
	return raw[randomIndex]
}

// Sample 样本 从一个元素数组中随机选择元素num个,每次选择后不放回,不会有重复的
// refer: https://www.cnblogs.com/xybaby/p/8280936.html
// 参考python的实现
func Sample[T any](seq *[]T, num int) *[]T {
	raw := *seq
	length := len(raw)
	if length <= 0 {
		panic("seq is empty")
	}
	if length < num {
		panic("seq_length < sample num")
	}

	//思路:
	// 1. 随机抽取 不放回,也就是每次抽取的都是没被抽过的,需要原序列删除.
	// 2, 随机抽取 放回,外部维护是否抽过,如果抽过,重新抽取,不影响原序列 o(NlogN)
	// 3. shuffle, 然后取前K个 o(N) 需要遍历数组个数
	// 4, 部分shuffle 得到k个元素就返回 o(k) 遍历输入的num次数就行
	//　因此，算法的实现主要考虑的是额外使用的内存如果list拷贝原序列内存占用少，那么用部分shuffle；如果set占用内存少，那么使用记录已选项的办法。

	//todo: 没有摸索时候go语言的特别数字,所以用比率来做计算, 目前 >= 0.6 就认为 抽出不重复的概率比较小 需要benckmark
	radio := float64(num) / float64(length)

	//rand.Seed(time.Now().UnixNano())
	if radio >= 0.6 {
		return SampleWithPartlyShuffle[T](seq, num)
	} else {
		// 当 抽选数量占总数比例较小的时候,其实这里选到重复的概率更高,而不需要拷贝数组
		// 参考:https://www.jianshu.com/p/32902050f7dd 只需要不重复,不需要随机性,可以用补集的思路
		return SampleWithFullSeqRandom[T](seq, num)
	}
}

// SampleWithReturn 从一个元素数组中随机选择元素num个,每次选择后放回,所以会有重复的[类似于选num次]
func SampleWithReturn[T any](seq *[]T, num int) *[]T {
	raw := *seq
	length := len(raw)
	if length <= 0 {
		panic("seq is empty")
	}
	if length < num {
		panic("seq_length < choice num")
	}
	res := make([]T, num, num)
	//rand.Seed(time.Now().UnixNano())
	for i := 0; i < num; i++ {
		// 随机坐标
		randomIndex := rand.Intn(length)
		res[i] = raw[randomIndex]
	}
	return &res
}

// SampleWithPartlyShuffle 部分洗牌 来获取不重复样本
// 参考python random.sample
func SampleWithPartlyShuffle[T any](input *[]T, num int) *[]T {
	raw := *input
	length := len(raw)
	if length <= 0 {
		return input
	}
	if length < num {
		panic("seq_length < sample num")
	}
	rawCopy := make([]T, length, length)
	// copy 一份,这样就不会修改原序列了
	copy(rawCopy, raw)
	res := make([]T, num, num)
	for i := 0; i < num; i++ {
		randomIndex := rand.Intn(length - i)       // 抽取范围会越来越小
		res[i] = rawCopy[randomIndex]              // 赋值
		rawCopy[randomIndex] = rawCopy[length-i-1] // 将下一轮不会选到的最后一个,复制到前面
	}
	return &res
}

// SampleWithFullSeqRandom 使用全seq长度的随机(拿出后放回)+ map_helper(去重) 来获取不重复样本
func SampleWithFullSeqRandom[T any](input *[]T, num int) *[]T {
	raw := *input
	length := len(raw)
	if length <= 0 {
		return input
	}
	if length < num {
		panic("seq_length < sample num")
	}
	res := make([]T, num, num)
	choiced := make(map[int]struct{})
	for i := 0; i < num; i++ {
		// 每次都是从全量中抽
		randomIndex := rand.Intn(length)
		_, exist := choiced[randomIndex]
		// 判断存不存在, 如果存在 知道不存在位置
		for exist {
			randomIndex = rand.Intn(length)
			_, exist = choiced[randomIndex]
		}
		choiced[randomIndex] = struct{}{}
		res[i] = raw[randomIndex]
	}
	return &res
}

// Choices 权重选择
// 思路的本质 是 某个范围内的随机数 通过各个范围 ,来判断目前这个随机数落于哪一块,就返回这一块
// 比如 A: 10 B:20 C:70
// 将这些权重求和: 10 + 20 + 70 = 100 然后 在[0,100)的范围内取随机数 看这个随机数落在哪个位置 类似于轮盘
// 结论:
// 1. 为什么要用累加权重? 可以理解成坐标的线段 ,可以吧 A + B + C 看成一个整体,ABC分别是上面的坐标点
// 2. 可以通过 二分查找算法 来快速找到上述所谓的位置.
/*
 概率分布转化为坐标线段表示如下图:

   0  10    30     100
    A     B      C
在转化的过程中实际上就相当于是对它做了一次累加权重.
看上去,求 100的概率的过程中,概率值必然会落到某个坐标点上,那么这个概率点在哪个位置,就说明命中了哪一个
所以:
A B C 累加权重是 10 30 100,
那怎么判断 概率值落在哪个区域?
使用二分查找:
将 A B C 的累加值 作为数组的 三个索引值 [10,30,100](其实2个就够了,因为概率是到100为止,所以不可能在100以后的)
所以其实 就2个桩  [10,30] 那么 你当前的累计值在哪个区间 就可以用 二分插入查找算法来处理.
比如 值是20,那么 插入20后 为了使它有序 他必须放在 10 和 30 中间,也就说明这个位置落在哪一个索引了.还要考虑到直接打到庄上取哪的问题.
*/
// refer: https://www.cnblogs.com/Mishell/p/14009383.html 随机权重抽奖
// refer: https://www.liujiangblog.com/course/python/57 二分查找和插入算法
// refer: https://segmentfault.com/q/1010000021810370 补足
// 参考python
// -seq 数据序列
// - weights 相对权重列表,数量要与序列对应
// - cumWeights 累加权重列表,数量要与序列对应
// - weights
// - num 需要选择的个数
func Choices[T any](seq *[]T, weights []float64, cumWeights []float64, num int) *[]T {
	raw := *seq
	rawLength := len(raw)
	if rawLength <= 0 {
		panic("seq is empty!")
	}
	weightsLength := len(weights)
	cumWeightsLength := len(cumWeights)
	// 不加num 给1
	if num == 0 {
		num = 1
	}

	// 结果数组
	res := make([]T, num, num)
	// 如果cumWeights 为 0
	if cumWeightsLength == 0 {
		if weightsLength == 0 {
			// 如果weights也为0,随机返回
			for i := 0; i < num; i++ {
				// 随机浮点数(数组下标) 类似于 SampleWithReturn
				randomIndex := int(math.Floor(Float64n(float64(rawLength))))
				res[i] = raw[randomIndex]
			}
			return &res
		} else {
			// 如果weights 不为0,则转化为cumWeights,后面统一用cumWeights 来算
			cumWeightsP := collection_helper.AccumulateWithSameType[float64](&weights, func(index int, x, y float64) float64 {
				if index == 0 {
					return y
				}
				return x + y
			})
			cumWeights = *cumWeightsP
			cumWeightsLength = len(cumWeights) // 更新 length
		}
	} else if weightsLength != 0 {
		// weightsLength 与 cumWeightsLength 二者 择其一
		if weightsLength != 0 && cumWeightsLength != 0 {
			panic("Cannot specify both weights and cumulative weights")
		}
	}
	if cumWeightsLength != rawLength {
		panic("The number of weights does not match the seq")
	}
	fmt.Println(cumWeights)
	// 取累计总和,也就是最后一个
	cumTotal := cumWeights[cumWeightsLength-1]
	end := rawLength - 1 // 排除掉最后累积的100,因为不需要100的这个桩
	for i := 0; i < num; i++ {
		// 0-cumTotal 范围内的随机值(转盘指针)
		randomVal := Float64n(cumTotal)
		// 观察这个值在 如果要插入到这个序列的话,在累积权重这个序列(被认为是有序的)插入点的位置(省去最后一项),相同返回右侧
		//为什么要取后面? 取前后的差别主要是打到同一个桩上的时候,假设这个是有序的,那么打到桩上,桩后的概率是要大于装前的所以取桩前(其实这个就是临界点时的取值问题),而且100也取不到
		insertIndex := search_helper.BisectRight(&cumWeights, randomVal, 0, end)
		res[i] = raw[insertIndex]
	}
	return &res
}

// ChoicesByFunc 权重参数来源于对象本身,不传递额外的参数,本质利用Choices
//   - seq *[]T 数据序列
//   - getWeightsFunc func(T) float64 从每个对象中获得相对权重
//     -handleWeightsHook func(*[]float64) *[]float64 在真正choice之前处理这些权重,可以不填,可以用于做一些概率反转等
//
// - num 需要选择的个数
func ChoicesByFunc[T any](seq *[]T, getWeightsFunc func(T) float64, handleWeightsHook func(*[]float64) *[]float64, num int) *[]T {
	raw := *seq
	rawLength := len(raw)
	if rawLength <= 0 {
		panic("seq is empty!")
	}
	if getWeightsFunc == nil {
		panic("you need getWeightsFunc to get weight from seq data")
	}

	//从数组中剥离出weights
	weights := collection_helper.Map[T, float64](seq, func(index int, input T) float64 {
		return getWeightsFunc(input)
	})
	//需不需要在进行处理
	if handleWeightsHook != nil {
		weights = handleWeightsHook(weights)
	}
	// 进行choices
	return Choices[T](seq, *weights, nil, num)
}

// RandInt 随机生成一个start,stop之间的整数 因为默认的 rand.Intn方法只能生成0-n
func RandInt(start, stop int) int {
	if start > stop {
		error_helper.PanicWithStr("RandInt: start[%d] > stop[%d]", start, stop)
	}
	//rand.Seed(time.Now().UnixNano())
	return rand.Intn(stop-start) + start
}

// RandFloat 随机生成一个start,stop之间的浮点数
func RandFloat(start, stop float64) float64 {
	if start > stop {
		error_helper.PanicWithStr("RandFloat: start[%d] > stop[%d]", start, stop)
	}
	//rand.Seed(time.Now().UnixNano())
	return Float64n(stop-start) + start
}

// Float64n float64的浮点数范围 [0.0 ~ 1.0 * n)
func Float64n(n float64) float64 {
	return rand.Float64() * n
}

// RandIntWithSameSeed 随机生成一个start,stop之间的整数 因为默认的 rand.Intn方法只能生成0-n,种子数保持不变
func RandIntWithSameSeed(start, stop int) int {
	if start > stop {
		error_helper.PanicWithStr("RandInt: start[%d] > stop[%d]", start, stop)
	}
	return rand.Intn(stop-start) + start
}

// GenerateRandomNumberWithoutRepeat 从 start-> end 中 生成count个不重复的序列
func GenerateRandomNumberWithoutRepeat(start, end, count int) *[]int {
	if start < 0 || end < 0 || count < 0 || end < start || (end-start) < count {
		return nil
	}
	nums := make(map[int]struct{}, count)
	//rand.Seed(time.Now().UnixNano())
	for len(nums) < count {
		num := RandInt(start, end)

		if _, exist := nums[num]; !exist {
			nums[num] = struct{}{}
		}
	}
	res := maps.Keys(nums)
	return &res
}
