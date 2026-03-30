// Package memory contains a file system implementation [Dir] stores data on
// memory, and its utilities that together allow its use with the standard
// library utilities such as [fs.WalkDir] or [template.Template]. It is
// intended to be used in testing only.
//
// [Dir] partially conforms the fstest.TestFS expectations. It supports
// absolute paths and redundant segments as builder utilizes. It doesn't
// support file timestamps and custom permissions.
package memory
