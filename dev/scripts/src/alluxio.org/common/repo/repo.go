/*
 * The Alluxio Open Foundation licenses this work under the Apache License, version 2.0
 * (the "License"). You may not use this work except in compliance with the License, which is
 * available at www.apache.org/licenses/LICENSE-2.0
 *
 * This software is distributed on an "AS IS" basis, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied, as more fully set forth in the License.
 *
 * See the NOTICE file distributed with this work for information regarding copyright ownership.
 */

package repo

import (
	"log"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/palantir/stacktrace"
)

var (
	root         string
	repoRootOnce sync.Once
)

func FindRepoRoot() string {
	repoRootOnce.Do(func() {
		// navigate 7 parent directories to reach repo root,
		// assuming this go file is located in <repoRoot>/dev/scripts/src/alluxio.org/common/repo/repo.go
		const repoRootDepth = 7
		_, r, _, ok := runtime.Caller(0)
		if !ok {
			panic(stacktrace.NewError("error getting call stack to find repo root"))
		}
		for i := 0; i < repoRootDepth; i++ {
			r = filepath.Dir(r)
		}
		log.Printf("Repository root at directory: %v", r)
		root = r
	})
	return root
}
