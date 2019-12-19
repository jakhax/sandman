from __future__ import absolute_import
from django.test.runner import DiscoverRunner
from .test_runner import CodewarsTestRunner
# author: https://github.com/Codewars/python-unittest

class CodewarsDjangoRunner(DiscoverRunner):
    def run_suite(self, suite, **kwargs):
        return CodewarsTestRunner(group_by_module=True).run(suite)
