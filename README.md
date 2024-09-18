# Radix Tree Search

## Overview

This library implements a reverse pattern search algorithm using a Radix tree data structure. Unlike conventional
pattern matching where a string is matched against a pattern, this reverse search identifies patterns that match a given
input string.

## Use Case Example

Consider a set of patterns:

```
abc
*cd
*bc
b*d
*
```

For an input string "abc", the search would yield:

```
abc
*bc
*
```

## Applications and Implementation

Reverse pattern search is particularly useful for browser database lookups, such as [Browscap](https://browscap.org/).
The library utilizes a Radix tree as a storage and backtracking to traverse the tree. Since the whole tree structure is
stored in memory, it's possible to convert tree into a Directed Acyclic Word Graph (DAWG) which can significantly reduce
memory usage.

## Installation

```bash
go get github.com/eugeniypetrov/radix-tree
```

## Usage

```go
r := radix.NewRadix()
r.Add("abc")
r.Add("*cd")
r.Add("*bc")
r.Add("b*d")
r.Add("*")

for _, v := range r.Find("abc") {
    log.Println(v)
}
```

To convert into DAWG:

```go
d := r.ToDAWG()
d.Find("abc")
```