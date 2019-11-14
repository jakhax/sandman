import unittest
from solution import *

results1 = two_oldest_ages([1,5,87,45,8,8])
results2 = two_oldest_ages([6,5,83,5,3,18])

class TestSolution(unittest.TestCase):
    
    def test_oldest(self):
        self.assertEquals(results1[0], 45)

    @timeout_decorator.timeout(3)
    def test_random_tests(self):
        import random
        for i in range(100**100):
            l=random.sample(range(1,10000),random.randint(2,100))
            l1=sorted(l)[-1]
            self.assertEquals(two_oldest_ages(l)[1],l1)

