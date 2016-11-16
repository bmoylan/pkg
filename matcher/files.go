/*
Copyright 2016 Palantir Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package matcher

import (
	"fmt"
	"os"
	"path/filepath"
)

// ListFiles returns the files in the provided directory (relative or absolute path) that match the provided include
// matcher but do not match the exclude matcher. The provided directory is used as the base directory and the listing is
// done recursively. The paths that are returned are relative to the input directory.
func ListFiles(dir string, include, exclude Matcher) ([]string, error) {
	dirAbsPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to convert path %s to absolute path", dir)
	}

	if fileInfo, err := os.Stat(dirAbsPath); err != nil {
		return nil, fmt.Errorf("failed to stat %s", dirAbsPath)
	} else if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dirAbsPath)
	}

	var paths []string
	if err := filepath.Walk(dirAbsPath, func(path string, currInfo os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk failed at %s", path)
		}

		currRelPath, err := filepath.Rel(dirAbsPath, path)
		if err != nil {
			return fmt.Errorf("failed to resolve %s to relative path against base %s", path, dirAbsPath)
		}

		// if current path matches an include and does not match any excludes, include
		if include != nil && include.Match(currRelPath) && (exclude == nil || !exclude.Match(currRelPath)) {
			paths = append(paths, currRelPath)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return paths, nil
}
