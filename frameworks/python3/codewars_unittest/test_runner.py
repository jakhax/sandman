import sys
import unittest
import warnings
from itertools import groupby
# author: https://github.com/Codewars/python-unittest

# Use timeit.default_timer for Python 2 compatibility.
# default_timer is time.perf_counter on 3.3+
from timeit import default_timer as perf_counter

from .test_result import CodewarsTestResult


class CodewarsTestRunner(object):
    def __init__(self, stream=None, group_by_module=False, warnings=None):
        if stream is None:
            stream = sys.stdout
        self.stream = _WritelnDecorator(stream)
        self.group_by_module = group_by_module
        self.results = []
        if warnings is None and not sys.warnoptions:
            self.warnings = "default"
        else:
            # Set to given `warnings` or use `None` to respect values passed to `-W`
            self.warnings = warnings

    def run(self, test):
        with warnings.catch_warnings():
            if self.warnings:
                warnings.simplefilter(self.warnings)
                # Minimize noise from deprecated assertion methods
                if self.warnings in ["default", "always"]:
                    warnings.filterwarnings(
                        "module",
                        category=DeprecationWarning,
                        message=r"Please use assert\w+ instead.",
                    )
            if isinstance(test, unittest.TestSuite):
                self._run_modules(_to_tree(_flatten(test)))
            else:
                self._run_case(test)
        return self._make_result()

    def _make_result(self):
        accum = unittest.TestResult()
        for result in self.results:
            accum.failures.extend(result.failures)
            accum.errors.extend(result.errors)
            accum.testsRun += result.testsRun
        return accum

    def _run_modules(self, modules):
        for mod in modules:
            name = ""
            if self.group_by_module:
                name = mod.group_name
                # Don't group on ImportError
                if name == "unittest.loader":
                    name = ""
                if name:
                    self.stream.writeln(_group(name))

            startTime = perf_counter()
            for cases in mod:
                self._run_cases(cases)

            if name:
                self.stream.writeln(_completedin(startTime, perf_counter()))

    def _run_cases(self, test):
        name = test.group_name
        # Don't group when errored before running tests, e.g., ImportError
        if name == "_FailedTest":
            name = ""
        if name:
            self.stream.writeln(_group(name))
        startTime = perf_counter()
        result = CodewarsTestResult(self.stream)
        try:
            test(result)
        finally:
            pass
        if name:
            self.stream.writeln(_completedin(startTime, perf_counter()))
        self.results.append(result)

    def _run_case(self, test):
        result = CodewarsTestResult(self.stream)
        try:
            test(result)
        finally:
            pass
        self.results.append(result)


class _NamedTestSuite(unittest.TestSuite):
    def __init__(self, tests=(), group_name=None):
        super(_NamedTestSuite, self).__init__(tests)
        self.group_name = group_name


def _group(name):
    return "\n<DESCRIBE::>{}".format(name)


def _completedin(start, end):
    return "\n<COMPLETEDIN::>{:.4f}".format(1000 * (end - start))


# Flatten nested TestSuite by collecting all test cases.
def _flatten(suites):
    tests = []
    for test in suites:
        if isinstance(test, unittest.TestSuite):
            tests.extend(_flatten(test))
        else:
            tests.append(test)
    return tests


# Group by module name and then by class name
def _to_tree(suite):
    tree = unittest.TestSuite()
    for k, ms in groupby(suite, _module_name):
        sub_trees = _NamedTestSuite(group_name=k)
        for c, cs in groupby(ms, _class_name):
            sub_trees.addTest(_NamedTestSuite(tests=cs, group_name=c))
        tree.addTest(sub_trees)
    return tree


def _module_name(x):
    return x.__class__.__module__


def _class_name(x):
    return x.__class__.__name__


# https://github.com/python/cpython/blob/289f1f80ee87a4baf4567a86b3425fb3bf73291d/Lib/unittest/runner.py#L13
class _WritelnDecorator(object):
    """Used to decorate file-like objects with a handy 'writeln' method"""

    def __init__(self, stream):
        self.stream = stream

    def __getattr__(self, attr):
        if attr in ("stream", "__getstate__"):
            raise AttributeError(attr)
        return getattr(self.stream, attr)

    def writeln(self, arg=None):
        if arg:
            self.write(arg)
        self.write("\n")  # text-mode streams translate to \r\n if needed
