# Sorted Arrays in Go

[![CircleCI](https://img.shields.io/circleci/build/github/jfsmig/go-bags/main)](https://app.circleci.com/pipelines/github/jfsmig/go-bags)
[![Codacy](https://app.codacy.com/project/badge/Grade/aa58726a923b40e6a92fdacd77a344ae)](https://www.codacy.com/gh/jfsmig/go-bags/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=jfsmig/go-bags&amp;utm_campaign=Badge_Grade)
[![CodeCov](https://img.shields.io/codecov/c/github/jfsmig/go-bags)](https://app.codecov.io/gh/jfsmig/go-bags)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

Sorted array provide a complexity profile which make it suitable for collection with significantly more lookups that modifying operations:
*   a compact memory footprint
*   an efficient lookup complexity in O(log N)
*   an efficient scan complexity since it depends on the lookup followed by a sequential scan of the array
*   an insertion in O(N * log N) which is rather inefficient but remains acceptable if the operation is rather rare

3 flavors of generic sorted arrays for efficient lookup and paginated scans.
