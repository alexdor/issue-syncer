package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/alexdor/issue-syncer/parser"
	"github.com/alexdor/issue-syncer/storer"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

var (
	Path           string
	WordsToLookFor []string
	DirsToSkip     []string
	UseGitIgnore   bool
	Storer         string
	DryRun         bool
	Version        string
)

func init() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.RFC3339,
		}),
	))
	rootCmd.Flags().StringVarP(&Path, "path", "p", ".", "Path to the folder to scan")
	rootCmd.Flags().StringSliceVarP(
		&WordsToLookFor, "words", "w", []string{"FIXME", "TODO", "HACK"}, "Words to look for in comments",
	)
	rootCmd.Flags().StringSliceVarP(&DirsToSkip, "dirs-to-skip", "d", defaultDirsToSkip, "Dirs to skip by default")
	// TODO: add a flag to use the content of gitignore to not search in them
	// rootCmd.Flags().BoolVarP(
	// 	&UseGitIgnore, "use-gitignore", "g", true, "Wether to use the content of gitignore to not search in them, or not",
	// )
	rootCmd.Flags().StringVarP(&Storer, "storer", "s", "github",
		"Storer to use for checking and updating issues, available: "+
			strings.Join(slices.Collect(maps.Keys(storer.AvailableStorer)), ","),
	)
	rootCmd.Flags().BoolVar(&DryRun, "dry-run", false, "Wether to do a dry run or not")
}

var defaultDirsToSkip = []string{
	".git", "node_modules", ".cache", ".next", "_next", ".vscode", "dist", "out", "build", ".tmp", ".idea",
}

var rootCmd = &cobra.Command{
	Use:     "issue-syncer",
	Short:   "A tool to sync TODO comments with issues on GitHub",
	Version: Version,
	Long:    ``,
	RunE: func(cob *cobra.Command, _ []string) error {
		storerToUse, ok := storer.AvailableStorer[strings.ToLower(Storer)]
		if !ok {
			return errors.New(
				"storer " + Storer + " is not available, available: " +
					strings.Join(slices.Collect(maps.Keys(storer.AvailableStorer)), ","),
			)
		}

		if DryRun {
			slog.Info("Running in dry run mode, no changes will be made")
			cob.SetContext(context.WithValue(cob.Context(), storer.DryRunKey, true))
		}

		for i := range WordsToLookFor {
			WordsToLookFor[i] = strings.ToLower(WordsToLookFor[i])
		}
		comments, err := parser.ParseDirectory(Path, WordsToLookFor, DirsToSkip, UseGitIgnore)
		if err != nil {
			return fmt.Errorf("failed to parse directory: %w", err)
		}

		err = storerToUse.Init(cob.Context())
		if err != nil {
			return fmt.Errorf("failed to create storer: %w", err)
		}

		currentIssues, err := storerToUse.FetchCurrentOpenIssues(cob.Context())
		if err != nil {
			return fmt.Errorf("failed to fetch current open issues: %w", err)
		}

		return storer.UpdateIssues(cob.Context(), storerToUse, currentIssues, comments)
	},
}

func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-signalChan:
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		os.Exit(2)
	}()

	err := rootCmd.ExecuteContext(ctx)
	cancel()
	if err != nil {
		os.Exit(1)
	}
}
