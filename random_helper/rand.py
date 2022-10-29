import random
import time

# random.seed()
# choice samples
# random.sample()

'''
相对权重:

累积权重:

累积权重的作用:
1. 保持升序,一般情况下 累积权重被认为是升序的
Return the index where to insert item x in list a, assuming a is sorted.
python 有一个 bisect 模块，用于维护有序列表。
bisect 模块实现了一个算法用于插入元素到有序列表。在一些情况下，这比反复排序列表或构造一个大的列表再排序的效率更高。Bisect 是二分法的意思，这里使用二分法来排序，它会将一个元素插入到一个有序列表的合适位置，这使得不需要每次调用 sort 的方式维护有序列表。
'''
# random.choices

if __name__ == '__main__':
    start = time.time()
    a = [[1, 2, 3], [2, 3, 4], [3, 4, 5]]
    b = random.choices(a, [1.0, 2.0, 7.0], k=100000000)
    print(time.time() - start)
