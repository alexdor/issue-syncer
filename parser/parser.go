package parser

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type MultiLineMarker struct {
	Start string
	End   string
}

// TODO: flip these types to use bytes instead of string
type CommentMarker struct {
	SingleLine []string
	MultiLine  []MultiLineMarker
}

type commentType int

const (
	singleLine commentType = 0
	multiLine  commentType = 1
)

type Comment struct {
	FilePath      string
	LineNumber    int
	LineNumberEnd int
	Text          string
	CommentType   commentType
}

func shouldIncludeTheComment(line string, wordsToLookFor []string) bool {
	line = strings.ToLower(line)
	for _, word := range wordsToLookFor {
		// TODO: this should be optimized to look only in the start of the string
		if strings.Contains(line, word) {
			return true
		}
	}
	return false
}

func parseFile(filePath string, marker CommentMarker, wordsToLookFor []string) ([]Comment, error) {
	var comments []Comment
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Warn("failed to open file, skipping", "file", filePath)
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	var multiLineComment *Comment
	var currentMultiLineMarker *MultiLineMarker

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// If we're in a multiline comment
		if multiLineComment != nil {
			multiLineEnd := strings.Index(line, currentMultiLineMarker.End)
			if multiLineEnd != -1 {
				multiLineComment.LineNumberEnd = lineNum
				multiLineComment.Text += line[:multiLineEnd]
				comments = append(comments, *multiLineComment)
				multiLineComment = nil
				currentMultiLineMarker = nil
				continue
			}
			multiLineComment.Text += line + "\n"
			continue
		}

		// Check for new multiline comments
		for _, marker := range marker.MultiLine {
			multiLineStart := strings.Index(line, marker.Start)
			if multiLineStart == -1 || !shouldIncludeTheComment(line[multiLineStart:], wordsToLookFor) {
				continue
			}

			currentMultiLineMarker = &marker
			multiLineComment = &Comment{
				FilePath:    filePath,
				LineNumber:  lineNum,
				Text:        line[multiLineStart:] + "\n",
				CommentType: multiLine,
			}

			// If start and end are on the same line
			multiLineEnd := strings.Index(multiLineComment.Text, marker.End)
			if multiLineEnd != -1 {
				multiLineComment.LineNumberEnd = lineNum
				multiLineComment.Text = multiLineComment.Text[:multiLineEnd]
				comments = append(comments, *multiLineComment)
				multiLineComment = nil
				currentMultiLineMarker = nil
			}
			break
		}

		if multiLineComment != nil {
			continue
		}

		// Check for single-line comments
		for _, marker := range marker.SingleLine {
			commentStart := strings.Index(line, marker)
			if commentStart == -1 || !shouldIncludeTheComment(line[commentStart:], wordsToLookFor) {
				continue
			}

			comments = append(comments, Comment{
				FilePath:      filePath,
				LineNumber:    lineNum,
				LineNumberEnd: lineNum,
				Text:          line[commentStart:],
				CommentType:   singleLine,
			})
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return comments, fmt.Errorf("error from scanner on file %s, err: %w", filePath, err)
	}

	return comments, nil
}

// ParseDirectory walks through a directory and finds all comments
func ParseDirectory(rootPath string, wordsToLookFor []string, pathsToSkip []string, useGitIgnore bool) ([]Comment, error) {
	var comments []Comment

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if slices.Contains(pathsToSkip, d.Name()) {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		marker, ok := commentMarkersWithEnds[ext]
		if !ok {
			return nil
		}

		fileComments, err := parseFile(path, marker, wordsToLookFor)
		if err != nil {
			if !errors.Is(err, bufio.ErrTooLong) {
				return err
			}
			slog.Warn("failed to parse file because of too long of a string, continuing", "path", path)
		}

		comments = append(comments, fileComments...)
		return nil
	})

	return comments, err
}
