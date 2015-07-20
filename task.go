/*
 * Copyright (c) 2015 Alex Yatskov <alex@foosoft.net>
 * Author: Alex Yatskov <alex@foosoft.net>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import "fmt"

type task struct {
	Deps  []string
	Links [][]string
	Cmds  [][]string
	Envs  [][]string
}

func (t *task) process(srcDir, dstDir string, conf *config, flags int) error {
	for _, currTask := range t.Deps {
		if err := processTask(currTask, srcDir, dstDir, conf, flags); err != nil {
			return err
		}
	}

	for _, currEnv := range t.Envs {
		if err := processEnv(currEnv, flags); err != nil {
			return err
		}
	}

	if flags&flagNoLink == 0 {
		for _, currLink := range t.Links {
			if err := processLink(currLink, srcDir, dstDir, flags); err != nil {
				return err
			}
		}
	}

	if flags&flagNoCmd == 0 {
		for _, currCmd := range t.Cmds {
			if err := processCmd(currCmd, dstDir, conf, flags); err != nil {
				return err
			}
		}
	}

	return nil
}

func processTask(taskName, srcDir, dstDir string, conf *config, flags int) error {
	handled, ok := conf.tasksHandled[taskName]
	if ok && handled {
		return nil
	}

	conf.tasksHandled[taskName] = true

	t, ok := conf.Tasks[taskName]
	if !ok {
		return fmt.Errorf("task not found %s", taskName)
	}

	return t.process(srcDir, dstDir, conf, flags)
}
