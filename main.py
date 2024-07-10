import os
import decimal
from decimal import Decimal, getcontext

# 设置初始Decimal精度
initial_prec = 100
getcontext().prec = initial_prec

# 文件路径
FILENAME = 'pi_digits.txt'


# Chudnovsky算法计算圆周率
def chudnovsky_algorithm(prec):
    getcontext().prec = prec
    C = 426880 * Decimal(10005).sqrt()
    K = 6
    M = 1
    X = 1
    L = 13591409
    S = L

    for k in range(1, prec // 14 + 1):  # 迭代次数取决于精度
        M = (K ** 3 - 16 * K) * M // (k ** 3)
        L += 545140134
        X *= -262537412640768000
        S += Decimal(M * L) / X
        K += 12

    pi = C / S
    return pi


# 读取已有的圆周率位数
def read_existing_digits():
    if os.path.exists(FILENAME):
        with open(FILENAME, 'r') as f:
            return f.read().replace('\n', '')
    return ''


# 写入新的圆周率位数
def write_digits(digits):
    with open(FILENAME, 'a') as f:
        f.write(digits)


# 计算圆周率函数
def calculate_pi():
    existing_digits = read_existing_digits()
    current_position = len(existing_digits)

    current_prec = max(initial_prec, current_position + 100)  # 根据已有位数动态设置初始精度

    while True:
        current_prec += 50  # 增加精度
        pi = str(chudnovsky_algorithm(current_prec))[2:]  # 去除"3."

        while current_position < len(pi):
            new_digits = pi[current_position:current_position + 50]
            if not new_digits:
                break
            write_digits(new_digits)  # 写入新计算的位数并换行
            current_position += 50


if __name__ == "__main__":
    calculate_pi()
