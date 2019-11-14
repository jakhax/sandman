import unittest
import sys
sys.path.insert(0,"/usr/lib/python3.6/code_kombat_test_frameworks")
from code_kombat_test_frameworks import CodeKombatTestRunner
import timeout_decorator
def load_tests(loader, tests, pattern):
    return loader.discover(".")
GLOBAL_TIMEOUT=3
timeout_decorator.timeout(GLOBAL_TIMEOUT)(unittest.main)(testRunner=CodeKombatTestRunner())