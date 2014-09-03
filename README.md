gar - Go application archiver
=============================

**gar** is an application archiver for Go inspired by JAR for Java.  All
external resources required for running an application, such as templates and
public assets for a web application. are archived directly with the binary to
create a completely portable application.

Comparison with go.rice
=======================

[go.rice](https://github.com/GeertJohan/go.rice) is another Go archiver solving
the same problems gar aims to solve, however there are some design differences
between the two.

* go.rice uses zip for both archiving and compression, whereas gar uses
  [tar archiving](http://en.wikipedia.org/wiki/Tar_\(computing\)) with optional
  [gzip compression](http://en.wikipedia.org/wiki/Gzip) so the packager can
  choose to optimise either for loading performance or package size.  
* go.rice relies on an external tool for archiving and compressing and
  currently does not work under Windows.  gar uses Go's core
  [tar](http://golang.org/pkg/archive/tar/) and
  [gzip](http://golang.org/pkg/compress/gzip/) packages which makes it portable
  to all systems supported by Go.
* go.rice loads all resources into memory which can be problematic for
  applications with large resources.  gar supports loading into memory by
  default, but can configured to load to the file system instead to save
  memory.
* go.rice supports source generation which provides another packaging option,
  gar only supports appending an archive to the binary.
